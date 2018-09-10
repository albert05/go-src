package main

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/service"
	"kd.explorer/util/dates"
	"os"
	"kd.explorer/exception"
	"kd.explorer/util/logger"
	"kd.explorer/service/transfer"
)

func main() {
	service.ConfigInit()

	if !common.Lock() {
		os.Exit(0)
	}
	defer exception.Handle(true)

	startTime := dates.NowTime()
	now := startTime

	for now-startTime < config.RunDURATION {
		// run analyse
		transfer.RunTA()

		dates.SleepSecond(config.SleepT)
		logger.Info(fmt.Sprintf("sleep %f second", config.SleepT))
		now = dates.NowTime()
	}
}
