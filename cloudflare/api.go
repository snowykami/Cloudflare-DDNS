package cloudflare

import (
	"Cloudflare-DDNS/config"
	"Cloudflare-DDNS/model"
	"Cloudflare-DDNS/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetDnsRecords() (*model.GetDnsRecordsResp, error) {
	// Get the DNS records
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", config.Config.ZoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Email", config.Config.ApiEmail)
	req.Header.Add("X-Auth-Key", config.Config.ApiKey)
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// Unmarshal
	var records model.GetDnsRecordsResp

	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return nil, err
	}
	return &records, nil
}

func CreateDNSRecord(record model.PostDNSRecord) (*model.PostDnsRecordsReps, error) {
	// Create the DNS record
	record.Comment = utils.FormatTime() + " " + record.Comment
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", config.Config.ZoneId)
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	// 使用 NewRequestWithBody 创建包含请求体的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Auth-Email", config.Config.ApiEmail)
	req.Header.Add("X-Auth-Key", config.Config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Unmarshal
	var newRecord model.PostDnsRecordsReps
	err = json.NewDecoder(resp.Body).Decode(&newRecord)
	if err != nil {
		return nil, err
	}

	return &newRecord, nil
}

func UpdateDNSRecord(record model.PostDNSRecord, id string) (*model.PostDnsRecordsReps, error) {
	// Update the DNS record
	record.Comment = utils.FormatTime() + " " + record.Comment

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", config.Config.ZoneId, id)
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Email", config.Config.ApiEmail)
	req.Header.Add("X-Auth-Key", config.Config.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// Unmarshal
	var newRecord model.PostDnsRecordsReps
	err = json.NewDecoder(resp.Body).Decode(&newRecord)
	if err != nil {
		return nil, err
	}
	return &newRecord, nil
}
