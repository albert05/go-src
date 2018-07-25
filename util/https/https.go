package https

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultContentTYPE = "application/x-www-form-urlencoded"
const JsonContentTYPE = "application/json"
const HttpSUCCESS = 0
const Timeout = 5

func Post(uri string, params map[string]string, cookie string) ([]byte, error) {
	v := url.Values{}
	for k, val := range params {
		v.Set(k, val)
	}

	//form数据编码
	body := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * Timeout,
		},
	}
	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	// set cookie
	if cookie != "" {
		request.Header.Add("cookie", "SESSIONID="+cookie)
	}

	//给一个key设定为响应的value
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;param=value")

	//发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(content))
	return content, nil
}

func PostWithoutCookie(url, params string) ([]byte, error) {
	resp, err := http.Post(url, DefaultContentTYPE, strings.NewReader(params))
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

func PostJson(url, params string) ([]byte, error) {
	resp, err := http.Post(url, JsonContentTYPE, strings.NewReader(params))
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

func PJson(url string, b io.Reader) ([]byte, error) {
	resp, err := http.Post(url, JsonContentTYPE, b)
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
