package main

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/config"
	"kd.explorer/model"
	"kd.explorer/util/dates"
	"os"
	"strings"
	"time"
	"kd.explorer/service"
	"kd.explorer/exception"
)

func main() {
	service.ConfigInit()

	if !common.Lock() {
		os.Exit(0)
	}
	defer exception.Handle(true)

	startTime := dates.NowTime()
	status := 0
	workId := `"exchange"`
	currentDir := common.GetPwd()
	var logPath string

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
				logPath = common.GetLogPath(workId)
				model.UpdateMultiTask(list, map[string]string{
					"status": "1",
				})

				cmdStr := common.GetCmdStr(workId, map[string]string{"ids": strings.Join(list, ","), "curDir": currentDir, "logDir": logPath})
				common.Cmd(cmdStr)
			}
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = dates.NowTime()
	}
}
