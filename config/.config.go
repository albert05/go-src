package config

const DSN  = "user:pwd@tcp(127.0.0.1:3306)/db"
const RunDURATION = 290

var SmsConfig = map[string]string {
	"userid": "***",
	"account": "***",
	"password": "***",
}

var SmsReceiverList = []int {
	18721809992,
	13042160232,
}
