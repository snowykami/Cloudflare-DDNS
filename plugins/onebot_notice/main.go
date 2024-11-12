package main

import (
	"Cloudflare-DDNS/config"
	"Cloudflare-DDNS/model"

	"github.com/sirupsen/logrus"
)

// 在收到变化通知后通过支持OneBot标准的机器人发送到超级用户

type OneBotNoticePlugin struct{}

var PluginInstance OneBotNoticePlugin
var (
	OneBotInstance string
	OneBotUser     string
	OneBotToken    string
)

func (p *OneBotNoticePlugin) Entry() {
	OneBotInstance = config.Get("onebot_instance")
	OneBotUser = config.Get("onebot_user")
	OneBotToken = config.Get("onebot_token")

	logrus.Println("Load OneBot Notice Plugin: ", OneBotInstance)
}

func (p *OneBotNoticePlugin) OnIPChange(event model.IPChangeEvent) {
	// 在变动时执行sudo systemctl restart nginx
	// 等待dns记录更新后再重启nginx，120s
	logrus.Infof("IP变动，发送到OneBot: %s -> %s", event.OldIP, event.NewIP)
}
