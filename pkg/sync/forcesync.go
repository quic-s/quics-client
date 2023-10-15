package sync

import (
	"log"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-protocol/pkg/stream"
)

func ForceSyncMain() {

	UUID := badger.GetUUID()
	err := Conn.OpenTransaction("FORCESYNC", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client: [FORCESYNC] transaction start")

		// Recv ForceSync
		req, Beforepath, err := qclient.ForceSyncRecvHandler(stream)
		if err != nil {
			return err
		}
		path := filepath.Join(Beforepath, req.AfterPath)
		//sync meta update
		syncMetaForcelyUpdate := types.SyncMetadata{
			BeforePath:          Beforepath,
			AfterPath:           req.AfterPath,
			LastSyncTimestamp:   req.LatestSyncTimestamp,
			LastSyncHash:        req.LatestHash,
			LastUpdateTimestamp: req.LatestSyncTimestamp,
			LastUpdateHash:      req.LatestHash,
		}
		err = badger.Update(path, syncMetaForcelyUpdate.Encode())
		if err != nil {
			return err
		}
		// Send ForceSyncRes
		err = qclient.ForceSyncHandler(stream, UUID, req.AfterPath, req.LatestSyncTimestamp, req.LatestHash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics-client: ", err)
	}

}
