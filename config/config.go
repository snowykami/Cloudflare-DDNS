package config

import (
	log "github.com/sirupsen/logrus"
	"os"

	"gopkg.in/yaml.v3"
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

var Config *BaseConfig

var KeyValue = make(map[string]any)

var (
	EnableIPv4 bool = false
	EnableIPv6 bool = false
)

// ddns切片
var AAAADDNSTasks = make([]DDNSConfig, 0)
var ADDNSTasks = make([]DDNSConfig, 0)

func NewConfig() BaseConfig {
	config := BaseConfig{
		ApiKey:   "your_api_key",
		ApiEmail: "user@example.com",
		Duration: 300,
		ZoneId:   "your_zone_id",
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

func Get(key string) string {
	if value, ok := KeyValue[key]; ok {
		return value.(string)
	} else {
		return ""
	}
}

func GetWithDefault(key string, defaultValue string) string {
	if value, ok := KeyValue[key]; ok {
		return value.(string)
	} else {
		return defaultValue
	}
}

func init() {
	msg, err := ReadConfig()
	if err != nil {
		log.Println(msg)
		os.Exit(1)
	}
	log.Println(msg)

	for _, ddnsTask := range Config.DDNS {
		if ddnsTask.Type == "A" {
			ADDNSTasks = append(ADDNSTasks, ddnsTask)
			EnableIPv4 = true
		} else if ddnsTask.Type == "AAAA" {
			AAAADDNSTasks = append(AAAADDNSTasks, ddnsTask)
			EnableIPv6 = true
		} else {
			log.Printf("Unsupported DNS record type: %s\n", ddnsTask.Type)
			os.Exit(1)
		}
	}

	if !EnableIPv4 && !EnableIPv6 {
		log.Println("No DNS record is enabled. 未启用任何DNS记录")
		os.Exit(1)
	}
}

func ReadConfig() (string, error) {
	_, err := os.Stat("config.yml")
	if os.IsNotExist(err) {
		// create config file
		file, err := os.Create("config.yml")
		if err != nil {
			return "Create config file error. 创建配置文件错误。", err
		}
		yamlData, err := yaml.Marshal(NewConfig())
		if err != nil {
			return "Marshal config error. 配置文件序列化错误。", err
		}
		_, err = file.Write(yamlData)
		if err != nil {
			return "Write config file error. 写入配置文件错误。", err
		}
		// end and tell user to edit config
		err = file.Close()
		if err != nil {
			return "Close config file error. 关闭配置文件错误。", err
		}
		log.Println("Configuration file has been created. Please edit it then enter to continue. 配置文件已经被创建，请编辑后按下回车继续。")
		if _, err = os.Stdin.Read(make([]byte, 1)); err != nil {
			return "Read from stdin error. 从标准输入读取错误。", err
		}
	}

	file, err := os.ReadFile("config.yml")
	if err != nil {
		return "Read config file error. 读取配置文件错误。", err
	}
	// Unmarshal
	err = yaml.Unmarshal(file, &Config)
	if err != nil {
		return "Unmarshal config file error. 配置文件反序列化错误。", err
	}
	err = yaml.Unmarshal(file, &KeyValue)
	if err != nil {
		return "Unmarshal config file error. 动态配置文件反序列化错误。", err
	}
	return "Read config file successfully. 读取配置文件成功。", nil
}
