package kd

import (
	"fmt"
	"kd.explorer/tools/mysql"
	"encoding/json"
	"strconv"
	"kd.explorer/tools/https"
)

const TransferListURL = "https://deposit.koudailc.com/credit/market-for-app-v2?appVersion=6.7.5&osVersion=11.300000&clientType=ios&deviceName=iPhone%20X&page=1&pageSize=5&sortRuleType=2"
const DefaultUserKEY = "cwf"
const TransLoginSuccessSTATUS = 1
const RetryCNT = 3

var Cookie string

type TransferItem struct {
	Id string `json:"id"`
	InvestId string `json:"invest_id"`
	AssignFee string `json:"assign_fee"`
	AssignRate string `json:"assign_rate"`
	RestDays int `json:"rest_days"`
}

func (item *TransferItem) GetFee() float64 {
	fee, err := strconv.ParseFloat(item.AssignFee, 64)
	if err != nil {
		return 0
	}

	return fee /100
}

func (item *TransferItem) GetRate() float64 {
	rate, err := strconv.ParseFloat(item.AssignRate, 64)
	if err != nil {
		return 0
	}

	return rate
}

func (item *TransferItem) GetKey() string {
	return item.Id + "-" + item.InvestId
}

func (item *TransferItem) String() string {
	b, _ := json.Marshal(item)

	return string(b)
}

type TransTmp struct {
	Code int ``
	Items []TransferItem `json:"creditItems"`
}

type TransList struct {
	Code int
	List TransTmp `json:"recentlyPublishedItems"`
	IsLogin int `json:"is_login"`
}

func InitCookie() {
	userInfo, err := mysql.Conn.FindOne(fmt.Sprintf("SELECT * FROM userinfos WHERE user_key = '%s'", DefaultUserKEY))
	if err != nil || userInfo == nil {
		// log
	}

	username := userInfo.GetAttrString("user_name")
	password := userInfo.GetAttrString("password")
	Cookie, err = Login(username, password)
	if err != nil {
		//
	}
}

func GetTransferList() (*TransList, error) {
	body, err := https.Post(TransferListURL, nil, Cookie)
	if err != nil {
		return nil, err
	}

	var result TransList
	json.Unmarshal(body, &result)

	return &result, nil
}

func RetryTransList() *TransList {
	var i = 0
	for i < RetryCNT {
		InitCookie()
		list, err := GetTransferList()
		if err != nil {
			return nil
		}

		if TransLoginSuccessSTATUS == list.IsLogin {
			return list
		}
		i++
	}

	return nil
}
