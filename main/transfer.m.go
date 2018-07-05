package main

import (
	"kd.explorer/service"
	"kd.explorer/tool"
	"time"
	"fmt"
	"kd.explorer/config"
)

func main() {
	startTime := tool.NowTime()
	now := startTime

	for now - startTime < config.RunDURATION {
		// run analyse
		service.RunTA()

		time.Sleep(10 * time.Second)
		fmt.Println("sleep 10 second")
		now = tool.NowTime()
	}
}
