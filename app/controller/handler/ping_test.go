package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bartick/go-task/app/controller/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := gin.New()
	router.GET("/ping", handler.HandlerPing) // Register your Ping route

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"pong"`)
}
