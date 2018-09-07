package main

import (
	"kd.explorer/model"
	"kd.explorer/service"
	"kd.explorer/util/mysql"
	"kd.explorer/config"
)

func main() {
	service.ConfigInit()

	list := model.FindTaskListByIds(config.JobList)

	runTaskList := make([]mysql.MapModel, 0)
	for _, task := range list {
		runTaskList = append(runTaskList, task)
	}

	service.GoRunTask(runTaskList)
}
