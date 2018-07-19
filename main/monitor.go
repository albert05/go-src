package main

import (
	"fmt"
	"time"
	"kd.explorer/common"
	"os"
	"kd.explorer/config"
	"kd.explorer/tools/dates"
	"kd.explorer/model"
	"strings"
)

const LockCODE  = "RUN:MONITOR:EXCHANGE"

func main() {
	if !common.Lock(LockCODE) {
		fmt.Println(LockCODE + " is running...")
		os.Exit(0)
	}
	defer func() {
		common.UnLock(LockCODE)
		os.Exit(0)
	}()

	startTime := dates.NowTime()
	status := 0
	workId := `"exchange"`
	currentDir := common.GetPwd()
	var logPath string

	n := dates.NowTime()
	for n - startTime < config.RunDURATION {
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
				model.UpdateMultiTask(list, map[string]string {
					"status": "1",
				})

				cmdStr := common.GetCmdStr(workId, map[string]string {"ids": strings.Join(list, ","), "curDir": currentDir, "logDir": logPath})
				common.Cmd(cmdStr)
			}
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = dates.NowTime()
	}
}
