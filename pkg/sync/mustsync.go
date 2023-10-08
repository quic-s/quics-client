package sync

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"

	"github.com/quic-s/quics-client/pkg/types"
	qp "github.com/quic-s/quics-protocol"
)

func MustSyncMain() {
	UUID := badger.GetUUID()

	err := QPClient.RecvTransactionHandleFunc("MUSTSYNC", func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		mustSyncReq, err := qclient.MustSyncRecvHandler(stream)
		if err != nil {
			return err
		}
		// get paths
		afterPath := mustSyncReq.AfterPath
		beforePath := badger.GetBeforePathWithAfterPath(afterPath)
		path := filepath.Join(beforePath, afterPath)

		// If New File in coming, then make new sync meta in badger
		if !badger.IsSyncMetadataExisted(path) {
			syncMetadata := types.SyncMetadata{
				BeforePath:          beforePath,
				AfterPath:           afterPath,
				LastUpdateTimestamp: mustSyncReq.LatestSyncTimestamp,
				LastUpdateHash:      mustSyncReq.LatestHash,
				LastSyncTimestamp:   mustSyncReq.LatestSyncTimestamp,
				LastSyncHash:        mustSyncReq.LatestHash,
				Conflict:            types.ConflictMetadata{},
			}
			err := badger.Update(path, syncMetadata.Encode())
			if err != nil {
				return err
			}
		}

		syncMetadata := badger.GetSyncMetadata(path)

		if !CheckCanOverwrite(syncMetadata.LastSyncTimestamp, syncMetadata.LastUpdateTimestamp, mustSyncReq.LatestSyncTimestamp) {
			err := qclient.MustSyncHandler(stream, UUID, "", 0, "")
			if err != nil {
				return err
			}
			return fmt.Errorf("transaction fail, cannot overwrite")

		}
		err = qclient.MustSyncHandler(stream, UUID, afterPath, syncMetadata.LastSyncTimestamp, syncMetadata.LastSyncHash)
		if err != nil {
			return err
		}

		_, err = qclient.GiveYouRecvHandler(stream, path)
		if err != nil {
			return err
		}

		//update syncmeta
		updatedSyncMeta := types.SyncMetadata{
			BeforePath:          beforePath,
			AfterPath:           afterPath,
			LastUpdateTimestamp: syncMetadata.LastUpdateTimestamp,
			LastUpdateHash:      syncMetadata.LastUpdateHash,
			LastSyncTimestamp:   mustSyncReq.LatestSyncTimestamp,
			LastSyncHash:        mustSyncReq.LatestHash,
		}
		badger.Update(path, updatedSyncMeta.Encode())

		err = qclient.GiveYouHandler(stream, UUID, afterPath, syncMetadata.LastSyncTimestamp, syncMetadata.LastSyncHash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("[MUSTSYNC] ", err)
	}

}

func CheckCanOverwrite(LastSyncTimestamp uint64, LastUpdateTimestamp uint64, LastestSyncTimestamp uint64) bool {
	if LastSyncTimestamp == LastUpdateTimestamp && LastestSyncTimestamp > LastSyncTimestamp {
		return true
	}
	return false
}
