package route

import (
	"github.com/bartick/go-task/app/controller/handler"
	"github.com/bartick/go-task/app/route/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

const (
	// Ping
	pathPing = "/ping"

	// Tasks
	pathTasks    = "/tasks"
	pathTasksID  = "/tasks/:id"
	pathSubTasks = "/tasks/:id/subtasks"
)

func AddAPIRouter(db *sqlx.DB) *gin.Engine {
	router := gin.New()
	router.Use(middleware.LogRequest(log))
	router.Use(middleware.CORSMiddleware())

	// Ping
	router.GET(pathPing, handler.HandlerPing)

	router.Use(middleware.Config(db))

	// Tasks
	router.GET(pathTasks, handler.HandlerGetTasks)
	router.GET(pathSubTasks, handler.HandlerGetSubTasks)
	router.POST(pathTasks, handler.HandlerCreateTasks)
	router.GET(pathTasksID, handler.HandlerGetTask)
	router.PATCH(pathTasksID, handler.HandlerUpdateTask)
	router.DELETE(pathTasksID, handler.HandlerDeleteTask)

	return router
}
