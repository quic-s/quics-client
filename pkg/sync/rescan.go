package sync

import (
	"fmt"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"
)

// @URL /api/v1/rescan
// ex) Rescan()
func Rescan() error {

	rootdirlist := badger.GetRootDirList()
	rootdirafterpathlist := []string{}
	for _, rootdir := range rootdirlist {
		rootdirafterpathlist = append(rootdirafterpathlist, rootdir.AfterPath)
	}
	UUID := badger.GetUUID()

	err := Conn.OpenTransaction("RESCAN", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		rescanRes, err := qclient.SendRescan(stream, UUID, rootdirafterpathlist)
		if err != nil {
			return err
		}
		if rescanRes.UUID == "" {
			return fmt.Errorf("Server cannot rescan client")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
