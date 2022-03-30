package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type feishuContent struct {
	MsgType   string                 `json:"msg_type"`
	Content   map[string]interface{} `json:"content"`
	Timestamp int64                  `json:"timestamp"`
	Sign      string                 `json:"sign"`
}

func sendMessageToFeishu(msg string) (err error) {
	var fc feishuContent
	var fcBytes []byte
	// data, _ := ioutil.ReadFile("feishu_msg.json")
	// 生成签名
	timestamp := time.Now().Unix()
	sign, err := genFeishuMessageSign(conf.HashKey, timestamp)
	if err != nil {
		err = fmt.Errorf("生成签名失败, 错误:%s", err.Error())
		logger.Error(err.Error())
		return
	}
	// 生成data
	fc.Timestamp = timestamp
	fc.MsgType = "text"
	fc.Sign = sign
	fc.Content = make(map[string]interface{})
	fc.Content["text"] = msg

	fcBytes, err = json.Marshal(fc)
	if err != nil {
		err = fmt.Errorf("json序列化失败, 错误:%s", err.Error())
		logger.Error(err.Error())
		return
	}
	// 发送请求
	resp, err := http.Post(conf.BOT_API, "application/json;charset=utf-8", bytes.NewBuffer(fcBytes))
	if err != nil {
		err = fmt.Errorf("POST failed, err:%s", err.Error())
		logger.Error(err.Error())
		return
	}
	// 检查返回码
	if resp.StatusCode != 200 {
		logger.Debug("无法推送到飞书, 状态码:", resp.StatusCode)
	} else {
		logger.Debug("飞书机器人推送成功")
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("json读取失败, 错误:%s", err.Error())
		logger.Error(err.Error())
		return
	}
	logger.Debug(string(respData))
	return nil
}

// 生成飞书的消息签名
func genFeishuMessageSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
