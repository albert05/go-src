package service

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"kd.explorer/config"
	"kd.explorer/tool"
	"kd.explorer/mysql"
)

const MonitorMaxFEE = 100000  // 10万以下
const MonitorMinRATE = 15  // 15% 以上
const MonitorMaxRestDAY = 300    // 300天以内

var CheckRule = map[string]func(item *TransferItem, v float64) bool {
	"fee-" + strconv.Itoa(MonitorMaxFEE): CheckFee,
	"rate-" + strconv.Itoa(MonitorMinRATE): CheckRate,
	"restdays-" + strconv.Itoa(MonitorMaxRestDAY): CheckRestDays,
}

func (list *TransList) Analyse() {
	monitorMsg := make([]string, 0)
	for _, item := range list.List.Items {
		if true == item.Check() && !CheckIsSended(item.GetKey()) {
			monitorMsg = append(monitorMsg, item.GetMonitorMsg())
		}
	}

	// is send monitor msg
	if len(monitorMsg) > 0 {
		msg := "高息转让项目提醒 >> " + strings.Join(monitorMsg, "@@")
		fmt.Println(msg)
		for _, phone := range config.SmsReceiverList {
			tool.Send(phone, msg)
		}
	}
}

func (item *TransferItem) GetMonitorMsg() string {
	return fmt.Sprintf("转让年化：%.2f%s, 金额：%.2f, 剩余天数：%d", item.GetRate(), "%", item.GetFee(), item.RestDays)
}

func (item *TransferItem) Check() bool {
	for e := range CheckRule {
		val := strings.Split(e, "-")
		v, _ := strconv.ParseFloat(val[1], 64)
		if false == CheckRule[e](item, v) {
			return false
		}
	}

	return true
}

func CheckFee(item *TransferItem, value float64) bool {
	if item.GetFee() > value {
		return false
	}
	return true
}

func CheckRate(item *TransferItem, value float64) bool {
	if item.GetRate() < value {
		return false
	}
	return true
}

func CheckRestDays(item *TransferItem, value float64) bool {
	if item.RestDays > int(value) {
		return false
	}
	return true
}

func CheckIsSended(transId string) bool {
	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM monitor WHERE trans_id = '%s'", transId))
	if err != nil {
		log.Fatal(err)
	}

	if len(userInfo) <= 0 {
		mysql.Conn.Exec(fmt.Sprintf("INSERT INTO monitor(trans_id, created_at) VALUES('%s', %d)", transId, tool.NowTime()))
		return false
	}

	return true
}

func RunTA() {
	list := RetryTransList()
	list.Analyse()
}

