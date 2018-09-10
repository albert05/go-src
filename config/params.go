package config

var TaskList = map[string]map[string]string{
	"exchange": {
		"scriptName": "task",
		"params":     " -tp %s -l %s",
	},
	"daily": {
		"scriptName": "task",
		"params":     " -tp %s -l %s",
	},
}

var SecKillList = []string{
	"cwf",
}
