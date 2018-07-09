package kd

import (
	"encoding/json"
	"fmt"
	"kd.explorer/common"
	"strconv"
	"io"
	"bytes"
	"os"
	"strings"
	"kd.explorer/tools/http"
	"kd.explorer/tools/dates"
)

const CodeURL = "http://deposit.koudailc.com%s"
const RefreshURL = "http://deposit.koudailc.com/user/captcha?refresh"

var ImagePATH string
var ImagePrefixPATH string

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
	ImagePrefixPATH = "/root/nginx/www/laravel/public/"
	ImagePATH = ImagePrefixPATH + "goimg/"
}

func (code *Code) setCookie(cookie string) {
	code.Cookie = cookie
}

func (code *Code) getFileName() string {
	return strings.TrimPrefix(code.FileName, ImagePrefixPATH)
}

func (code *Code) Refresh() {
	params := map[string]string{}

	body, err := http.Post(RefreshURL, params, code.Cookie)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result CodeResponse
	json.Unmarshal(body, &result)

	code.Url = fmt.Sprintf(CodeURL, result.Url)
}

func (code *Code) CreateImage() {
	params := map[string]string{}
	body, err := http.Post(code.Url, params, code.Cookie)
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
	dateStr := strconv.FormatInt(dates.NowTime(), 32)
	code.FileName = ImagePATH + "captcha_" + dateStr + "_" + strconv.Itoa(randNum) + ".png"
}
