package config

var TaskList = map[string]map[string]string{
	"exchange": {
		"scriptName": "exchange",
		"params":     " -tp %s -l %s",
	},
}

var SecKillList = []string{
	"cwf",
}
