package utils

import (
	"fmt"
	"net"
	"os"
)

func GetIp() string {
	// 현재 IP 주소 중 첫 번째 것을 선택한다
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	ip := addrs[0].String()

	return ip
}
