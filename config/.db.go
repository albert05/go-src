package config

const DSN = "user:pwd@tcp(127.0.0.1:3306)/db"

var SmsConfig = map[string]string{
	"userid":   "...",
	"account":  "...",
	"password": "...",
}

var MailConfig = map[string]string{
	"host":     "smtp.qq.com",
	"port":     "465",
	"username": "...",
	"password": "...",
}
