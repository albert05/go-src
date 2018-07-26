package main

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/service"
	"kd.explorer/util/dates"
	"os"
	"kd.explorer/exception"
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
		service.RunTA()

		dates.SleepSecond(config.SleepT)
		fmt.Println(fmt.Sprintf("sleep %f second", config.SleepT))
		now = dates.NowTime()
	}
}
