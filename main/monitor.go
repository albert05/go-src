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
	defer common.UnLock(LockCODE)

	startTime := tool.NowTime()
	status := 0
	workId := "exchange"
	currentDir := common.GetPwd()

	n := tool.NowTime()
	for n - startTime < RunDURATION {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id='%s' limit 10", status, workId)
		fmt.Println(mysql.Conn)
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
			cmdStr := fmt.Sprintf("cd %s;./%s -t %s -l %s", currentDir, TaskScriptName, "exchange", ids)
			common.Cmd(cmdStr)
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = tool.NowTime()
	}
}

