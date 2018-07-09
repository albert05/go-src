package icbc

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"kd.explorer/tools/https"
	"bytes"
)

const TaskListURL = "http://icbc2.91laomi.cn:8888/aigou/index"
const SignDataURL = "https://elife.icbc.com.cn/OFSTCUST/lifeIndex/custInfoToAPI.action?actId=%s&t_k=%s"
const GiftURL = "http://icbc2.91laomi.cn:8888/usb/prize/awards"

const ResultSuccess = "0"
const RetryDataCNT = 3
const RetryGiftCNT = 10

type Gift struct {
	Cookie string
	ActID string
	Data string
	TaskNumber string
	SubTaskNumber string
}

type TaskResp struct {
	Data []TaskMonth ``
}

type TaskMonth struct {
	MainTaskNo string `json:"mainTaskNo"`
	List []TaskItem
}

type TaskItem struct {
	SubTaskNo string `json:"subTaskNo"`
	Id string `json:"id"`
}

func InitG(cookie, actId string) *Gift {
	gift := &Gift {
		Cookie: cookie,
		ActID: actId,
	}
	gift.InitTaskID()

	return gift
}

func (gift *Gift) InitTaskID() bool {
	resp, err := http.Get(TaskListURL)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result TaskResp
	json.Unmarshal(body, &result)

	if len(result.Data) <= 0 {
		log.Fatal(string(body))
	}

	return gift.match(result)
}

func (gift *Gift) match(task TaskResp) bool {
	for _, data := range task.Data {
		for _, item := range data.List {
			if gift.ActID == item.Id {
				gift.TaskNumber = data.MainTaskNo
				gift.SubTaskNumber = item.SubTaskNo
				return true
			}
		}
	}

	return false
}

// ******************** //

type SignData struct {
	Status string `json:"res"`
	Data string `json:"data"`
}

func (gift *Gift) refreshData() bool {
	url := fmt.Sprintf(SignDataURL, gift.ActID, gift.Cookie)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var result SignData
	json.Unmarshal(body, &result)

	if ResultSuccess != result.Status {
		fmt.Println(string(body))
		return false
	}

	gift.Data = result.Data
	return true
}

func (gift *Gift) RetryRData() bool {
	cnt := 0

	for cnt < RetryDataCNT {
		if gift.refreshData() {
			return true
		}
		cnt++
	}

	return false
}

// ******************** //

type GiftResp struct {
	IsSuccess bool `json:"is_success"`
	Msg string `json:"msg"`
}

func (gift *Gift) start() bool {
	params := map[string]string {
		"sub_task_number": gift.SubTaskNumber,
		"task_number": gift.TaskNumber,
		"cus_id": gift.Data,
	}
	bytesData, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
		return false
	}
	reader := bytes.NewReader(bytesData)
	body, err := https.PJson(GiftURL, reader)
	if err != nil {
		fmt.Println(err, string(body))
		return false
	}

	var result GiftResp
	json.Unmarshal(body, &result)

	if false == result.IsSuccess {
		fmt.Println(string(body))
	}

	return result.IsSuccess
}

func (gift *Gift) RetryStart() bool {
	cnt := 0
	for cnt < RetryGiftCNT {
		if gift.start() {
			return true
		}
		cnt++
	}

	return false
}

func (gift *Gift) RUN() bool {
	if !gift.RetryRData() {
		return false
	}

	return gift.RetryStart()
}
