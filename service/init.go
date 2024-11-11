package service

import (
	"Cloudflare-DDNS/model"
)

var eventChan = make(chan *model.IPChangeEvent, 100)

func Start() {
	// 从云端同步记录
	go CloudMonitor()
	// 启动IP监控
	go IPMonitor()

	select {}
}
