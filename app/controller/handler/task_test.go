package handler_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bartick/go-task/app/controller/handler"
	"github.com/bartick/go-task/app/model"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-nulltype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerGetTask_Success(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Dummy task to return
	dummyTask := &model.Task{
		ID:    1,
		Title: "Test Task",
	}

	// Expect GetByID to be called inside the handler
	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(dest interface{}, query string, args ...interface{}) error {
			// dest is a pointer to a Task
			taskPtr := dest.(*model.Task)
			*taskPtr = *dummyTask
			return nil
		})

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks/:id", handler.HandlerGetTask)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"title":"Test Task"`)
}

func TestHandlerGetTask_Failure(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Expect GetByID to be called and return an error
	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, mock.Anything).
		Return(sql.ErrNoRows)

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks/:id", handler.HandlerGetTask)

	// Make HTTP request with invalid ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/999", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestHandlerGetTask_InvalidID(t *testing.T) {
	// Create router without DB since it won't be used
	router := gin.New()
	router.GET("/tasks/:id", handler.HandlerGetTask)

	// Make HTTP request with invalid ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid task ID")
}

func TestHandlerGetTasks_Success(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Dummy tasks to return
	dummyTasks := []model.TaskWithCategory{
		{Task: model.Task{ID: 1, Title: "Task 1"}},
		{Task: model.Task{ID: 2, Title: "Task 2"}},
	}

	// Expect Select to be called inside the handler
	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(dest interface{}, query string, args ...interface{}) error {
			tasksPtr := dest.(*[]model.TaskWithCategory) // <- important
			*tasksPtr = dummyTasks
			return nil
		})

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks", handler.HandlerGetTasks)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"title":"Task 1"`)
	assert.Contains(t, w.Body.String(), `"title":"Task 2"`)
}

func TestHandlerGetTasks_Failure(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Expect Select to be called and return an error
	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything, mock.Anything).
		Return(sql.ErrConnDone)

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks", handler.HandlerGetTasks)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve tasks")
}

func TestHandlerGetSubTasks_Success(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Dummy subtasks to return (slice, not struct)
	dummySubTasks := []model.TaskHierarchy{
		{
			Task: model.Task{ID: 1, Title: "Parent Task"},
		},
		{
			Task: model.Task{ID: 2, Title: "Subtask 1", ParentTaskID: nulltype.NullInt64Of(1)},
		},
		{
			Task: model.Task{ID: 3, Title: "Subtask 2", ParentTaskID: nulltype.NullInt64Of(1)},
		},
	}

	// Expect Select to be called inside the handler
	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(dest interface{}, query string, args ...interface{}) error {
			subTasksPtr := dest.(*[]model.TaskHierarchy) // must match
			*subTasksPtr = dummySubTasks
			return nil
		})

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks/:id/subtasks", handler.HandlerGetSubTasks)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1/subtasks", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"title":"Subtask 1"`)
	assert.Contains(t, w.Body.String(), `"title":"Subtask 2"`)
}

func TestHandlerGetSubTasks_Failure(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Expect Select to be called and return an error
	mockDB.EXPECT().
		Select(mock.Anything, mock.Anything, mock.Anything).
		Return(sql.ErrConnDone)

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.GET("/tasks/:id/subtasks", handler.HandlerGetSubTasks)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1/subtasks", nil)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve task subtasks")
}

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

func TestHandlerCreateTasks_Success(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Dummy request and response
	reqBody := `{"title":"New Task","description":"Task description"}`
	createdTask := &model.Task{
		ID:          1,
		Title:       "New Task",
		Description: nulltype.NullStringOf("Task description"),
	}

	// Expect NamedExec (because your INSERT uses named parameters)
	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(&mockResult{lastInsertID: 1, rowsAffected: 1}, nil)

	// Expect Get to fetch the created task
	mockDB.EXPECT().
		Get(mock.Anything, mock.Anything, mock.Anything).
		RunAndReturn(func(dest interface{}, query string, args ...interface{}) error {
			taskPtr := dest.(*model.Task)
			*taskPtr = *createdTask
			return nil
		})

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.POST("/tasks", handler.HandlerCreateTasks)

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"title":"New Task"`)
}

func TestHandlerCreateTasks_InvalidJSON(t *testing.T) {
	// Create router without DB since it won't be used
	router := gin.New()
	router.POST("/tasks", handler.HandlerCreateTasks)

	// Make HTTP request with invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(`{"title":`)) // malformed JSON
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `{"error":"unexpected EOF"}`)
}

func TestHandlerUpdateTask_Success(t *testing.T) {
	// Initialize mock DB
	mockDB := model.NewMockDBTX(t)

	// Expect only NamedExec (Update)
	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 1}, nil)

	// Create router with the mock DB
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB) // Inject mock DB into context
	})
	router.PATCH("/tasks/:id", handler.HandlerUpdateTask)

	// Dummy request body
	reqBody := `{"title":"Updated Task"}`

	// Make HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Task updated successfully"`)
}

func TestHandlerUpdateTask_InvalidID(t *testing.T) {
	router := gin.New()
	router.PATCH("/tasks/:id", handler.HandlerUpdateTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/abc", strings.NewReader(`{"title":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"Invalid task ID"`)
}

func TestHandlerUpdateTask_InvalidJSON(t *testing.T) {
	router := gin.New()
	router.PATCH("/tasks/:id", handler.HandlerUpdateTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", strings.NewReader(`{invalid json}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestHandlerUpdateTask_NoFields(t *testing.T) {
	router := gin.New()
	router.PATCH("/tasks/:id", handler.HandlerUpdateTask)

	// Empty body (all fields invalid)
	reqBody := `{}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"No fields to update"`)
}

func TestHandlerUpdateTask_NotFound(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	// Update executes but affects 0 rows
	mockDB.EXPECT().
		NamedExec(mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 0}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB)
	})
	router.PATCH("/tasks/:id", handler.HandlerUpdateTask)

	reqBody := `{"title":"Updated Task"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/tasks/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"Task not found"`)
}

func TestHandlerDeleteTask_Success(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	// Expect delete to affect 1 row
	mockDB.EXPECT().
		Exec(mock.Anything, mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 1}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB)
	})
	router.DELETE("/tasks/:id", handler.HandlerDeleteTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"Task deleted successfully"`)
}

func TestHandlerDeleteTask_InvalidID(t *testing.T) {
	router := gin.New()
	router.DELETE("/tasks/:id", handler.HandlerDeleteTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"Invalid task ID"`)
}

func TestHandlerDeleteTask_NotFound(t *testing.T) {
	mockDB := model.NewMockDBTX(t)

	// Exec returns 0 rows affected
	mockDB.EXPECT().
		Exec(mock.Anything, mock.Anything, mock.Anything).
		Return(&mockResult{rowsAffected: 0}, nil)

	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("db", mockDB)
	})
	router.DELETE("/tasks/:id", handler.HandlerDeleteTask)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"Task not found"`)
}
