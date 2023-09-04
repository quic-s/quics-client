package quic

import (
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"time"

	qp "github.com/quic-s/quics-protocol"
)

// message type
const (
	DOWNLOAD string = "DOWNLOAD"
	CREATE   string = "CREATE"
	DELETE   string = "DELETE"
	WRITE    string = "WRITE"
	RENAME   string = "RENAME"
	CLIENT   string = "CLIENT"
	ROOTDIR  string = "ROOTDIR" // /root/path/like/this
	RESCAN   string = "RESCAN"
	HISTORY  string = "HISTORY"
	SHARING  string = "SHARING"
	FILE     string = "FILE"
)

func ClientMessage(msgtype string, message []byte) {
	quicClient, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quics-protocol"},
	}
	// start client
	parsedPort, _ := strconv.Atoi(host)
	conn, err := quicClient.Dial(&net.UDPAddr{IP: net.IP(host), Port: parsedPort}, tlsConf)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	err = conn.SendMessage(msgtype, message)
	if err != nil {
		log.Println(err)
	}
	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}

// ex) ClientFile("/Users/username/Desktop/test.txt")
func ClientFile(filepath string) {

	quicClient, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quics-protocol"},
	}
	// start client
	parsedPort, _ := strconv.Atoi(host)
	conn, err := quicClient.Dial(&net.UDPAddr{IP: net.IP(host), Port: parsedPort}, tlsConf)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	err = conn.SendFile(FILE, filepath)
	if err != nil {
		log.Println(err)
	}

	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}

// ex) ClientFileWithMessage("/Users/username/Desktop/test.txt", "CREATE", []byte("hello"))
func ClientFileWithMessage(filepath string, msgtype string, message []byte) {
	quicClient, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quics-protocol"},
	}
	// start client
	parsedPort, _ := strconv.Atoi(host)
	conn, err := quicClient.Dial(&net.UDPAddr{IP: net.IP(host), Port: parsedPort}, tlsConf)
	if err != nil {
		log.Println("quics-protocol: ", err)
	}
	if err != nil {
		log.Println("quics-protocol: ", err)
	}

	err = conn.SendFileMessage(msgtype, message, filepath)
	if err != nil {
		log.Println(err)
	}

	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}
