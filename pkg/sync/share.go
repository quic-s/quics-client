package sync

import (
	"log"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"
)

func GetShareLink(path string, MaxCnt uint) (string, error) {
	UUID := badger.GetUUID()
	_, AfterPath := badger.SplitBeforeAfterRoot(path)
	link := ""

	err := Conn.OpenTransaction("SHARE", func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		// Send PleaseSync
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
