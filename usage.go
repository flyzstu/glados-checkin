package main

// 获得当前用量
import (
	"fmt"
	"time"
)

type usageResponse struct {
	Code int       `json:"code"`
	Data [][]int64 `json:"data"`
}

// 获取时间戳和用量
func getTimeAndUsage(data [][]int64) (timestr, usage string) {
	var sum int64
	for index := range data {
		sum += data[index][1]
	}
	usage = fmt.Sprintf("%.2fGB", ((float32(sum)/1024)/1024)/1024)

	// 时间转换
	timestamp := data[0][0]
	timeobj := time.Unix(0, timestamp*1000000) // 毫秒转纳秒
	timestr = timeobj.Format("2006-01-02 15:04:05")

	return timestr, usage
}

func usage(cookie string) (timestr, usage string) {
	resp, err := newRequest("GET", "https://glados.one/api/user/usage", "https://glados.one/console/report", cookie, nil)
	if err != nil {
		fmt.Printf("request failed, err :%s", err.Error())
		return
	}
	defer resp.Body.Close()

	respObj := &usageResponse{}
	parseBody(resp.Body, respObj)
	if len(respObj.Data) == 0 { // 如果还没开始用
		timestr, usage = time.Now().Format("2006-01-02 15:04:05"), "0GB"
		return
	}
	if respObj.Code != 0 {
		fmt.Printf("get usage data failed, please check whether the cookie has expired...")
		return
	}

	// 输出
	timestr, usage = getTimeAndUsage(respObj.Data)
	return
}
