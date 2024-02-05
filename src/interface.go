package src

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger
var LogFile *os.File

func CreateLogger() (*log.Logger, *os.File, error) {
	//
	// Check for logs directory
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			fmt.Println("Error creating logs directory:", err)
			return nil, nil, err
		}
	}
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("logs/%s.log", currentDate)

	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return nil, nil, err
	}

	logger := log.New(io.MultiWriter(file, os.Stdout), "Cloudflare-DDNS: ", log.LstdFlags)
	return logger, file, nil
}

func StartServer() error {
	// Get the DNS records
	_, err := GetDnsRecords()
	if err != nil {
		return err
	}
	// Start the server
	for {
		err = CheckForUpdates()
		if err != nil {
			Logger.Println("Error checking for updates:", err)
		}
		// Sleep
		time.Sleep(time.Duration(Config.Duration) * time.Second)
	}
}

func CheckForUpdates() error {
	// Get the current IP

	// Get the DNS records
	records, err := GetDnsRecords()
	if err != nil {
		return err
	}

	// If exists but different, update
	// If not exists, create
	// If exists and the same, do nothing
	for _, config := range Config.DDNS {
		// Flag to check if a matching record is found
		foundMatchingRecord := false

		for _, record := range records.Result {
			if config.Name == record.Name && config.Type == record.Type {
				foundMatchingRecord = true

				var currentIP string

				if config.Type == "A" {
					currentIP, err = GetCurrentIpv4()
				} else if config.Type == "AAAA" {
					currentIP, err = GetCurrentIpv6()
				}
				if err != nil {
					return err
				}

				if currentIP != record.Content {
					// Update existing record
					newRecord := PostDNSRecord{
						Type:    config.Type,
						Name:    config.Name,
						Content: currentIP,
						TTL:     config.TTL,
						Proxied: config.Proxied,
						Comment: config.Comment,
					}
					resp, err := UpdateDNSRecord(newRecord, record.Id)
					if err != nil {
						return err
					}
					if !resp.Success {
						return errors.New("error updating the DNS record")
					}
					Logger.Println("DNSUpdated", config.Name, "->", currentIP)
				} else {
					Logger.Println("No change", config.Name, "->", currentIP)
				}
				// delete the var currentIP from the memory
				currentIP = ""

				break // Break the inner loop after finding a matching record
			}
		}

		// If no matching record is found, add a new one
		var currentIP string
		if !foundMatchingRecord {
			if config.Type == "A" {
				currentIP, err = GetCurrentIpv4()
			}
			if config.Type == "AAAA" {
				currentIP, err = GetCurrentIpv6()
			}
			if err != nil {
				return err
			}

			newRecord := PostDNSRecord{
				Type:    config.Type,
				Name:    config.Name,
				Content: currentIP,
				TTL:     config.TTL,
				Proxied: config.Proxied,
				Comment: config.Comment,
			}
			resp, err := CreateDNSRecord(newRecord)
			if err != nil {
				return err
			}
			if !resp.Success {
				Logger.Fatalln("Error adding the DNS record", resp.Errors)
			}
			Logger.Println("DNSAdded", config.Name, "->", currentIP)
		}
	}

	return nil
}

func GetContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("could not get the current IP, If your device is not connected to the internet or not support ipv6")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func GetCurrentIpv4() (string, error) {
	url := "https://v4.ident.me"
	return GetContent(url)
}

func GetCurrentIpv6() (string, error) {
	url := "https://v6.ident.me"
	return GetContent(url)
}
