package kd

import (
	"encoding/json"
	"strconv"
	"kd.explorer/tools/https"
	"kd.explorer/config"
	"fmt"
)

const TransferListURL = "https://deposit.koudailc.com/credit/market-for-app-v2?appVersion=6.7.5&osVersion=11.300000&clientType=ios&deviceName=iPhone%20X&page=1&pageSize=2&sortRuleType=2"
const DefaultUserKEY = "cwf_2551"
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
	return item.Id + "-" + item.InvestId + "-" + config.CurUser
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
	Cookie string
}

func InitCookie(isFlush bool) {
	if !isFlush && Cookie != "" {
		return
	}

	user := DefaultUserKEY
	if config.CurUser != "" {
		user = config.CurUser
	}
	cookie, err := LoginK(user)
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
		fmt.Println(string(body))
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
			fmt.Println(err)
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
