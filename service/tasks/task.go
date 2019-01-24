package tasks

import (
	"fmt"
	"kd.explorer/common"
	"kd.explorer/model"
	"kd.explorer/util/dates"
	"kd.explorer/util/mysql"
	"strconv"
	"strings"
	"time"
	"kd.explorer/util/logger"
	"kd.explorer/service/base"
)

type TaskResponse struct {
	Code    int    ``
	Message string ``
}

func GoRunTask(taskList []mysql.MapModel) {
	ch := make(chan string)
	for _, task := range taskList {
		switch task.GetAttrString("work_id") {
			case "exchange":
				go runExchange(task, ch)
				break
			case "daily":
				go runDaily(task, ch)
				break
			default:
				continue
		}
	}

	for range taskList {
		<-ch
	}

	close(ch)
}

func runExchange(task mysql.MapModel, ch chan<- string) {
	taskId := task.GetAttrInt("id")
	logger.Info(fmt.Sprintf("taskID %d start work", taskId))
	userKey := task.GetAttrString("user_key")

	cookie, err := base.LoginK(userKey)
	if err != nil {
		ch <- "login failed"
		return
	}

	var code base.Code
	code.SetCookie(cookie)
	code.Refresh()
	code.RandFileName()
	code.CreateImage()

	logger.Info(cookie, code.GetFileName())

	model.UpdateTask(taskId, map[string]string{
		"img_url": code.GetFileName(),
	})

	timePoint := task.GetAttrFloat("time_point")
	imgCode := wait(timePoint, taskId)

	pid := task.GetAttrString("product_id")
	prizeNumber := task.GetAttrString("prize_number")

	params := map[string]string{
		"id":           pid,
		"imgcode":      imgCode,
		"prize_number": prizeNumber,
	}

	isOk, errMsg := Exchange(cookie, params)
	status := 3
	if isOk {
		status = 2
	}

	model.UpdateTask(taskId, map[string]string{
		"status": strconv.Itoa(status),
		"result": errMsg,
	})

	ch <- "success"
}

func runDaily(task mysql.MapModel, ch chan<- string) {
	taskId := task.GetAttrInt("id")
	logger.Info(fmt.Sprintf("taskID %d start work", taskId))
	userKey := task.GetAttrString("user_key")

	cookie, err := base.LoginK(userKey)
	if err != nil {
		ch <- "login failed"
		return
	}

	//params := map[string]string{
	//	"type":           "3",  // 冒险性
	//}
	_, errMsg := Earn(cookie, nil)
	logger.Info(fmt.Sprintf("taskID %d earn result: %s", taskId, errMsg))

	_, errMsg = Share(cookie, nil)
	logger.Info(fmt.Sprintf("taskID %d shake result: %s", taskId, errMsg))

	ch <- "success"
}

func GetUpdateData(workId string) map[string]string {

	if workId == "exchange" {
		return map[string]string{
			"status": "1",
		}
	} else if workId == "daily" {
		return map[string]string{
			"run_time": dates.DateAfterNDays(1),
		}
	}

	return nil
}

func wait(timePoint float64, taskId int) string {
	currTime := dates.TimeInt2float(dates.CurrentMicro())
	logger.Info(currTime, timePoint)

	var imgCode string

	for currTime < timePoint {
		time.Sleep(common.DefaultSleepTIME)

		if imgCode == "" {
			task := model.FindTask(taskId)

			imgCode = strings.Trim(task.GetAttrString("code"), " ")
		}

		currTime = dates.TimeInt2float(dates.CurrentMicro())
	}

	return imgCode
}
