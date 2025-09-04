package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	null "github.com/mattn/go-nulltype"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

func (s *TaskStatus) Scan(value interface{}) error {
	if value == nil {
		*s = StatusTodo
		return nil
	}

	if str, ok := value.([]byte); ok {
		*s = TaskStatus(str)
		return nil
	}
	return fmt.Errorf("cannot scan %T into TaskStatus", value)
}

func (s TaskStatus) Value() (driver.Value, error) {
	return string(s), nil
}

type Task struct {
	ID           int64           `json:"id" db:"id"`
	Title        string          `json:"title" db:"title"`
	Description  null.NullString `json:"description" db:"description"`
	Status       TaskStatus      `json:"status" db:"status"`
	Priority     int8            `json:"priority" db:"priority"`
	DueDate      null.NullTime   `json:"due_date" db:"due_date"`
	CompletedAt  null.NullTime   `json:"completed_at" db:"completed_at"`
	ParentTaskID null.NullInt64  `json:"parent_task_id" db:"parent_task_id"`
	CategoryID   null.NullInt64  `json:"category_id" db:"category_id"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at" db:"updated_at"`
}

type TaskWithCategory struct {
	Task
	CategoryName *string `json:"category_name" db:"category_name"`
}

type TaskHierarchy struct {
	Task
	CategoryName *string         `json:"category_name" db:"category_name"`
	Subtasks     []TaskHierarchy `json:"subtasks,omitempty"`
}

type CreateTaskRequest struct {
	Title        string          `json:"title"`
	Description  null.NullString `json:"description"`
	Status       *TaskStatus     `json:"status"`
	Priority     int8            `json:"priority"`
	DueDate      null.NullTime   `json:"due_date"`
	ParentTaskID null.NullInt64  `json:"parent_task_id"`
	CategoryName null.NullString `json:"category_name"`
}

type UpdateTaskRequest struct {
	Title        null.NullString `db:"title"`
	Description  null.NullString `db:"description"`
	Status       null.NullString `db:"status"`
	Priority     null.NullInt64  `db:"priority"`
	DueDate      null.NullTime   `db:"due_date"`
	CompletedAt  null.NullTime   `db:"completed_at"`
	ParentTaskID null.NullInt64  `db:"parent_task_id"`
	CategoryName null.NullString `db:"category_name"`
}

const (
	queryAllGetTasks = `
		SELECT 
			t.id, t.title, t.description, t.status, t.priority, 
			t.due_date, t.completed_at, t.parent_task_id, t.category_id,
			t.created_at, t.updated_at, c.name as category_name
		FROM tasks t
		LEFT JOIN categories c ON t.category_id = c.id
	`

	queryGetTaskHierarchy = `
	WITH RECURSIVE task_hierarchy AS (
		SELECT 
			id, title, description, status, priority, 
			due_date, completed_at, parent_task_id, category_id,
			created_at, updated_at
		FROM tasks
		WHERE id = ?

		UNION ALL

		SELECT 
			t.id, t.title, t.description, t.status, t.priority, 
			t.due_date, t.completed_at, t.parent_task_id, t.category_id,
			t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN task_hierarchy th ON t.parent_task_id = th.id
	)
	SELECT 
		th.*, c.name as category_name
	FROM task_hierarchy th
	LEFT JOIN categories c ON th.category_id = c.id
	ORDER BY th.priority DESC, th.created_at ASC;
	`
	queryCreateTask = `
	INSERT INTO tasks (title, description, status, priority, due_date, parent_task_id, category_id)
	VALUES (:title, :description, :status, :priority, :due_date, :parent_task_id, (SELECT id FROM categories WHERE name = :category_name))
	`

	queryUpdateTask = `
	UPDATE tasks SET 
	title = COALESCE(:title, title), 
	description = COALESCE(:description, description), 
	status = COALESCE(:status, status), 
	priority = COALESCE(:priority, priority), 
	due_date = COALESCE(:due_date, due_date), 
	parent_task_id = COALESCE(:parent_task_id, parent_task_id), 
	category_id = COALESCE((SELECT id FROM categories WHERE name = :category_name), category_id)
	WHERE id = :id
	`

	queryDeleteTask = `
	WITH RECURSIVE task_hierarchy AS (
		SELECT id, parent_task_id 
		FROM tasks 
		WHERE id = ?
		
		UNION ALL
		
		SELECT t.id, t.parent_task_id
		FROM tasks t
		INNER JOIN task_hierarchy th ON t.parent_task_id = th.id
	)
	DELETE FROM tasks 
	WHERE id IN (SELECT id FROM task_hierarchy);
	`
)

func GetAllTasks(db *sqlx.DB) ([]TaskWithCategory, error) {
	var tasks []TaskWithCategory
	err := db.Select(&tasks, queryAllGetTasks)
	return tasks, err
}

func GetTaskWithSubtasks(db *sqlx.DB, taskID int64) (*TaskHierarchy, error) {
	// Step 1: Fetch all tasks in the hierarchy (parent and all descendants) into a flat list.
	var flatTasks []TaskHierarchy
	if err := db.Select(&flatTasks, queryGetTaskHierarchy, taskID); err != nil {
		return nil, err // This could be sql.ErrNoRows or a real error
	}

	// Step 2: Handle the case where the root task ID doesn't exist.
	if len(flatTasks) == 0 {
		return nil, sql.ErrNoRows
	}

	// Step 3: Build the tree structure from the flat list.
	taskMap := make(map[int64]*TaskHierarchy)
	for i := range flatTasks {
		task := &flatTasks[i]
		taskMap[task.ID] = task
	}

	// Iterate through the tasks again to link children to their parents.
	for _, task := range taskMap {
		// Check if the task has a valid parent ID.
		if task.ParentTaskID.Valid() {
			// Find the parent in the map using the parent ID.
			if parent, ok := taskMap[task.ParentTaskID.Int64Value()]; ok {
				// **THE FIX**: Append the pointer to the task (`task`), not a copy (`*task`).
				parent.Subtasks = append(parent.Subtasks, *task)
			}
		}
	}

	// Find the root of the tree (the task with the original ID) and return it.
	if root, ok := taskMap[taskID]; ok {
		return root, nil
	}

	return nil, sql.ErrNoRows
}

func CreateTask(db *sqlx.DB, req *CreateTaskRequest) (*Task, error) {

	result, err := db.NamedExec(queryCreateTask, map[string]interface{}{
		"title":          req.Title,
		"description":    req.Description,
		"status":         req.Status,
		"priority":       req.Priority,
		"due_date":       req.DueDate,
		"parent_task_id": req.ParentTaskID,
		"category_name":  req.CategoryName,
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetByID(db, id)
}

func UpdateTask(db *sqlx.DB, taskID uint64, updates *UpdateTaskRequest) (int64, error) {

	res, err := db.NamedExec(queryUpdateTask, map[string]interface{}{
		"id":             taskID,
		"title":          updates.Title,
		"description":    updates.Description,
		"status":         updates.Status,
		"priority":       updates.Priority,
		"due_date":       updates.DueDate,
		"completed_at":   updates.CompletedAt,
		"parent_task_id": updates.ParentTaskID,
		"category_name":  updates.CategoryName,
	})
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteTask(db *sqlx.DB, taskID uint64) (int64, error) {
	req, err := db.Exec(queryDeleteTask, taskID)
	if err != nil {
		return 0, err
	}
	return req.RowsAffected()
}

func GetByID(db *sqlx.DB, taskID int64) (*Task, error) {
	query := `
        SELECT id, title, description, status, priority, due_date, 
               completed_at, parent_task_id, category_id, created_at, updated_at
        FROM tasks WHERE id = ?`

	var task Task
	err := db.Get(&task, query, taskID)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
