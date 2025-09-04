package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bartick/go-task/app/database"
	"github.com/bartick/go-task/app/route"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func runTaskAPI() {
	gin.SetMode(gin.DebugMode)

	fmt.Println(config)

	// Database Init
	db, err := database.InitDatabases(config.Database)
	if err != nil {
		log.Fatal("Database initialization failed", zap.String("err", err.Error()))
	}

	// Start HTTP server
	router := route.AddAPIRouter(db)

	serverAddr := config.Server.Address + ":" + config.Server.Port
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a goroutine so that it doesn't block.
	go func() {
		log.Info("Starting server at http://%s\n", zap.String("addr", serverAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen: %s\n", zap.String("err", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", zap.String("err", err.Error()))
	}

	if err := db.Close(); err != nil {
		log.Error("Database connection close failed:", zap.String("err", err.Error()))
	}

	log.Info("Server exiting")
}
