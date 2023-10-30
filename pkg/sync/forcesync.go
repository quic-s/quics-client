package sync

import (
	"log"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/types"
	qp "github.com/quic-s/quics-protocol"
	"github.com/quic-s/quics-protocol/pkg/stream"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func ForceSyncMain() {

	err := QPClient.RecvTransactionHandleFunc(qstypes.FORCESYNC, func(conn *qp.Connection, stream *stream.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client: [FORCESYNC] transaction start")
		UUID := badger.GetUUID()

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

		log.Println("quics-client: [FORCESYNC] transaction success")
		return nil
	})
	if err != nil {
		log.Println("quics-client: [FORCESYNC] ", err)
	}

}
