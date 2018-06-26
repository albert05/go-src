package main

import (
	"fmt"
	"log"
	"kd.explorer/mysql"
	"kd.explorer/tool"
	"time"
	"kd.explorer/common"
	"os"
	"strings"
)

const RunDURATION = 290
const LockCODE  = "RUN:MONITOR:EXCHANGE"
const TaskScriptName = "run"

func main() {
	if !common.Lock(LockCODE) {
		fmt.Println(LockCODE + " is running...")
		os.Exit(0)
	}
	defer func() {
		common.UnLock(LockCODE)
		os.Exit(0)
	}()

	startTime := tool.NowTime()
	status := 0
	workId := `"exchange","abcGift"`
	currentDir := common.GetPwd()
	var logPath string

	n := tool.NowTime()
	for n - startTime < RunDURATION {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id in(%s) limit 10", status, workId)
		list, err := mysql.Conn.FindAll(sql)
		if err != nil {
			log.Fatal(err)
		}

		now := tool.NowDateStr()
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
				ids := strings.Join(list, ",")
				mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id in (%s)", ids))

				cmdStr := common.GetCmdStr(workId, map[string]string {"ids": ids, "curDir": currentDir, "logDir": logPath})
				common.Cmd(cmdStr)
			}
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = tool.NowTime()
	}
}
