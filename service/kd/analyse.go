package kd

import (
	"fmt"
	"log"
	"strings"
	"kd.explorer/config"
	"kd.explorer/tools/mysql"
	"kd.explorer/tools/mail"
	"kd.explorer/tools/dates"
)

const MonitorMaxFEE = 100000  // 10万以下
const MonitorMinRATE = 15  // 15% 以上
const MonitorMaxRestDAY = 300    // 300天以内

const SecKillMaxFEE = 50000  // 5万以下
const SecKillMinRATE = 30  // 30% 以上
const SecKillMaxRestDAY = 90    // 60天以内

var CheckRule = map[string]func(item *TransferItem, v float64) bool {
	"fee": CheckFee,
	"rate": CheckRate,
	"restdays": CheckRestDays,
}

// 告警线
var MonitorLine = map[string]float64 {
	"fee": MonitorMaxFEE,
	"rate": MonitorMinRATE,
	"restdays": MonitorMaxRestDAY,
}

// 秒杀线
var SecKillLine = map[string]float64 {
	"fee": SecKillMaxFEE,
	"rate": SecKillMinRATE,
	"restdays": SecKillMaxRestDAY,
}

func (list *TransList) Analyse() {
	monitorMsg := make([]string, 0)
	for _, item := range list.List.Items {
		if !CheckIsSended(item.GetKey(), item.String()) {
			if true == item.Check(MonitorLine) {
				monitorMsg = append(monitorMsg, item.GetMonitorMsg())
			}
			if true == item.Check(SecKillLine) {
				item.RunKILL()
			}
		}
	}

	fmt.Println(monitorMsg)

	// is send monitor msg
	if len(monitorMsg) > 0 {
		msg := "高息转让项目提醒 >> " + strings.Join(monitorMsg, "@@")
		fmt.Println(msg)
		// send mail
		for _, receiver := range config.MailReceiverList {
			mail.SendSingle(receiver, "高息转让项目提醒", msg)
		}
		//ret := mail.Send(config.MailReceiverList, "高息转让项目提醒", msg)

		//send sms
		//for _, phone := range config.SmsReceiverList {
		//	tools.Send(phone, msg)
		//}
	}
}

func (item *TransferItem) GetMonitorMsg() string {
	return fmt.Sprintf("转让年化：%.2f%s, 金额：%.2f, 剩余天数：%d", item.GetRate(), "%", item.GetFee(), item.RestDays)
}

func (item *TransferItem) Check(params map[string]float64) bool {
	for e := range CheckRule {
		if false == CheckRule[e](item, params[e]) {
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

func CheckIsSended(transId string, data string) bool {
	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM trans_monitor_list WHERE trans_id = '%s'", transId))
	if err != nil {
		log.Fatal(err)
	}

	if len(userInfo) <= 0 {
		mysql.Conn.Exec(fmt.Sprintf("INSERT INTO trans_monitor_list(trans_id, created_at, data) VALUES('%s', %d, '%s')", transId, dates.NowTime(), data))
		return false
	}

	return true
}

func RunTA() {
	list := RetryTransList()
	fmt.Println(list)
	list.Analyse()
}

