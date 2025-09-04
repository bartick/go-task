package main

import (
	"github.com/bartick/ringover-task/app/shared/utils"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log = utils.InitLogger()

	LoadConfig()
}

func main() {
	log.Debug("Start")
	runTaskAPI()
	log.Debug("End")
}
