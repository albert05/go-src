package main

import (
	"fmt"
	"log"
	"hello/mysql"
	"hello/service"
	"hello/tool"
	"hello/model"
	"time"
)

const RUN_DURATION = 300

func main() {

	startTime := tool.NowTime()
	status := 0
	workId := "exchange"

	for n := tool.NowTime(); n - startTime < RUN_DURATION; {
		sql := fmt.Sprintf("SELECT * FROM tasks WHERE status =%d and work_id='%s' limit 10", status, workId)
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
	}
}

