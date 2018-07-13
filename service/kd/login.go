package kd

import (
	"fmt"
	"encoding/json"
	"errors"
	"kd.explorer/tools/https"
	"kd.explorer/model"
)

type loginResponse struct {
	Code int ``
	Sessionid string ``
}

const LoginURL = "https://deposit.koudailc.com/user/login"

func Login(username, password string) (string, error) {
	params := fmt.Sprintf("username=%s&password=%s", username, password)
	body, err := https.PostWithoutCookie(LoginURL, params)
	if err != nil {
		return "", err
	}

	//fmt.Println(string(body))

	var result loginResponse
	json.Unmarshal(body, &result)

	if https.HttpSUCCESS == result.Code {
		return result.Sessionid, nil
	}

	return "", errors.New("login request result failed")
}

func LoginK(user string) (string, error) {
	userInfo := model.FindUser(user)

	username := userInfo.GetAttrString("user_name")
	password := userInfo.GetAttrString("password")
	cookie, err := Login(username, password)
	if err != nil {
		return "", err
	}

	return cookie, nil
}