package model_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bartick/go-task/app/model"
	"github.com/mattn/go-nulltype"
	mock "github.com/stretchr/testify/mock"
	"github.com/zeebo/assert"
)

func TestGetByID_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	expected := &model.Task{
		ID:          1,
		Title:       "Mock Task",
		Description: nulltype.NullStringOf("This is a mock task"),
		Status:      model.StatusTodo,
		Priority:    1,
		DueDate:     nulltype.NullTimeOf(time.Now().Add(24 * time.Hour)),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, []interface{}{int64(1)}).
		Run(func(dest interface{}, query string, args ...interface{}) {
			d := dest.(*model.Task)
			*d = *expected
		}).
		Return(nil)

	got, err := model.GetByID(mockDB, 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetByID_NotFound(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, []interface{}{int64(999)}).
		Return(sql.ErrNoRows)

	got, err := model.GetByID(mockDB, 999)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestGetAllTasks_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	expected := []model.TaskWithCategory{
		{
			Task: model.Task{
				ID:          1,
				Title:       "Task 1",
				Description: nulltype.NullStringOf("First task"),
				Status:      model.StatusTodo,
				Priority:    1,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			CategoryName: func() *string { s := "Backend"; return &s }(),
		},
		{
			Task: model.Task{
				ID:          2,
				Title:       "Task 2",
				Description: nulltype.NullStringOf("Second task"),
				Status:      model.StatusInProgress,
				Priority:    2,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			CategoryName: nil,
		},
	}

	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything).
		Run(func(dest interface{}, query string, args ...interface{}) {
			d := dest.(*[]model.TaskWithCategory)
			*d = expected
		}).
		Return(nil)

	tasks, err := model.GetAllTasks(mockDB)

	assert.NoError(t, err)
	assert.Equal(t, expected, tasks)
}

func TestGetAllTasks_DBError(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything).
		Return(sql.ErrConnDone)

	tasks, err := model.GetAllTasks(mockDB)

	assert.Error(t, err)
	assert.Nil(t, tasks)
	assert.Equal(t, sql.ErrConnDone, err)
}

// mockResult implements sql.Result for testing
type mockResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (m *mockResult) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}
func (m *mockResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func TestCreateTask_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	status := model.StatusTodo
	req := &model.CreateTaskRequest{
		Title:        "New Task",
		Description:  nulltype.NullStringOf("This is a new task"),
		Status:       &status,
		Priority:     2,
		DueDate:      nulltype.NullTimeOf(time.Now().Add(48 * time.Hour)),
		ParentTaskID: nulltype.NullInt64{},
		CategoryName: nulltype.NullStringOf("Work"),
	}

	// Step 1: Stub NamedExec to return fake result
	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(&mockResult{lastInsertID: 1}, nil)

	// Step 2: Stub GetByID call (CreateTask calls GetByID after insert)
	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, []interface{}{int64(1)}).
		Run(func(dest interface{}, query string, args ...interface{}) {
			task := dest.(*model.Task)
			task.ID = 1
			task.Title = req.Title
			task.Description = req.Description
			task.Status = *req.Status
			task.Priority = req.Priority
		}).
		Return(nil)

	task, err := model.CreateTask(mockDB, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), task.ID)
	assert.Equal(t, "New Task", task.Title)
}

func TestCreateTask_DBError(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	status := model.StatusTodo
	req := &model.CreateTaskRequest{
		Title:        "New Task",
		Description:  nulltype.NullStringOf("This is a new task"),
		Status:       &status,
		Priority:     2,
		DueDate:      nulltype.NullTimeOf(time.Now().Add(48 * time.Hour)),
		ParentTaskID: nulltype.NullInt64{},
		CategoryName: nulltype.NullStringOf("Work"),
	}

	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(nil, sql.ErrConnDone)
	task, err := model.CreateTask(mockDB, req)

	assert.Error(t, err)
	assert.Nil(t, task)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestUpdateTask_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	status, _ := model.StatusInProgress.Value()
	completedAt := nulltype.NullTimeOf(time.Now())
	req := &model.UpdateTaskRequest{
		Title:       nulltype.NullStringOf("Updated Task"),
		Description: nulltype.NullStringOf("This task has been updated"),
		Status:      nulltype.NullStringOf(status.(string)),
		Priority:    nulltype.NullInt64Of(3),
		DueDate:     nulltype.NullTimeOf(time.Now().Add(72 * time.Hour)),
		CompletedAt: completedAt,
	}

	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 1}, nil)

	rowsAffected, err := model.UpdateTask(mockDB, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
}

func TestUpdateTask_DBError(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	status, _ := model.StatusInProgress.Value()
	req := &model.UpdateTaskRequest{
		Title:       nulltype.NullStringOf("Updated Task"),
		Description: nulltype.NullStringOf("This task has been updated"),
		Status:      nulltype.NullStringOf(status.(string)),
		Priority:    nulltype.NullInt64Of(3),
		DueDate:     nulltype.NullTimeOf(time.Now().Add(72 * time.Hour)),
	}

	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(nil, sql.ErrConnDone)

	rowsAffected, err := model.UpdateTask(mockDB, 1, req)

	assert.Error(t, err)
	assert.Equal(t, int64(0), rowsAffected)
	assert.Equal(t, sql.ErrConnDone, err)
}

func TestDeleteTask_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	mockDB.EXPECT().
		Exec(mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 1}, nil)
	rowsAffected, err := model.DeleteTask(mockDB, 1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)
}

func TestDeleteTask_DBError(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	mockDB.EXPECT().
		Exec(mock.Anything, mock.Anything).
		Return(nil, sql.ErrConnDone)
	rowsAffected, err := model.DeleteTask(mockDB, 1)

	assert.Error(t, err)
	assert.Equal(t, int64(0), rowsAffected)
	assert.Equal(t, sql.ErrConnDone, err)
}
