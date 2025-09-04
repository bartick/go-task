package main

import (
	"github.com/bartick/ringover-task/app/model"
	"github.com/bartick/ringover-task/app/shared/utils"
	"go.uber.org/zap"
)

var (
	log    *zap.Logger
	config *model.Configuration
)

func init() {
	log = utils.InitLogger()

	var err error
	config, err = utils.LoadConfig()
	if err != nil {
		log.Fatal("Configuration loading failed", zap.String("err", err.Error()))
	}
}

func main() {
	log.Debug("Start")
	runTaskAPI()
	log.Debug("End")
}
