package utils

import (
	"io/ioutil"
	"net/http"
)

// 从指定的url获取文本内容
func GetContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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