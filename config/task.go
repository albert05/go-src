package config

//
var TaskList = map[string]map[string]string{
	"exchange": {
		"scriptName": "run",
		"params":     " -t %s -l %s",
	},
}
