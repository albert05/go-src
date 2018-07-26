package service

import (
	"encoding/json"
	"fmt"
	"kd.explorer/config"
	"kd.explorer/model"
	"kd.explorer/util/https"
	"kd.explorer/util/mail"
)

const TransOrderURL = "https://deposit.koudailc.com/credit/apply-assignment"

type OrderResp struct {
	Code int `json:"code"`
	Uid  int `json:"uid"`
}

// 单账号秒杀
func (item *TransferItem) RunKill(cookie string) {
	params := item.MakeOrderParams(config.CurUser)
	body, err := https.Post(TransOrderURL, params, cookie)
	if err != nil {
		mail.SendSingle(config.AdminMailer, "高息转让项目提醒", err.Error())
		fmt.Println(err)
		return
	}

	var result OrderResp
	json.Unmarshal(body, &result)

	msg := fmt.Sprintf("user:%s 购买转让项目invest_id：%s 结果：%s", config.CurUser, item.InvestId, string(body))
	fmt.Println(msg)
	email := model.FindUser(config.CurUser).GetAttrString("email")
	if email != config.AdminMailer {
		mail.SendSingle(config.AdminMailer, "高息转让项目提醒", msg)
	}
	mail.SendSingle(email, "高息转让项目抢购提醒", msg)
}

// 多线程多账号异步秒杀
func (item *TransferItem) SyncRunKill() {
	ch := make(chan bool)
	for _, user := range config.SecKillList {
		go item.runT(user, ch)
	}

	for range config.SecKillList {
		<-ch
	}

	close(ch)
}

func (item *TransferItem) runT(user string, ch chan<- bool) {
	cookie, err := LoginK(user)
	if err != nil {
		fmt.Println(err)
		ch <- false
		return
	}

	params := item.MakeOrderParams(user)
	body, err := https.Post(TransOrderURL, params, cookie)
	if err != nil {
		fmt.Println(err)
		ch <- false
		return
	}

	var result OrderResp
	json.Unmarshal(body, &result)

	if result.Code == 0 && result.Uid != 0 {
		msg := fmt.Sprintf("user:%s 购买转让项目invest_id：%s 成功", user, item.InvestId)
		fmt.Println(msg)
		mail.SendSingle(config.AdminMailer, "高息转让项目抢购成功提醒", msg)
		ch <- true
		return
	}

	fmt.Println(fmt.Sprintf("user:%s 购买转让项目invest_id：%s 失败", user, item.InvestId))
	ch <- false
	return
}

func (item *TransferItem) MakeOrderParams(user string) map[string]string {
	userInfo := model.FindUser(user)
	paypasswd := userInfo.GetAttrString("pay_passwd")

	return map[string]string{
		"invest_id":    item.InvestId,
		"pay_password": paypasswd,
		"use_remain":   "1",
		"is_kdb_pay":   "0",
	}
}
