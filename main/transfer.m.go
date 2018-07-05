package main

import (
	"kd.explorer/service"
	"kd.explorer/tool"
	"time"
	"fmt"
	"kd.explorer/config"
	"os"
	"kd.explorer/common"
)

const LockTransferCODE = "RUN.MONITOR.TRANSFER"

func main() {
	if !common.Lock(LockTransferCODE) {
		fmt.Println(LockTransferCODE + " is running...")
		os.Exit(0)
	}
	defer func() {
		common.UnLock(LockTransferCODE)
		os.Exit(0)
	}()


	startTime := tool.NowTime()
	now := startTime

	for now - startTime < config.RunDURATION {
		// run analyse
		service.RunTA()

		time.Sleep(10 * time.Second)
		fmt.Println("sleep 10 second")
		now = tool.NowTime()
	}
}
