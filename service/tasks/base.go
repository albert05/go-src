package tasks

import (
	"kd.explorer/util/https"
	"encoding/json"
	"kd.explorer/util/logger"
)

func Go(url string, cookie string, params map[string]string) (bool, string) {
	body, err := https.Post(url, params, cookie)
	if err != nil {
		logger.Info(err)
		return false, "http failed"
	}

	var result TaskResponse
	json.Unmarshal(body, &result)
	logger.Info(string(body))

	return https.HttpSUCCESS == result.Code, result.Message
}
