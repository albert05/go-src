package config

const DSN  = "user:pwd@tcp(127.0.0.1:3306)/db"
const RunDURATION = 290

var SmsConfig = map[string]string {
	"userid": "***",
	"account": "***",
	"password": "***",
}

var SmsReceiverList = []int {
	1000,
	10001,
}

var MailConfig = map[string]string {
	"host": "smtp.qq.com",
	"port": "465",
	"username": "***",
	"password": "***",
}

var MailReceiverList = []string {
	"***",
}

var SecKillList = []string {
	"***",
}

var CurUser string

var SecKillFee float64
var SecKillRate float64
var SecKillRestDay int
