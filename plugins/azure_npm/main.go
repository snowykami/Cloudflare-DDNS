package main

import (
	"Cloudflare-DDNS/config"
	"Cloudflare-DDNS/model"
	"os"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

// Azure Nginx Proxy Manager Plugin
type AzureNpmPlugin struct{}

var PluginInstance AzureNpmPlugin

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func (p *AzureNpmPlugin) Entry() {
	logrus.Println("Load Azure Nginx Proxy Manager Plugin/蔚蓝Nginx代理管理器插件")
	logrus.Printf("当前ssh: %s@%s", config.Get("ssh_user"), config.Get("ssh_host"))
}

func (p *AzureNpmPlugin) OnIPChange(event model.IPChangeEvent) {
	// 在变动时执行sudo systemctl restart nginx
	// 等待dns记录更新后再重启nginx，120s
	logrus.Infof("IP变动，等待120s后重启代理nginx: %s -> %s", event.OldIP, event.NewIP)
	time.Sleep(120 * time.Second)
	RestarNginx()
}

func RestarNginx() {
	sshHost := config.Get("ssh_host")
	sshPort := config.GetWithDefault("ssh_port", "22")
	sshUser := config.Get("ssh_user")
	sshKeyPath := config.Get("ssh_key_path")

	cmd := exec.Command("ssh", "-4", "-vvv", "-i", sshKeyPath, "-p", sshPort, "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", sshUser+"@"+sshHost, "sudo systemctl restart nginx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Errorf("Failed to restart nginx/重启nginx失败: %v, output: %s", err, string(output))
	} else {
		logrus.Println("Nginx has been restarted successfully/成功重启nginx")
	}
}
