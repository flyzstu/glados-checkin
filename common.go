// 新建一个网络请求
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// 新建一个网络请求
func newRequest(method, url, referer string, cookie string, reqBodyBytes []byte) (resp *http.Response, err error) {
	var req *http.Request
	client := &http.Client{}
	if reqBodyBytes != nil {
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(reqBodyBytes))
		req.Header.Add("Content-Type", "application/json;charset=utf-8") // post json
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Referer", referer)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:94.0) Gecko/20100101 Firefox/94.0")

	resp, err = client.Do(req)

	if err != nil {
		fmt.Printf("request failed, err :%s", err.Error())
		return nil, err
	}
	return resp, err
}

// 解析body
func parseBody(body io.ReadCloser, respObj interface{}) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Printf("read body failed, err :%s", err.Error())
		return
	}
	err = json.Unmarshal(data, respObj)
	if err != nil {
		fmt.Printf("parse body to json failed, err :%s", err.Error())
		return
	}
}
