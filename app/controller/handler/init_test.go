package handler_test

import (
	"sync"

	"github.com/bartick/go-task/app/model"
	"github.com/bartick/go-task/app/route"
	"github.com/bartick/go-task/app/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	config *model.Configuration
	db     model.DBTX
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
		db, err = model.InitDatabases(config.Database)
		if err != nil {
			log.Fatal("Database initialization failed", zap.String("err", err.Error()))
		}
		router = route.AddAPIRouter(db.(*sqlx.DB))

		log.Info("Test environment initialized")
	})
}
