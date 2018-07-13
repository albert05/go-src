package main

import (
	"kd.explorer/config"
	"kd.explorer/tools/dates"
	"kd.explorer/service/kd"
	"flag"
	"time"
	"fmt"
	"os"
	"kd.explorer/common"
)

const LockTransferCODE = "RUN.MONITOR.TRANSFERS"

func main() {
	flag.StringVar(&config.CurUser, "u", "", "current user")
	flag.Parse()

	code := LockTransferCODE + config.CurUser
	if !common.Lock(code) {
		fmt.Println(code + " is running...")
		os.Exit(0)
	}
	defer func() {
		common.UnLock(code)
		os.Exit(0)
	}()

	startTime := dates.NowTime()
	now := startTime

	for now - startTime < config.RunDURATION {
		// run analyse
		kd.RunTA()

		time.Sleep(1 * time.Second)
		fmt.Println("sleep 1 second")
		now = dates.NowTime()
	}
}
