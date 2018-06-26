package main

import (
	"fmt"
	"log"
	"kd.explorer/mysql"
	"kd.explorer/tool"
	"time"
	"kd.explorer/common"
	"os"
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
		//taskList := make([]string, 0)
		for _, task := range list {
			runTime := task.GetAttrString("run_time")
			if runTime <= now {
				//ids := strings.Join(taskList, ",")
				mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id in (%s)", task.GetAttrString("id")))

				logPath = common.GetLogPath(task.GetAttrString("work_id"))
				cmdStr := common.GetCmdStr(workId, map[string]string {"ids": task.GetAttrString("id"), "curDir": currentDir, "logDir": logPath})
				common.Cmd(cmdStr)
				//taskList = append(taskList, task.GetAttrString("id"))
			}
		}

		//if len(taskList) > 0 {
		//	ids := strings.Join(taskList, ",")
		//	mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id in (%s)", ids))
		//
		//	cmdStr := common.GetCmdStr(workId, map[string]string {"ids": ids, "curDir": currentDir, "logDir": logPath})
		//	common.Cmd(cmdStr)
		//}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = tool.NowTime()
	}
}
