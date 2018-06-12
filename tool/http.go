package tool

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
)

const DEFAULT_CONTENT_TYPE = "application/x-www-form-urlencoded"
const HTTP_SUCCESS = 0

func Post(urlx string, params map[string]string, cookie string) ([]byte, error) {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}

	//form数据编码
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	request, err := http.NewRequest("POST", urlx, body)
	if err != nil {
		return nil, err
	}

	// set cookie
	if cookie != "" {
		request.Header.Add("cookie", "SESSIONID=" + cookie)
	}

	//给一个key设定为响应的value
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")

	//发送请求
	resp, err := client.Do(request)
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(content))
	return content, nil
}


func PostWithoutCookie(url, params string) ([]byte, error) {
	resp, err := http.Post(url, DEFAULT_CONTENT_TYPE, strings.NewReader(params))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}