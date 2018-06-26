//url := "https://enjoy.abchina.com/yh-web/customer/giftTokenDraw"
//params := `{"sessionId":"{ps_db86dd27589e11ea2b07b3018f10c8e8}_common","ruleNo":"064119","actNo":"999999CXE00064117","discType":"12","actType":"E","appr":"10"}`
//actNo := "999999CXE00064117"
//sessionID := "d46d835dab7daf16ccbc8bc27d5f995e"
package main

import (
	"kd.explorer/service"
	"fmt"
	"flag"
	"kd.explorer/common"
	"log"
	"kd.explorer/mysql"
	"kd.explorer/tool"
)

func main() {
	var jobList string
	flag.StringVar(&jobList, "l", "", "jobList")
	flag.Parse()

	sql := fmt.Sprintf("SELECT * FROM tasks WHERE id in (%s)", jobList)
	job, err := mysql.Conn.FindOne(sql)
	if err != nil {
		log.Fatal(err)
	}

	giftItem, err := service.GetGiftDetail(job.GetAttrString("product_id"))
	if err != nil {
		mysql.Conn.Exec(fmt.Sprintf("update tasks set status=2,result='%s' where id=%d", err.Error(), job.GetAttrInt("id")))
		log.Fatal(err)
	}

	common.Wait(job.GetAttrFloat("time_point"))

	giftItem.SetSession(job.GetAttrString("code"))

	i := 0
	for i < 3 {
		ret := giftItem.RunGift()

		status := 3
		if !ret {
			status = 2
		}
		mysql.Conn.Exec(fmt.Sprintf("update tasks set status=%d,result='run gift failed' where id=%d", status, job.GetAttrInt("id")))
		tool.SleepSecond(5)
		i++
	}
}
