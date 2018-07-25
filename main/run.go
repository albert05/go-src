package main

import (
	"flag"
	"fmt"
	"kd.explorer/model"
	"kd.explorer/service"
	"kd.explorer/util/mysql"
)

func main() {
	var jobType string
	var jobList string
	flag.StringVar(&jobType, "t", "exchange", "jobType")
	flag.StringVar(&jobList, "l", "", "jobList")
	flag.Parse()

	fmt.Println(jobList)
	list := model.FindTaskListByIds(jobList)

	runTaskList := make([]mysql.MapModel, 0)
	for _, task := range list {
		runTaskList = append(runTaskList, task)
	}

	service.GoRunTask(runTaskList)
}
