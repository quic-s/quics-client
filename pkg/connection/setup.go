package connection

import (
	"io"
	"log"
	"os"
	"sync"

	qp "github.com/quic-s/quics-protocol"
)

var (
	wg                   = sync.WaitGroup{}
	QuicResponseListener *qp.QP
)

func init() {
	var err error
	QuicResponseListener, err = initializeServer()
	if err != nil {
		log.Fatal(err)
	}

}

func initializeServer() (*qp.QP, error) {
	// initialize server
	quicProto, err := qp.New(qp.LOG_LEVEL_INFO)
	if err != nil {
		return nil, err
	}
	err = quicProto.RecvMessageHandleFunc("test", func(conn *qp.Connection, msgType string, data []byte) {
		defer wg.Done()
		log.Println("quics-protocol: ", "message received ", conn.Conn.RemoteAddr().String())
		log.Println("quics-protocol: ", msgType, string(data))
	})
	if err != nil {
		return nil, err
	}

	err = quicProto.RecvFileHandleFunc("test", func(conn *qp.Connection, fileType string, fileInfo *qp.FileInfo, fileReader io.Reader) {
		defer wg.Done()
		log.Println("quics-protocol: ", "file received ", fileInfo.Name)
		file, err := os.Create("received.txt")
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(file, fileReader)
		if err != nil {
			log.Fatal(err)
		}
		if n != fileInfo.Size {
			log.Fatalf("quics-protocol: read only %dbytes", n)
		}
		log.Println("quics-protocol: ", "file saved with ", n, "bytes")
	})
	if err != nil {
		return nil, err
	}

	err = quicProto.RecvFileMessageHandleFunc("test", func(conn *qp.Connection, fileMsgType string, data []byte, fileInfo *qp.FileInfo, fileReader io.Reader) {
		defer wg.Done()
		log.Println("quics-protocol: ", "message received ", conn.Conn.RemoteAddr().String())
		log.Println("quics-protocol: ", fileMsgType, string(data))

		log.Println("quics-protocol: ", "file received ", fileInfo.Name)
		file, err := os.Create("received2.txt")
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(file, fileReader)
		if err != nil {
			log.Fatal(err)
		}
		if n != fileInfo.Size {
			log.Println("quics-protocol: ", "read only ", n, "bytes")
			log.Fatal(err)
		}
		log.Println("quics-protocol: ", "file saved")
	})
	if err != nil {
		return nil, err
	}

	return quicProto, nil
}
