//url := "https://enjoy.abchina.com/yh-web/customer/giftTokenDraw"
//params := `{"sessionId":"{ps_db86dd27589e11ea2b07b3018f10c8e8}_common","ruleNo":"064119","actNo":"999999CXE00064117","discType":"12","actType":"E","appr":"10"}`
//actNo := "999999CXE00064117"
//sessionID := "d46d835dab7daf16ccbc8bc27d5f995e"
package main

import (
	"fmt"
	"flag"
	"kd.explorer/common"
	"log"
	"kd.explorer/tools/dates"
	"kd.explorer/service/abc"
	"kd.explorer/model"
	"strconv"
)

func main() {
	var ids string
	flag.StringVar(&ids, "l", "", "job")
	flag.Parse()

	id ,_ := strconv.Atoi(ids)
	job := model.FindTask(id)

	giftItem, err := abc.GetGiftDetail(job.GetAttrString("product_id"))
	if err != nil {
		model.UpdateTask(job.GetAttrInt("id"), map[string]string {
			"status": "2",
			"result": err.Error(),
		})
		log.Fatal(err)
	}

	common.Wait(job.GetAttrFloat("time_point"))

	giftItem.SetSession(job.GetAttrString("code"))

	i := 0
	for i < 3 {
		giftRep := giftItem.RunGift()

		status := 3
		if abc.GiftStatusSUCCESS != giftRep.Status {
			status = 2
		}

		model.UpdateTask(job.GetAttrInt("id"), map[string]string {
			"status": fmt.Sprintf("%d", status),
			"result": giftRep.Result,
		})
		dates.SleepSecond(5)
		i++
	}
}
