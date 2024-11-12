package main

import (
	"Cloudflare-DDNS/model"

	"github.com/sirupsen/logrus"
)

type ExamplePlugin struct{}

func (p *ExamplePlugin) Entry() {
	logrus.Println("Load Example Plugin")
}

func (p *ExamplePlugin) OnIPChange(event model.IPChangeEvent) {
	logrus.Printf("Example Plugin: IP changed from %s to %s\n", event.OldIP, event.NewIP)
}

var PluginInstance ExamplePlugin // 导出插件实例
