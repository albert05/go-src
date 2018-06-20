package service

import (
	"kd.explorer/mysql"
	"fmt"
	"kd.explorer/tool"
	"time"
	"kd.explorer/model"
	"encoding/json"
	"log"
	"strings"
)

type TaskResponse struct {
	Code int ``
	Message string ``
}

const DefaultSleepTIME = time.Millisecond * 10
const ExchangeURL = "https://deposit.koudailc.com/user-order-form/convert"

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

	var code Code
	code.setCookie(cookie)
	code.Refresh()
	code.RandFileName()
	code.CreateImage()

	fmt.Println(cookie, code.getFileName())

	mysql.Conn.Exec(fmt.Sprintf("update tasks set img_url='%s' where id=%d", code.getFileName(), taskId))

	timePoint := task.GetAttrFloat("time_point")
	imgCode := wait(timePoint, taskId)

	pid := task.GetAttrString("product_id")
	prizeNumber := task.GetAttrString("prize_number")

	params := map[string]string{
		"id": pid,
		"imgcode": imgCode,
		"prize_number": prizeNumber,
	}

	body, err := tool.Post(ExchangeURL, params, cookie)
	if err != nil {
		fmt.Println(err)
		ch <- err.Error()
		return
	}

	var result TaskResponse
	json.Unmarshal(body, &result)

	status := 3
	msg := ""
	if tool.HttpSUCCESS != result.Code {
		status = 2
		msg = result.Message
	}
	mysql.Conn.Exec(fmt.Sprintf("update tasks set status=%d,result='%s' where id=%d", status, msg, taskId))

	fmt.Println(userKey + " -- " + string(body))
	ch <- "success"
}

func wait(timePoint float64, taskId int) string {
	currTime := tool.TimeInt2float(tool.CurrentMicro())
	fmt.Println(currTime, timePoint)

	var imgCode string

	for currTime < timePoint {
		time.Sleep(DefaultSleepTIME)

		if imgCode == "" {
			sql := fmt.Sprintf("SELECT * FROM tasks WHERE id =%d", taskId)
			task, err := mysql.Conn.FindOne(sql)
			if err != nil {
				log.Fatal(err)
			}

			imgCode = strings.Trim(task.GetAttrString("code"), " ")
		}

		currTime = tool.TimeInt2float(tool.CurrentMicro())
	}

	return imgCode
}