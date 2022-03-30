package main

import (
	"flag"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/flyzstu/mylog"
)

var opts []func(*chromedp.ExecAllocator)
var wg sync.WaitGroup
var logger = mylog.New()
var userChan = make(chan *User, 100)
var messageChan = make(chan *resultContent, 100)
var userObj *resultContent
var conf *Config

func main() {
	opts = append(chromedp.DefaultExecAllocatorOptions[:],
		// chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("headless", true), // 无头模式
	)
	var (
		confName = flag.String("f", "user.yaml", "Your configuration.")
	)
	flag.Parse()
	logger.Info("解析配置文件成功")
	conf = loadConf(*confName)
	producer(conf)
	go consumer()
	go sendMessage()
	wg.Wait()
}
