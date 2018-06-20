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

const RunDURATION = 300
const LockCODE  = "RUN.MONITOR.EXCHANGE.TEST"
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
	workId := "exchange"
	currentDir := common.GetPwd()
	logPath := common.GetLogPath(workId)

	n := tool.NowTime()
	for n - startTime < RunDURATION {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id='%s' limit 10", status, workId)
		list, err := mysql.Conn.FindAll(sql)
		if err != nil {
			log.Fatal(err)
		}

		now := tool.NowDateStr()
		taskList := make([]string, 0)
		for _, task := range list {
			runTime := task.GetAttrString("run_time")
			if runTime <= now {
				taskList = append(taskList, task.GetAttrString("id"))
			}
		}

		if len(taskList) > 0 {
			ids := strings.Join(taskList, ",")
			mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id in (%s)", ids))
			cmdStr := fmt.Sprintf("cd %s;./%s -t %s -l %s %s", currentDir, TaskScriptName, workId, ids, logPath)
			common.Cmd(cmdStr)
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = tool.NowTime()
	}
}



