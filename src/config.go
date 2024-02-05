package src

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type DDNSConfig struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	TTL     int    `yaml:"ttl"`
	Proxied bool   `yaml:"proxied"`
	Comment string `yaml:"comment"`
}

type BaseConfig struct {
	ApiKey   string       `yaml:"api_key"`
	ApiEmail string       `yaml:"api_email"`
	Duration int          `yaml:"duration"`
	ZoneId   string       `yaml:"zone_id"`
	DDNS     []DDNSConfig `yaml:"ddns"`
}

var Config BaseConfig

func NewConfig() BaseConfig {
	config := BaseConfig{
		ApiKey:   "xxxxxxxx",
		ApiEmail: "user@example.com",
		Duration: 300,
		ZoneId:   "xxxxxxxx",
		DDNS: []DDNSConfig{
			{
				Name:    "v4.example.com",
				Type:    "A",
				TTL:     60,
				Proxied: true,
				Comment: "DDNS",
			},
			{
				Name:    "v6.example.com",
				Type:    "AAAA",
				TTL:     60,
				Proxied: true,
				Comment: "DDNS",
			},
		},
	}
	return config
}

func ReadConfig() (bool, error) {
	// Detect file exist
	// Read file
	// If not exist, Create file and write default
	_, err := os.Stat("config.yml")
	if os.IsNotExist(err) {
		file, err := os.Create("config.yml")
		if err != nil {
			return false, err
		}
		yamlData, err := yaml.Marshal(NewConfig())
		if err != nil {
			return false, err
		}
		_, err = file.WriteString(string(yamlData))
		// end and tell user to edit config
		err = file.Close()
		if err != nil {
			return false, err
		}
		fmt.Println("Configuration file has been created. Please edit it and run the program again.")
		os.Exit(0)
	}
	file, err := os.ReadFile("config.yml")
	if err != nil {
		return false, err
	}
	// Unmarshal
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		return false, err
	}
	return true, nil
}
