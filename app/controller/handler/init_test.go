package handler_test

import (
	"sync"

	"github.com/bartick/ringover-task/app/database"
	"github.com/bartick/ringover-task/app/model"
	"github.com/bartick/ringover-task/app/route"
	"github.com/bartick/ringover-task/app/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	config *model.Configuration
	db     *sqlx.DB
	router *gin.Engine
	once   sync.Once
	log    *zap.Logger
)

func initTestEnvironment() {
	once.Do(func() {
		var err error
		log = utils.InitLogger()
		config, err = utils.LoadConfig()
		if err != nil {
			log.Fatal("Configuration loading failed", zap.String("err", err.Error()))
		}
		db, err = database.InitDatabases(config.Database)
		if err != nil {
			log.Fatal("Database initialization failed", zap.String("err", err.Error()))
		}
		router = route.AddAPIRouter(db)

		log.Info("Test environment initialized")
	})
}
