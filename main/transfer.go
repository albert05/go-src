package main

import (
	"kd.explorer/config"
	"kd.explorer/service/kd"
	"flag"
	"time"
	"fmt"
	"strconv"
)

func main() {
	var t string
	flag.StringVar(&config.CurUser, "u", "zhoushan_5781", "current user")
	flag.StringVar(&t, "t", "", "sleep time")
	flag.Float64Var(&config.SecKillFee, "fee", kd.SecKillMaxFEE, "")
	flag.Float64Var(&config.SecKillRate, "rate", kd.SecKillMinRATE, "")
	flag.IntVar(&config.SecKillRestDay, "rest", kd.SecKillMaxRestDAY, "")
	flag.Parse()

	st := 1
	if t != "" {
		st, _ = strconv.Atoi(t)
	}

	for {
		// run analyse
		kd.RunTA()

		time.Sleep(time.Duration(st) * time.Second)
		fmt.Println(fmt.Sprintf("sleep %d second", st))
	}
}
