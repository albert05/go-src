package kd

import (
	"fmt"
	"encoding/json"
	"errors"
	"kd.explorer/tools/http"
)

type loginResponse struct {
	Code int ``
	Sessionid string ``
}

const LoginURL = "http://deposit.koudailc.com/user/login"

func Login(username, password string) (string, error) {
	params := fmt.Sprintf("username=%s&password=%s", username, password)
	body, err := http.PostWithoutCookie(LoginURL, params)
	if err != nil {
		return "", err
	}

	//fmt.Println(string(body))

	var result loginResponse
	json.Unmarshal(body, &result)

	if http.HttpSUCCESS == result.Code {
		return result.Sessionid, nil
	}

	return "", errors.New("login request result failed")
}