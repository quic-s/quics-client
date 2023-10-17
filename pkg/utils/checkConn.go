package utils

import (
	"net"
	"time"
)

// CheckInternetConnection checks internet connection
func CheckInternetConnection() bool {

	conn, err := net.DialTimeout("tcp", "google.com:80", time.Second)
	if err != nil {

		return false
	}

	defer conn.Close()
	return true
}
