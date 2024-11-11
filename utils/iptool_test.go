package utils

import "testing"

func TestGetContent(t *testing.T) {
	url := "https://v4.ident.me"
	_, err := GetContent(url)
	if err != nil {
		t.Error(err)
	}
}

func TestGetCurrentIpv4(t *testing.T) {
	ip, err := GetCurrentIpv4()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("IPv4 is supported and the current IP is", ip)
	}
}

func TestGetCurrentIpv6(t *testing.T) {
	ip, err := GetCurrentIpv6()
	if err != nil {
		t.Error(err)
	} else {
		t.Log("IPv6 is supported and the current IP is", ip)
	}
}
