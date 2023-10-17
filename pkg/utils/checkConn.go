package utils

import (
	"log"
	"net"
	"time"
)

func CheckInternetConnection() bool {
	// 1초의 타임아웃을 가진 TCP 연결을 생성합니다.
	conn, err := net.DialTimeout("tcp", "google.com:80", time.Second)
	if err != nil {
		log.Println("오프라인")
		return false
	}
	log.Println("온라인")
	// 연결을 닫습니다.
	defer conn.Close()
	return true
}
