package utils

import (
	"Cloudflare-DDNS/model"
	"os"
	"plugin"

	"github.com/sirupsen/logrus"
)

var Plugins = make(map[string]Plugin)

func InitPlugins() {
	// 检测插件文件夹是否存在不存在则创建
	if !PathExists("plugins") {
		err := os.Mkdir("plugins", os.ModePerm)
		if err != nil {
			logrus.Errorf("Failed to create the plugin folder/创建插件文件夹失败: %v", err)
		} else {
			logrus.Println("The plugin folder has been created successfully/插件文件夹创建成功")
		}
	}
	// 读取插件文件夹下的所有插件so
	files, err := os.ReadDir("plugins")
	if err != nil {
		logrus.Errorf("Failed to read the plugin folder/读取插件文件夹失败: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if file.Name()[len(file.Name())-3:] == ".so" {
			// 加载插件
			p, err := plugin.Open("plugins/" + file.Name())
			if err != nil {
				logrus.Errorf("Failed to load the plugin %s/加载插件失败: %v", file.Name(), err)
			} else {
				// 获取插件入口

				// 断言插件实例
				var pluginInstance Plugin
				sym, err := p.Lookup("PluginInstance")
				if err != nil {
					logrus.Errorf("Failed to lookup the plugin instance/查找插件实例失败: %v", err)
				}
				pluginInstance, ok := sym.(Plugin)
				if !ok {
					logrus.Errorf("Failed to assert the plugin instance/插件实例断言失败")
				} else {
					Plugins[file.Name()] = pluginInstance
					pluginInstance.Entry()
				}

			}
		}
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type Plugin interface {
	Entry()
	OnIPChange(model.IPChangeEvent)
}
