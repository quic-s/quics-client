package quic

import (
	"log"
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
)

type CRDBody struct {
	Uuid     string `json:"uuid"`
	FilePath string `json:"filepath"`
}

type ClientBody struct {
	Ip string `json:"ip"`
}

type RootdirBody struct {
	Uuid     string `json:"uuid"`
	RootPath string `json:"rootpath"`
}

func ClientMessage(msgtype string, message []byte) {
	// initialize client
	quicClient, err := qp.New()
	if err != nil {
		log.Panicln(err)
	}

	// start quics client
	err = quicClient.Dial(host + ":" + port)
	if err != nil {
		log.Panicln(err)
	}

	quicClient.SendMessage(msgtype, message)

	// delay for waiting message sent to server
	time.Sleep(3 * time.Second)
	quicClient.Close()
}
