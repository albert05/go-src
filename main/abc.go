package main

import (
	"kd.explorer/service"
	"fmt"
	"flag"
	"kd.explorer/common"
)

func main() {
	var actNo string
	var sessionID string
	var startTime float64
	flag.StringVar(&actNo, "a", "", "actNo")
	flag.StringVar(&sessionID, "s", "", "sessionID")
	flag.Float64Var(&startTime, "t", 100000, "startTime")
	flag.Parse()

	//url := "https://enjoy.abchina.com/yh-web/customer/giftTokenDraw"
	//params := `{"sessionId":"{ps_db86dd27589e11ea2b07b3018f10c8e8}_common","ruleNo":"064119","actNo":"999999CXE00064117","discType":"12","actType":"E","appr":"10"}`
	//actNo := "999999CXE00064117"
	//sessionID := "d46d835dab7daf16ccbc8bc27d5f995e"

	giftItem, err := service.GetGiftDetail(actNo)
	if err != nil {
		fmt.Println(err)
	}

	common.Wait(startTime)

	giftItem.SetSession(sessionID)
	ret := giftItem.RunGift()
	fmt.Println(ret)
}

