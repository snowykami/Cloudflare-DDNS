package main

import (
	"Cloudflare-DDNS/service"
	"Cloudflare-DDNS/utils"
)

func main() {
	utils.InitPlugins()
	service.Start()
}
