package main

// 签到

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type checkinRequestBody struct {
	Token string `json:"token"`
}

type checkResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type statusResponse struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
}

// 返回用户名，使用天数，剩余天数
func checkin(cookie string) (string, string, string, string) {
	// 构造结构体
	reqBody := checkinRequestBody{
		Token: "glados_network",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	resp, err := newRequest("POST", "https://glados.one/api/user/checkin", "https://glados.one/console/checkin", cookie, reqBodyBytes)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respObj := &checkResponse{}
	parseBody(resp.Body, respObj)
	resp, err = newRequest("GET", "https://glados.one/api/user/status", "https://glados.one/console/checkin", cookie, nil)
	if err != nil {
		panic(err)
	}
	respObj2 := &statusResponse{}
	if respObj2.Code != 0 {
		panic("get status error")
	}
	parseBody(resp.Body, respObj2)
	// 解析内容
	if len(respObj2.Data) == 0 {
		panic("user info error, please check it!")
	}
	email := respObj2.Data["email"].(string)
	days := fmt.Sprintf("%.0f", respObj2.Data["days"].(float64))
	leftDays, _ := strconv.ParseFloat(respObj2.Data["leftDays"].(string), 64)

	return respObj.Message, email, days, fmt.Sprintf("%.0f", leftDays)
}

// func main() {
// 	Run()
// }
