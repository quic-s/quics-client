package quic

import (
	"crypto/tls"
	"log"
	"net"
	"strconv"
	"time"

	qp "github.com/quic-s/quics-protocol"
)

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

	err = conn.SendFileWithMessage(msgtype, message, filepath)
	if err != nil {
		log.Println(err)
	}

	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}
