package service

import (
	"Cloudflare-DDNS/cloudflare"
	"Cloudflare-DDNS/config"
	"Cloudflare-DDNS/model"
	log "github.com/sirupsen/logrus"
)

// SyncedRecords 用于同步云端记录
var SyncedRecords *model.GetDnsRecordsResp // SyncRecords 用于同步云端记录

// SyncRoutine 同步云端记录
func SyncRecords() error {
	// Get the DNS records
	records, err := cloudflare.GetDnsRecords()
	if err != nil {
		return err
	}
	SyncedRecords = records
	return nil
}

func CloudMonitor() {
	// 预检查冲突内容并警告
Sync:
	err := SyncRecords()
	if err != nil {
		log.Errorf("Failed to sync the records from the cloud/从远端同步记录失败: %v", err)
		goto Sync
	} else {
		log.Println("Synced the records from the cloud/成功从远端同步记录")
	}

	for _, record := range SyncedRecords.Result {
		for _, task := range config.Config.DDNS {
			if record.Name == task.Name {
				if record.Type != "A" && record.Type != "AAAA" {
					log.Warnf("Domain has a conflict record type/域名有会和A/AAAA冲突的记录类型: %s\n", record.Type)
				}
			}
		}
	}

	// 收到本地变动通知后检查云端记录是否与本地新地址一致
	for event := range eventChan {
		// 云端更新
		err := SyncRecords()
		if err != nil {
			log.Errorf("Failed to sync the records from the cloud/从远端同步记录失败: %v", err)
			continue
		} else {
			log.Println("Synced the records from the cloud/成功从远端同步记录")
		}
		var Tasks []config.DDNSConfig
		if event.Type == "A" {
			Tasks = config.ADDNSTasks
		} else if event.Type == "AAAA" {
			Tasks = config.AAAADDNSTasks
		} else {
			log.Errorf("Unsupported DNS record type/不支持的DNS记录类型: %s", event.Type)
		}

		for _, task := range Tasks {
			needCreate := true
			for _, record := range SyncedRecords.Result {
				// 有则更新
				if record.Name == task.Name {
					if record.Type == event.Type {
						// 类型相同
						if record.Content != event.NewIP {
							needCreate = false
							newReqData := model.PostDNSRecord{
								Content: event.NewIP,
								Name:    task.Name,
								Proxied: task.Proxied,
								Type:    event.Type,
								TTL:     task.TTL,
								Comment: task.Comment,
							}
							resp, err := cloudflare.UpdateDNSRecord(newReqData, record.Id)
							if err != nil {
								log.Errorf("Failed to update the DNS record/更新DNS记录失败: %v", err)
							} else if !resp.Success {
								log.Errorf("Failed to update the DNS record/更新DNS记录失败: %v", resp.Errors)
							} else {
								log.Printf("Updated the DNS record/更新DNS记录成功: %s -> %s", resp.Result.Name, resp.Result.Content)
							}
						} else {
							needCreate = false
						}
					}
				}
			}
			if needCreate {
				// 无则添加
				newReqData := model.PostDNSRecord{
					Content: event.NewIP,
					Name:    task.Name,
					Proxied: task.Proxied,
					Type:    event.Type,
					TTL:     task.TTL,
					Comment: task.Comment,
				}
				resp, err := cloudflare.CreateDNSRecord(newReqData)
				if err != nil {
					log.Errorf("Failed to create the DNS record/创建DNS记录失败: %v", err)
				} else if !resp.Success {
					log.Errorf("Failed to create the DNS record/创建DNS记录失败: %v", resp.Errors)
				} else {
					log.Printf("Created the DNS record/创建DNS记录成功: %s -> %s", resp.Result.Name, resp.Result.Content)
				}
			}
		}
		// 每次事件处理后同步一次云端记录

	}
}
