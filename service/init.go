package service

import (
	"flag"
	"kd.explorer/config"
	"kd.explorer/common"
	"kd.explorer/util/logger"
)

func ConfigInit() {
	// transfer.m.go
	flag.StringVar(&config.CurUser, "u", "", "")
	flag.StringVar(&config.SecUser, "su", "", "")
	flag.Float64Var(&config.SleepT, "t", 1, "sleep time")
	flag.Float64Var(&config.SecKillFee, "fee", 50000, "")
	flag.Float64Var(&config.SecKillRate, "rate", 30, "")
	flag.IntVar(&config.SecKillRestDay, "rest", 150, "")
	flag.StringVar(&config.RuleKey, "rkey", "", "")
	flag.Float64Var(&config.SecKillTime, "st", 3, "")

	// exchange.go
	flag.StringVar(&config.JobType, "tp", "exchange", "jobType")
	flag.StringVar(&config.JobList, "l", "", "jobList")
	flag.Parse()

	// 获取本机IP
	localIp, err := common.GetLocalIp()
	if err != nil {
		panic("GetLocalIp Err:" + err.Error())
	}
	config.LocalIp = localIp

	// log init
	logger.SetConsole(false)
	logger.SetLevel(logger.INFO)
	logger.SetRollingDaily(common.GetLogPath())
}
