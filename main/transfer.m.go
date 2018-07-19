package main

import (
	"kd.explorer/config"
	"kd.explorer/tools/dates"
	"kd.explorer/service"
	"flag"
	"time"
	"fmt"
	"os"
	"kd.explorer/common"
	"strconv"
)

const LockTransferCODE = "RUN.MONITOR.TRANSFERS"

func main() {
	var t string
	flag.StringVar(&config.CurUser, "u", "", "current user")
	flag.StringVar(&t, "t", "", "sleep time")
	flag.Float64Var(&config.SecKillFee, "fee", service.SecKillMaxFEE, "")
	flag.Float64Var(&config.SecKillRate, "rate", service.SecKillMinRATE, "")
	flag.IntVar(&config.SecKillRestDay, "rest", service.SecKillMaxRestDAY, "")
	flag.StringVar(&config.RuleKey, "rkey", "", "")
	flag.Int64Var(&config.SecKillTime, "st", 4, "")
	flag.Parse()

	st := 1
	if t != "" {
		st, _ = strconv.Atoi(t)
	}

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

		time.Sleep(time.Duration(st) * time.Second)
		fmt.Println(fmt.Sprintf("sleep %d second", st))
		now = dates.NowTime()
	}
}
