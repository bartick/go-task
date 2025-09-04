package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAllTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/9999", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetTaskInvalidID(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/abc", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetSubTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/tasks/1/subtasks", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreateTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":"New Task","description":"Task description","priority":3,"status":"in_progress"}`
	req, _ := http.NewRequest("POST", "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreateTaskInvalidBody(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":123,"description":true}`
	req, _ := http.NewRequest("POST", "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":"Updated Task","description":"Updated description","priority":2,"status":"done"}`
	req, _ := http.NewRequest("PATCH", "/tasks/3", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateTaskNotFound(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":"Updated Task","description":"Updated description","priority":2,"status":"done"}`
	req, _ := http.NewRequest("PATCH", "/tasks/9999", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateTaskInvalidID(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":"Updated Task","description":"Updated description","priority":2,"status":"done"}`
	req, _ := http.NewRequest("PATCH", "/tasks/abc", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTaskInvalidBody(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{"title":123,"description":true}`
	req, _ := http.NewRequest("PATCH", "/tasks/2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateTaskNoFields(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	reqBody := `{}`
	req, _ := http.NewRequest("PATCH", "/tasks/2", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(reqBody))
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/9999", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteTaskInvalidID(t *testing.T) {
	initTestEnvironment()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/tasks/abc", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
