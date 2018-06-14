package service

import (
	"kd.explorer/tool"
	"encoding/json"
	"fmt"
	"kd.explorer/common"
	"strconv"
	"io"
	"bytes"
	"os"
)

const CODE_URL = "http://deposit.koudailc.com%s";
const REFRESH_URL = "http://deposit.koudailc.com/user/captcha?refresh";

var IMAGE_PATH string

type CodeResponse struct {
	Hash1 int ``
	Url string ``
}

type Code struct {
	Cookie string ``
	Url string ``
	FileName string ``
}

func init() {
	// 目录暂设置在laravel
	IMAGE_PATH = "/root/nginx/www/laravel/public/goimg/"
}

func (code *Code) setCookie(cookie string) {
	code.Cookie = cookie
}

func (code *Code) Refresh() {
	params := map[string]string{}

	body, err := tool.Post(REFRESH_URL, params, code.Cookie)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result CodeResponse
	json.Unmarshal(body, &result)

	code.Url = fmt.Sprintf(CODE_URL, result.Url)
}

func (code *Code) CreateImage() {
	params := map[string]string{}
	body, err := tool.Post(code.Url, params, code.Cookie)
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := os.Create(code.FileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (code *Code) RandFileName() {
	randNum := common.GenerateRangeNum(10000, 99999)
	code.FileName = IMAGE_PATH + "captcha_" + strconv.Itoa(randNum) + ".png"
}