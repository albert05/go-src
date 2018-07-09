package main

import (
	"fmt"
	"log"
	"kd.explorer/tools/mysql"
	"time"
	"kd.explorer/common"
	"os"
	"strings"
	"kd.explorer/config"
	"kd.explorer/tools/dates"
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
	workId := `"exchange","abcGift"`
	currentDir := common.GetPwd()
	var logPath string

	n := dates.NowTime()
	for n - startTime < config.RunDURATION {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id in(%s) limit 10", status, workId)
		list, err := mysql.Conn.FindAll(sql)
		if err != nil {
			log.Fatal(err)
		}

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
				ids := strings.Join(list, ",")
				mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id in (%s)", ids))

				cmdStr := common.GetCmdStr(workId, map[string]string {"ids": ids, "curDir": currentDir, "logDir": logPath})
				common.Cmd(cmdStr)
			}
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = dates.NowTime()
	}
}
