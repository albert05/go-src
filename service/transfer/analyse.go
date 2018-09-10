package transfer

import (
	"encoding/json"
	"kd.explorer/config"
	"kd.explorer/model"
	"kd.explorer/util/mysql"
	"kd.explorer/util/logger"
)

// 告警线
const MonitorMaxFEE = 100000  // 10万以下
const MonitorMinRATE = 15     // 15% 以上
const MonitorMaxRestDAY = 300 // 300天以内

var MonitorRule *Rule
var SecKillRules *Rules

func Init() {
	MonitorRule = InitRule()
	MonitorRule.SetFee(MonitorMaxFEE)
	MonitorRule.SetRate(MonitorMinRATE)
	MonitorRule.SetRestdays(MonitorMaxRestDAY)

	SecKillRules = &Rules{}

	var rule mysql.MapModel
	if config.RuleKey != "" {
		rule = model.FindRule(config.RuleKey)
	}

	if rule == nil {
		secKillRule := InitRule()
		secKillRule.SetFee(config.SecKillFee)
		secKillRule.SetRate(config.SecKillRate)
		secKillRule.SetRestdays(config.SecKillRestDay)
		SecKillRules.R = []*Rule{secKillRule}
	} else {
		var r []*Rule
		json.Unmarshal([]byte(rule.GetAttrString("rule")), &r)
		SecKillRules.R = r
	}
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
	Init()
	list := RetryTransList()
	if list != nil {
		logger.Info(list)
		list.Analyse()
	}
}
