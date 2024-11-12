package service

import (
	"Cloudflare-DDNS/config"
	"Cloudflare-DDNS/model"
	"Cloudflare-DDNS/utils"
	log "github.com/sirupsen/logrus"
	"time"
)

// IPMonitor 监控本地IP变动
func IPMonitor() {
	log.Println("Start monitoring the local IP address. 开始监控本地IP地址变动")
	if config.EnableIPv4 {
		log.Println("IPv4 is enabled. IPv4已在DDNS启用")
	}
	if config.EnableIPv6 {
		log.Println("IPv6 is enabled. IPv6已在DDNS启用")
	}
	var currentIPv4, currentIPv6 string
	if config.EnableIPv4 {
		currentIPv4, _ = utils.GetCurrentIpv4()
	}
	if config.EnableIPv6 {
		currentIPv6, _ = utils.GetCurrentIpv6()
	}

	// Get the current IP address
	for {
		if config.EnableIPv4 {
			v4, err := utils.GetCurrentIpv4()
			if err != nil {
				log.Errorf("Failed to get the current IPv4 address: %v", err)
			} else {
				if v4 != currentIPv4 {
					eventChan <- &model.IPChangeEvent{
						OldIP: currentIPv4,
						NewIP: v4,
						Type:  "A",
					}
					log.Printf("IPv4 changed/IPv4变动从: %s -> %s\n", currentIPv4, v4)
					currentIPv4 = v4
				} else {
					log.Printf("IPv4 remains unchanged/IPv4还是这个: %s\n", currentIPv4)
				}
			}

		}

		if config.EnableIPv6 {
			v6, err := utils.GetCurrentIpv6()
			if err != nil {
				log.Errorf("Failed to get the current IPv6 address: %v, maybe your network does not support IPv6", err)
			} else {
				if v6 != currentIPv6 {
					eventChan <- &model.IPChangeEvent{
						OldIP: currentIPv6,
						NewIP: v6,
						Type:  "AAAA",
					}
					log.Printf("IPv6 changed/IPv6变动 %s -> %s\n", currentIPv6, v6)
					currentIPv6 = v6
				} else {
					log.Printf("IPv6 remains unchanged/IPv6还是这个: %s\n", currentIPv6)
				}
			}
		}
		time.Sleep(time.Duration(config.Config.Duration) * time.Second)
	}
}
