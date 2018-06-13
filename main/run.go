package main

import (
	"fmt"
	"log"
	"kd.explorer/mysql"
	"kd.explorer/service"
	"kd.explorer/tool"
	"kd.explorer/model"
	"time"
	"kd.explorer/common"
	"os"
)

const RUN_DURATION = 300
const LOCK_CODE  = "RUN.EXCHANGE.TEST"

func main() {
	if !common.Lock(LOCK_CODE) {
		fmt.Println(LOCK_CODE + " is running...")
		os.Exit(0)
	}
	defer common.UnLock(LOCK_CODE)

	startTime := tool.NowTime()
	status := 0
	workId := "exchange"

	n := tool.NowTime()
	for n - startTime < RUN_DURATION {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id='%s' limit 10", status, workId)
		fmt.Println(mysql.Conn)
		list, err := mysql.Conn.FindAll(sql)
		if err != nil {
			log.Fatal(err)
		}

		now := tool.NowDateStr()
		runTaskList := make([]model.MapModel, 0)
		for _, task := range list {
			runTime := task.GetAttrString("run_time")
			if runTime <= now {
				runTaskList = append(runTaskList, task)
			}
		}

		if len(runTaskList) > 0 {
			service.GoRunTask(runTaskList)
		}

		time.Sleep(5 * time.Second)
		fmt.Println("sleep 5 second")
		n = tool.NowTime()
	}
}

