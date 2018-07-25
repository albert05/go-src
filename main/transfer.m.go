package main

import (
	"kd.explorer/config"
	"kd.explorer/util/dates"
	"kd.explorer/service"
	"flag"
	"fmt"
	"os"
	"kd.explorer/common"
)

const LockTransferCODE = "RUN.MONITOR.TRANSFERS"

func main() {
	var t float64
	flag.StringVar(&config.CurUser, "u", "", "current user")
	flag.Float64Var(&t, "t", 1, "sleep time")
	flag.Float64Var(&config.SecKillFee, "fee", service.SecKillMaxFEE, "")
	flag.Float64Var(&config.SecKillRate, "rate", service.SecKillMinRATE, "")
	flag.IntVar(&config.SecKillRestDay, "rest", service.SecKillMaxRestDAY, "")
	flag.StringVar(&config.RuleKey, "rkey", "", "")
	flag.Float64Var(&config.SecKillTime, "st", 3, "")
	flag.Parse()

	code := fmt.Sprintf(LockTransferCODE + "_%s_%f_%f_%d", config.CurUser, config.SecKillFee, config.SecKillRate, config.SecKillRestDay)
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
		service.RunTA()

		dates.SleepSecond(t)
		fmt.Println(fmt.Sprintf("sleep %f second", t))
		now = dates.NowTime()
	}
}
