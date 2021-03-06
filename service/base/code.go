package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kd.explorer/common"
	"kd.explorer/util/dates"
	"kd.explorer/util/https"
	"os"
	"strconv"
	"strings"
)

const CodeURL = "https://deposit.koudailc.com%s"
const RefreshURL = "https://deposit.koudailc.com/user/captcha?refresh"

var ImagePATH string
var ImagePrefixPATH string

type CodeResponse struct {
	Hash1 int    ``
	Url   string ``
}

type Code struct {
	Cookie   string ``
	Url      string ``
	FileName string ``
}

func init() {
	// 目录暂设置在laravel
	ImagePrefixPATH = "/www/laravel/public/"
	ImagePATH = ImagePrefixPATH + "goimg/"
}

func (code *Code) SetCookie(cookie string) {
	code.Cookie = cookie
}

func (code *Code) GetFileName() string {
	return strings.TrimPrefix(code.FileName, ImagePrefixPATH)
}

func (code *Code) Refresh() {
	params := map[string]string{}

	body, err := https.Post(RefreshURL, params, code.Cookie)
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
	body, err := https.Post(code.Url, params, code.Cookie)
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
