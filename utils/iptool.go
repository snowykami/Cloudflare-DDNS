package utils

import (
	"io"
	"net/http"
	"time"
)

// 从指定的url获取文本内容
func GetContent(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetCurrentIpv4() (string, error) {
	url := "https://v4.ident.me"
	return GetContent(url)
}

func GetCurrentIpv6() (string, error) {
	url := "https://v6.ident.me"
	return GetContent(url)
}
