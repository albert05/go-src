package main

import (
	"fmt"
	"kd.explorer/tools/mysql"
	"kd.explorer/model"
	"flag"
	"kd.explorer/service"
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

