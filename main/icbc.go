package main

import (
	"kd.explorer/service/icbc"
	"flag"
	"fmt"
	"log"
	"kd.explorer/tools/dates"
	"kd.explorer/tools/mysql"
	"kd.explorer/common"
)

func main()  {
	var jobList string
	flag.StringVar(&jobList, "l", "", "jobList")
	flag.Parse()

	sql := fmt.Sprintf("SELECT * FROM tasks WHERE id in (%s)", jobList)
	job, err := mysql.Conn.FindOne(sql)
	if err != nil {
		log.Fatal(err)
	}

	actId := job.GetAttrString("product_id")
	cookie := job.GetAttrString("code")

	common.Wait(job.GetAttrFloat("time_point"))
	gift := icbc.InitG(cookie, actId)

	i := 0
	for i < 100 {
		result := gift.RUN()
		if result {
			mysql.Conn.Exec(fmt.Sprintf("update tasks set status=%d,result='%s' where id=%d", 3, "success", job.GetAttrInt("id")))
			break
		}

		dates.SleepSecond(5)
		i++
	}

	mysql.Conn.Exec(fmt.Sprintf("update tasks set status=%d,result='%s' where id=%d", 2, "failed", job.GetAttrInt("id")))
}
