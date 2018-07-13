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

// 告警线
const MonitorMaxFEE = 100000  // 10万以下
const MonitorMinRATE = 15  // 15% 以上
const MonitorMaxRestDAY = 300    // 300天以内

// 秒杀线
const SecKillMaxFEE = 50000  // 5万以下
const SecKillMinRATE = 30  // 30% 以上
const SecKillMaxRestDAY = 150    // 60天以内

var MonitorRule *Rule
var SecKillRule *Rule

func init() {
	MonitorRule = InitRule()
	MonitorRule.SetFee(MonitorMaxFEE)
	MonitorRule.SetRate(MonitorMinRATE)
	MonitorRule.SetRestdays(MonitorMaxRestDAY)

	SecKillRule = InitRule()
	SecKillRule.SetFee(SecKillMaxFEE)
	SecKillRule.SetRate(SecKillMinRATE)
	SecKillRule.SetRestdays(SecKillMaxRestDAY)
}

func (list *TransList) Analyse() {
	monitorMsg := make([]string, 0)
	for _, item := range list.List.Items {
		if true == SecKillRule.Check(item) {
			if list.Cookie == "" {
				item.SyncRunKill()
			} else {
				item.RunKill(list.Cookie)
			}
		}
		if true == MonitorRule.Check(item) && !CheckIsSended(item.GetKey(), item.String()) {
			monitorMsg = append(monitorMsg, item.GetMonitorMsg())
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

