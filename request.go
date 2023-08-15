package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// 推送消息结构
type resultContent struct {
	UserName     string
	Message      string
	UsedDays     string
	LeftDays     string
	Timestamp    string
	Usage        string
	scriptionURL string
}

// 响应消息结构
type statusResponse struct {
	Code    int                    `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

// 用量消息结构
type usageResponse struct {
	Code int       `json:"code"`
	Data [][]int64 `json:"data"`
}

const checkinURL = "https://glados.one/api/user/checkin"
const getStatusURL = "https://glados.one/api/user/status"
const usageURL = "https://glados.one/api/user/usage"
const consoleURL = "https://glados.one/console/clash"

// 解析json
func parseBody(body io.ReadCloser, respObj interface{}) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		// fmt.Printf("read body failed, err :%s", err.Error())
		panic(err)
	}
	err = json.Unmarshal(data, respObj)
	if err != nil {
		// fmt.Printf("parse body to json failed, err :%s", err.Error())
		panic(err)
	}
}

// 签到: 返回签到的信息
func checkin(cookie string) string {
	resObj := new(statusResponse)
	payload := strings.NewReader(`{"token": "glados.one"}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", checkinURL, payload)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Referer", "https://glados.one/console/checkin")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Origin", "https://glados.one")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// 解析
	parseBody(res.Body, resObj)
	if resObj.Code == 1 {
		logger.Info("Request succeed.")
	}
	return resObj.Message
}

// 获取订阅链接
func getSubscriptionURL(cookie string) (scriptionURL string) {
	var (
		node   []*cdp.Node
		header = map[string]interface{}{
			"Cookie": cookie,
		}
	)
	allocCtx, cancal := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancal()
	chromeCtx, cancal := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancal()
	logger.Debug("准备运行chrome浏览器")
	err := chromedp.Run(
		chromeCtx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers(header)),
		chromedp.Navigate(consoleURL),
		chromedp.Nodes(`//*[@type="button"]`, &node, chromedp.NodeVisible),
	)
	if len(node) == 0 {
		logger.Error("爬取订阅URL失败，请重试, 错误: %v", err.Error())
		return
	}
	if err != nil {
		logger.Error("浏览器执行任务失败，错误: %v", err.Error())
		return
	}
	scriptionURL = node[0].AttributeValue("data-clipboard-text")
	logger.Debug("退出chrome浏览器")
	return
}

// 获取时间和用量
func getTimeAndUsage(cookie string) (timestr, usage string) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", usageURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Referer", "https://glados.one/console/report")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Origin", "https://glados.one")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	respObj := new(usageResponse)
	parseBody(res.Body, respObj)
	if len(respObj.Data) == 0 { // 如果还没开始用
		timestr, usage = time.Now().Format("2006-01-02 15:04:05"), "0GB"
		return
	}
	if respObj.Code != 0 {
		fmt.Printf("get usage data failed, please check whether the cookie has expired...")
		return
	}

	// 输出
	data := respObj.Data
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

// 查看用户信息
func getStatus(cookie string) (string, string, string) {

	resObj := new(statusResponse)
	client := &http.Client{}
	req, err := http.NewRequest("GET", getStatusURL, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Referer", "https://glados.one/console/checkin")
	req.Header.Add("Cookie", cookie)
	req.Header.Add("Origin", "https://glados.one")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	parseBody(res.Body, resObj)

	// 解析body
	if len(resObj.Data) == 0 {
		// 失败
		return "", "", ""
	}

	email := resObj.Data["email"].(string)
	usedDays := fmt.Sprintf("%.0f", resObj.Data["days"].(float64))
	leftDays, err := strconv.ParseFloat(resObj.Data["leftDays"].(string), 64)
	if err != nil {
		panic("ParseFloat Failed.")
	}
	return email, usedDays, fmt.Sprintf("%.0f", leftDays)
}

func worker(cookie string) {
	mes := checkin(cookie)
	email, days, leftdays := getStatus(cookie)
	timestamp, usage := getTimeAndUsage(cookie)
	scriptionURL := getSubscriptionURL(cookie)
	userObj = &resultContent{
		Message:      mes,
		UserName:     email,
		UsedDays:     days,
		LeftDays:     leftdays,
		Timestamp:    timestamp,
		Usage:        usage,
		scriptionURL: scriptionURL,
	}
	// logger.Info("User: %v", userObj)
	messageChan <- userObj
}

func sendMessage() {
	for msg := range messageChan {
		noticeTofeishuMessage(msg)
		wg.Done()
	}
}

// 消费者
func consumer() {
	for user := range userChan {
		go worker(user.Cookie)
	}
}

// 生产者
func producer(conf *Config) {
	users := conf.Users
	for i := 0; i < len(users); i++ {
		wg.Add(1)
		userChan <- &users[i]
	}
	close(userChan)
}
