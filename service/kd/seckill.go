package kd

import (
	"kd.explorer/config"
	"fmt"
	"kd.explorer/tools/https"
	"encoding/json"
)

const TransOrderURL = "https://deposit.koudailc.com/credit/apply-assignment"

type OrderResp struct {
	Code int `json:"code"`
	Uid int `json:"uid"`
}

func (item *TransferItem) RunKILL() {
	for _, user := range config.SecKillList {
		cookie, err := LoginK(user)
		if err != nil {
			fmt.Println(err)
			continue
		}

		params := item.MakeOrderParams(user)
		body, err := https.Post(TransOrderURL, params, cookie)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result OrderResp
		json.Unmarshal(body, &result)

		if result.Code == 0 && result.Uid != 0 {
			fmt.Println(fmt.Sprintf("转让项目invest_id：%s 购买成功", item.InvestId))
		}
	}
}

func (item *TransferItem) MakeOrderParams(user string) map[string]string {
	userInfo := FindUser(user)
	paypasswd := userInfo.GetAttrString("pay_passwd")

	return map[string]string{
		"invest_id": item.InvestId,
		"pay_password": paypasswd,
		"use_remain": "1",
		"is_kdb_pay": "0",
	}
}
