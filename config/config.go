package config

const ProNAME = "kd.pro"
const LogPATH = "/root/nginx/www/logs/" + ProNAME + "/"
const RunDURATION = 290
const AdminMailer = "fengelom@163.com"

var CurUser string

// sec kill
var SecKillFee float64
var SecKillRate float64
var SecKillRestDay int
var SecKillTime float64
var SleepT float64
var RuleKey string

var LocalIp string
var LockCode string

var JobType string
var JobList string
