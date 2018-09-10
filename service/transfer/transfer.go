package transfer

import (
	"encoding/json"
	"fmt"
	"kd.explorer/config"
	"kd.explorer/util/https"
	"strconv"
	"kd.explorer/service/base"
	"strings"
	"kd.explorer/util/mail"
	"kd.explorer/model"
	"kd.explorer/util/dates"
	"kd.explorer/util/logger"
)

const TransferListURL = "https://deposit.koudailc.com/credit/market-for-app-v2?appVersion=6.7.5&osVersion=11.300000&clientType=ios&deviceName=iPhone%20X&page=1&pageSize=2&sortRuleType=2"
const DefaultUserKEY = "cwf_2551"
const TransLoginSuccessSTATUS = 1
const RetryCNT = 3

var Cookie string

type TransferItem struct {
	Id         string `json:"id"`
	InvestId   string `json:"invest_id"`
	AssignFee  string `json:"assign_fee"`
	AssignRate string `json:"assign_rate"`
	RestDays   int    `json:"rest_days"`
	UpdatedAt  int64  `json:"updated_at"`
}

func (item *TransferItem) GetFee() float64 {
	fee, err := strconv.ParseFloat(item.AssignFee, 64)
	if err != nil {
		return 0
	}

	return fee / 100
}

func (item *TransferItem) GetRate() float64 {
	rate, err := strconv.ParseFloat(item.AssignRate, 64)
	if err != nil {
		return 0
	}

	return rate
}

func (item *TransferItem) GetKey() string {
	return item.Id + "-" + item.InvestId + "-" + config.CurUser
}

func (item *TransferItem) String() string {
	b, _ := json.Marshal(item)

	return string(b)
}

func GetTransferMsg(item TransferItem) string {
	return fmt.Sprintf("转让年化：%.2f%s, 金额：%.2f, 剩余天数：%d", item.GetRate(), "%", item.GetFee(), item.RestDays)
}

type TransTmp struct {
	Code  int            ``
	Items []TransferItem `json:"creditItems"`
}

type TransList struct {
	Code    int
	List    TransTmp `json:"recentlyPublishedItems"`
	IsLogin int      `json:"is_login"`
	Cookie  string
}

func InitCookie(isFlush bool) {
	if !isFlush && Cookie != "" {
		return
	}

	user := DefaultUserKEY
	if config.CurUser != "" {
		user = config.CurUser
	}
	cookie, err := base.LoginK(user)
	if err == nil {
		Cookie = cookie
	}
}

func GetTransferList() (*TransList, error) {
	body, err := https.Post(TransferListURL, nil, Cookie)
	if err != nil {
		return nil, err
	}

	var result TransList
	json.Unmarshal(body, &result)

	if TransLoginSuccessSTATUS != result.IsLogin {
		logger.Info(string(body))
	}

	result.Cookie = Cookie
	return &result, nil
}

func RetryTransList() *TransList {
	var i = 0
	isFlush := false
	for i < RetryCNT {
		InitCookie(isFlush)
		list, err := GetTransferList()
		if err != nil {
			logger.Error(err)
			return nil
		}

		if TransLoginSuccessSTATUS == list.IsLogin {
			return list
		}
		i++
		isFlush = true
	}

	return nil
}

func (list *TransList) Analyse() {
	monitorMsg := make([]string, 0)
	for _, item := range list.List.Items {
		if !CheckIsSended(item.GetKey(), item.String()) {
			if true == SecKillRules.Check(item) {
				dates.SleepSecond(config.SecKillTime)
				if config.SecUser != "" {
					if cookie, err := base.LoginK(config.SecUser); err == nil {
						item.RunKill(cookie)
					}
				}
			}
			if true == MonitorRule.Check(item) {
				monitorMsg = append(monitorMsg, GetTransferMsg(item))
			}
		}
	}

	logger.Info(monitorMsg)

	// is send monitor msg
	if len(monitorMsg) > 0 {
		msg := "高息转让项目提醒 >> " + strings.Join(monitorMsg, "@@")
		logger.Info(msg)
		// send mail
		email := model.FindUser(config.CurUser).GetAttrString("email")
		mail.SendSingle(email, "高息转让项目提醒", msg)
	}
}

