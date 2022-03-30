package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func noticeTofeishuMessage(result *resultContent) {
	var msg string
	msg = "GLaDOS Checker 签到提醒\n"
	msg += fmt.Sprintf("触发时间：%s\n", time.Now().Format("2006-01-02 15:04:05"))
	msg += fmt.Sprintf("用户名：%s\n", result.UserName)
	msg += fmt.Sprintf("订阅URL：%s\n", result.scriptionURL)
	msg += fmt.Sprintf("标题：%s\n", result.Message)
	msg += fmt.Sprintf("服务使用了：%s天\n", result.UsedDays)
	msg += fmt.Sprintf("服务还剩：%s天\n", result.LeftDays)
	msg += fmt.Sprintf("已经使用了：%s\n", result.Usage)
	msg += fmt.Sprintf("上次使用时间：%s\n", result.Timestamp)
	msg += fmt.Sprintf("事件ID：%s\n", uuid.New().String())
	err := sendMessageToFeishu(msg)
	if err != nil {
		panic(err)
	}
}
