package main

import (
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/model"
	"kd.explorer/util/dates"
	"os"
	"time"
	"kd.explorer/service"
	"kd.explorer/exception"
	"kd.explorer/util/logger"
	"kd.explorer/service/tasks"
)

func main() {
	service.ConfigInit()

	if !common.Lock() {
		os.Exit(0)
	}
	defer exception.Handle(true)

	startTime := dates.NowTime()
	status := 0
	workId := `"exchange","daily"`
	currentDir := common.GetPwd()

	n := dates.NowTime()
	for n-startTime < config.RunDURATION {
		list := model.FindTaskListByStatus(status, workId)

		now := dates.NowDateStr()
		taskList := make(map[string][]string)
		for _, task := range list {
			runTime := task.GetAttrString("run_time")
			if runTime <= now {
				workId := task.GetAttrString("work_id")
				if len(taskList[workId]) <= 0 {
					taskList[workId] = make([]string, 0)
				}
				taskList[workId] = append(taskList[workId], task.GetAttrString("id"))
			}
		}

		if len(taskList) > 0 {
			for workId, list := range taskList {
				model.UpdateMultiTask(list, tasks.GetUpdateData(workId))

				for _,id := range list {
					cmdStr := common.GetCmdStr(workId, map[string]string{"ids": id, "curDir": currentDir})
					common.Cmd(cmdStr)
				}
			}
		}

		logger.Info("sleep 5 second")
		time.Sleep(5 * time.Second)
		n = dates.NowTime()
	}
}
