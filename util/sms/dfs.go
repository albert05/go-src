// 鼎丰所
package sms

import (
	"fmt"
	"kd.explorer/config"
	"strconv"
	"encoding/xml"
	"errors"
	"kd.explorer/util/https"
)

const SendURL = "https://115.29.242.32:8888/sms.aspx?action=send"
const SmsSuccessSTATUS  = "Success"
const DefaultSIGN = "【鼎丰所】"

type Sms struct {
	Status string `returnstatus`
	Message string `message`
}

type SmsResponse struct {
	XMLName     xml.Name `xml:"returnsms"`
	Version     string   `xml:"version"`
	Status string `xml:"returnstatus"`
	Message string `xml:"message"`
}

func Send(phone int, content string) error {
	params := config.SmsConfig
	params["mobile"] = strconv.Itoa(phone)
	params["content"] = DefaultSIGN + content
	body, err := https.Post(SendURL, params, "")
	if err != nil {
		return err
	}

	var result SmsResponse
	err = xml.Unmarshal(body, &result)

	if SmsSuccessSTATUS != result.Status {
		return errors.New(fmt.Sprintf("returnstatus:%s,message:%s", result.Status, result.Message))
	}

	return nil
}
