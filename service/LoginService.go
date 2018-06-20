package service

import (
	"fmt"
	"encoding/json"
	"errors"
	"kd.explorer/tool"
)

type loginResponse struct {
	Code int ``
	Sessionid string ``
}

const LoginURL = "http://deposit.koudailc.com/user/login"

func Login(username, password string) (string, error) {
	params := fmt.Sprintf("username=%s&password=%s", username, password)
	body, err := tool.PostWithoutCookie(LoginURL, params)
	if err != nil {
		return "", err
	}

	//fmt.Println(string(body))

	var result loginResponse
	json.Unmarshal(body, &result)

	if tool.HttpSUCCESS == result.Code {
		return result.Sessionid, nil
	}

	return "", errors.New("login request result failed")
}