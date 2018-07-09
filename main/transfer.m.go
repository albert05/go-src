package main

import (
	"time"
	"fmt"
	"kd.explorer/config"
	"os"
	"kd.explorer/common"
	"kd.explorer/tools/dates"
	"kd.explorer/service/kd"
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


	startTime := dates.NowTime()
	now := startTime

	for now - startTime < config.RunDURATION {
		// run analyse
		kd.RunTA()

		time.Sleep(10 * time.Second)
		fmt.Println("sleep 10 second")
		now = dates.NowTime()
	}
}
