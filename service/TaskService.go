package service

import (
	"kd.explorer/mysql"
	"fmt"
	"kd.explorer/tool"
	"time"
	"kd.explorer/model"
	"encoding/json"
)

type TaskResponse struct {
	Code int ``
	Message string ``
}

const DEFAULT_SLEEP_TIME = time.Millisecond * 10
const EXCHANGE_URL = "https://deposit.koudailc.com/user-order-form/convert"

func GoRunTask(taskList []model.MapModel) {
	ch := make(chan string)
	for _, task := range taskList {
		go runT(task, ch)
	}

	for range taskList {
		<-ch
	}

	close(ch)
}

func runT(task model.MapModel, ch chan<- string) {
	taskId := task.GetAttrInt("id")
	fmt.Println(fmt.Sprintf("taskID %d start work", taskId))
	mysql.Conn.Exec(fmt.Sprintf("update tasks set status=1 where id=%d", taskId))
	userKey := task.GetAttrString("user_key")

	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM userinfos WHERE user_key = '%s'", userKey))
	if err != nil || userInfo == nil {
		ch <- "get user info failed"
		return
	}

	username := userInfo.GetAttrString("user_name")
	password := userInfo.GetAttrString("password")
	cookie, err := Login(username, password)
	if err != nil {
		ch <- "login failed"
		return
	}

	timePoint := task.GetAttrFloat("time_point")

	wait(timePoint)

	pid := task.GetAttrString("product_id")
	imgCode := task.GetAttrString("code")
	prizeNumber := task.GetAttrString("prize_number")

	params := map[string]string{
		"id": pid,
		"imgcode": imgCode,
		"prize_number": prizeNumber,
	}

	body, err := tool.Post(EXCHANGE_URL, params, cookie)
	if err != nil {
		fmt.Println(err)
		ch <- err.Error()
		return
	}

	var result TaskResponse
	json.Unmarshal(body, &result)

	status := 3
	msg := ""
	if tool.HTTP_SUCCESS != result.Code {
		status = 2
		msg = result.Message
	}
	mysql.Conn.Exec(fmt.Sprintf("update tasks set status=%d,result='%s' where id=%d", status, msg, taskId))

	fmt.Println(userKey + " -- " + string(body))
	ch <- "success"
}

func wait(timePoint float64) {
	currTime := tool.TimeInt2float(tool.CurrentMicro())
	fmt.Println(currTime, timePoint)

	for currTime < timePoint {
		time.Sleep(DEFAULT_SLEEP_TIME)

		currTime = tool.TimeInt2float(tool.CurrentMicro())
	}

	return
}