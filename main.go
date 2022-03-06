package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

type userObj struct {
	Email  string `json:"email"`
	Cookie string `json:"cookie"`
}

// 请求
type RequestContent struct {
	Type    string `json:"msg_type"`
	Content string `json:"content"`
}

func main() {
	var (
		users        = []userObj{}
		reqBodyBytes = []byte{}
	)
	// 读取users.json
	data, err := ioutil.ReadFile("user.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < len(users); i++ {
		wg.Add(1)
		go func(u userObj) {
			w := RequestContent{}
			w.Type = "text"
			checkMsg, user, useDays, leftDays := checkin(u.Cookie)
			timeStr, usage := usage(u.Cookie)
			w.Content = fmt.Sprintf("GLaDOS任务提醒\n用户%s签到信息:%s，GlaDOS服务已经使用了%s天，剩余%s天，截至到%s已经使用了%s", user, checkMsg, useDays, leftDays, timeStr, usage)
			reqBodyBytes, err = json.Marshal(&w)
			if err != nil {
				panic(err)
			}
			newRequest("POST", "http://127.0.0.1:7890/sendMessage", "", "wf09", reqBodyBytes)
			wg.Done()
		}((users)[i])
	}
	wg.Wait()

}
