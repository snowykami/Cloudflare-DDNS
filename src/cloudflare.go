package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DNSRecord struct {
	Content   string `json:"content"`
	Name      string `json:"name"`
	Proxied   bool   `json:"proxied"`
	Type      string `json:"type"`
	Comment   string `json:"comment"`
	CreatedOn string `json:"created_on"`
	Id        string `json:"id"`
	Locked    bool   `json:"locked"`
	Meta      struct {
		AutoAdded bool   `json:"auto_added"`
		Source    string `json:"source"`
	}
	ModifiedOn string `json:"modified_on"`
	Proxiable  bool   `json:"proxiable"`
	Tags       []string
	TTL        int    `json:"ttl"`
	ZoneId     string `json:"zone_id"`
	ZoneName   string `json:"zone_name"`
}

type PostDNSRecord struct {
	Content string   `json:"content"`
	Name    string   `json:"name"`
	Proxied bool     `json:"proxied"`
	Type    string   `json:"type"`
	TTL     int      `json:"ttl"`
	Comment string   `json:"comment"`
	Tags    []string `json:"tags"`
}

// Use for response

type TypeGetDnsRecords struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   []DNSRecord   `json:"result"`
}

type TypePostDnsRecords struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []string      `json:"messages"`
	Result   DNSRecord     `json:"result"`
}

func GetDnsRecords() (TypeGetDnsRecords, error) {
	// Get the DNS records
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", Config.ZoneId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TypeGetDnsRecords{}, err
	}
	req.Header.Add("X-Auth-Email", Config.ApiEmail)
	req.Header.Add("X-Auth-Key", Config.ApiKey)
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TypeGetDnsRecords{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// Unmarshal
	var records TypeGetDnsRecords
	err = json.NewDecoder(resp.Body).Decode(&records)
	if err != nil {
		return TypeGetDnsRecords{}, err
	}
	return records, nil
}

func CreateDNSRecord(record PostDNSRecord) (TypePostDnsRecords, error) {
	// Create the DNS record
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", Config.ZoneId)
	data, err := json.Marshal(record)
	if err != nil {
		return TypePostDnsRecords{}, err
	}

	// 使用 NewRequestWithBody 创建包含请求体的请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return TypePostDnsRecords{}, err
	}

	req.Header.Add("X-Auth-Email", Config.ApiEmail)
	req.Header.Add("X-Auth-Key", Config.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TypePostDnsRecords{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Unmarshal
	var newRecord TypePostDnsRecords
	err = json.NewDecoder(resp.Body).Decode(&newRecord)
	if err != nil {
		return TypePostDnsRecords{}, err
	}

	return newRecord, nil
}

func UpdateDNSRecord(record PostDNSRecord, id string) (TypePostDnsRecords, error) {
	// Update the DNS record
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", Config.ZoneId, id)
	data, err := json.Marshal(record)
	if err != nil {
		return TypePostDnsRecords{}, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return TypePostDnsRecords{}, err
	}
	req.Header.Add("X-Auth-Email", Config.ApiEmail)
	req.Header.Add("X-Auth-Key", Config.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TypePostDnsRecords{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// Unmarshal
	var newRecord TypePostDnsRecords
	err = json.NewDecoder(resp.Body).Decode(&newRecord)
	if err != nil {
		return TypePostDnsRecords{}, err
	}
	return newRecord, nil
}
