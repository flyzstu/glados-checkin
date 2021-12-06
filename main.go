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
type wxMsg struct {
	// Reserved field to add some meta information to the API response
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
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
			w := wxMsg{}
			w.Title = "GLaDOS任务提醒"
			w.URL = "https://fly97.cn"
			checkFlag, user, useDays, leftDays := checkin(u.Cookie)
			timeStr, usage := usage(u.Cookie)
			if checkFlag {
				w.Description = fmt.Sprintf("用户%s签到成功，GlaDOS服务已经使用了%s天，剩余%s天，截至到%s已经使用了%s", user, useDays, leftDays, timeStr, usage)
			} else {
				w.Description = fmt.Sprintf("用户%s签到失败，GlaDOS服务已经使用了%s天，剩余%s天，截至到%s已经使用了%s", user, useDays, leftDays, timeStr, usage)
			}

			reqBodyBytes, err = json.Marshal(&w)
			if err != nil {
				panic(err)
			}
			newRequest("POST", "http://127.0.0.1:8080/wx-api/sendMessage", "", "wf09", reqBodyBytes)
			wg.Done()
		}((users)[i])
	}
	wg.Wait()

}
