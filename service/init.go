package service

import (
	"flag"
	"kd.explorer/config"
	"kd.explorer/common"
)

func ConfigInit() {
	// transfer.m.go
	flag.StringVar(&config.CurUser, "u", "", "")
	flag.Float64Var(&config.SleepT, "t", 1, "sleep time")
	flag.Float64Var(&config.SecKillFee, "fee", 50000, "")
	flag.Float64Var(&config.SecKillRate, "rate", 30, "")
	flag.IntVar(&config.SecKillRestDay, "rest", 150, "")
	flag.StringVar(&config.RuleKey, "rkey", "", "")
	flag.Float64Var(&config.SecKillTime, "st", 3, "")

	// run.go
	flag.StringVar(&config.JobType, "tp", "exchange", "jobType")
	flag.StringVar(&config.JobList, "l", "", "jobList")
	flag.Parse()

	// 获取本机IP
	localIp, err := common.GetLocalIp()
	if err != nil {
		panic("GetLocalIp Err:" + err.Error())
	}
	config.LocalIp = localIp
}
