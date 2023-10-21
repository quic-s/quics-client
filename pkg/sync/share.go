package sync

import (
	"fmt"
	"log"
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"
)

// @URL /api/v1/share/file
func GetShareLink(path string, MaxCnt uint) (string, error) {
	UUID := badger.GetUUID()
	_, AfterPath := badger.SplitBeforeAfterRoot(path)
	link := ""

	err := Conn.OpenTransaction("SHARE", func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		linkShareRes, err := qclient.SendLinkShare(stream, UUID, AfterPath, MaxCnt)
		if err != nil {
			return err
		}
		link = linkShareRes.Link

		log.Println("quics-client : [SHARE] transaction success")
		return nil

	})
	if err != nil {
		return "", err
	}

	return link, nil

}

// @URL /api/v1/share/stop
func StopShare(link string) error {

	err := Conn.OpenTransaction("SHARESTOP", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		stopShareRes, err := qclient.SendStopShare(stream, badger.GetUUID(), link)
		if err != nil {
			return err
		}
		if reflect.ValueOf(stopShareRes).IsZero() {
			return fmt.Errorf("[STOPSHARE] server cannot stop sharing")
		}
		log.Println("quics-client : [STOPSHARE] transaction success")
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
