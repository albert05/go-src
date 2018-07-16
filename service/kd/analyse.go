package kd

import (
	"fmt"
	"strings"
	"kd.explorer/config"
	"kd.explorer/tools/mail"
	"kd.explorer/model"
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
	SecKillRule.SetFee(config.SecKillFee)
	SecKillRule.SetRate(config.SecKillRate)
	SecKillRule.SetRestdays(config.SecKillRestDay)
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
		email := model.FindUser(config.CurUser).GetAttrString("email")
		mail.SendSingle(email, "高息转让项目提醒", msg)
	}
}

func (item *TransferItem) GetMonitorMsg() string {
	return fmt.Sprintf("转让年化：%.2f%s, 金额：%.2f, 剩余天数：%d", item.GetRate(), "%", item.GetFee(), item.RestDays)
}

func CheckIsSended(transId string, data string) bool {
	monitorInfo := model.FindMRecord(transId)

	if len(monitorInfo) <= 0 {
		model.InsertMRecord(transId, data)
		return false
	}

	return true
}

func RunTA() {
	list := RetryTransList()
	if list != nil {
		fmt.Println(list)
		list.Analyse()
	}
}

