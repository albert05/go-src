package main

import (
	"fmt"
	"log"
	"kd.explorer/mysql"
	"kd.explorer/service"
	"kd.explorer/model"
	"flag"
)

func main() {
	var jobType string
	var jobList string
	flag.StringVar(&jobType, "t", "exchange", "jobType")
	flag.StringVar(&jobList, "l", "", "jobList")
	flag.Parse()

	fmt.Println(jobList)

	sql := fmt.Sprintf("SELECT * FROM tasks WHERE id in (%s)", jobList)
	fmt.Println(mysql.Conn)
	list, err := mysql.Conn.FindAll(sql)
	if err != nil {
		log.Fatal(err)
	}

	runTaskList := make([]model.MapModel, 0)
	for _, task := range list {
		runTaskList = append(runTaskList, task)
	}

	service.GoRunTask(runTaskList)
}

