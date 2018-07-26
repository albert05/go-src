package config

var TaskList = map[string]map[string]string{
	"exchange": {
		"scriptName": "run",
		"params":     " -tp %s -l %s",
	},
}

var SecKillList = []string{
	"cwf",
}
