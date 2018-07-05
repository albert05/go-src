package mail

import (
	"gopkg.in/gomail.v2"
	"kd.explorer/config"
	"strconv"
)

const defaultPORT = 465 // qq SMTP

func Send(receivers []string, subject, content string) bool {
	mail := gomail.NewMessage()
	mail.SetAddressHeader("From", config.MailConfig["username"], "")
	for _, receiver := range receivers {
		mail.SetHeader("To",
			mail.FormatAddress(receiver, ""),
		)
	}

	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", content)

	port, _ := strconv.Atoi(config.MailConfig["port"])
	d := gomail.NewPlainDialer(config.MailConfig["host"], port, config.MailConfig["username"], config.MailConfig["password"])

	if err := d.DialAndSend(mail); err != nil {
		return false
	}

	return true
}
