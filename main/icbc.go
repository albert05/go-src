package main

import (
	"kd.explorer/service/icbc"
	"flag"
	"kd.explorer/tools/dates"
	"kd.explorer/common"
	"strconv"
	"kd.explorer/model"
)

func main()  {
	var ids string
	flag.StringVar(&ids, "l", "", "job")
	flag.Parse()

	id ,_ := strconv.Atoi(ids)
	job := model.FindTask(id)

	actId := job.GetAttrString("product_id")
	cookie := job.GetAttrString("code")

	common.Wait(job.GetAttrFloat("time_point"))
	gift := icbc.InitG(cookie, actId)

	i := 0
	for i < 100 {
		result := gift.RUN()
		if result {
			model.UpdateTask(job.GetAttrInt("id"), map[string]string {
				"status": "3",
				"result": "success",
			})
			break
		}

		dates.SleepSecond(5)
		i++
	}

	model.UpdateTask(job.GetAttrInt("id"), map[string]string {
		"status": "2",
		"result": "failed",
	})
}
