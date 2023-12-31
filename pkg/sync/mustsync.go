package sync

import (
	"crypto/sha1"
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"

	"github.com/quic-s/quics-client/pkg/types"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func NeedContentMain() {
	err := QPClient.RecvTransactionHandleFunc(qstypes.NEEDCONTENT, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		req, err := qclient.NeedContentRecvHandler(stream)
		if err != nil {
			return err
		}
		// get paths
		afterPath := req.AfterPath
		beforePath := badger.GetBeforePathWithAfterPath(afterPath)
		path := filepath.Join(beforePath, afterPath)

		syncMeta := badger.GetSyncMetadata(path)
		if reflect.ValueOf(syncMeta).IsZero() {
			return fmt.Errorf("cannot find sync metadata")
		}

		if syncMeta.LastUpdateTimestamp == req.LastUpdateTimestamp && syncMeta.LastUpdateHash == req.LastUpdateHash {
			err := qclient.NeedContentHandler(stream, path, badger.GetUUID(), afterPath, syncMeta.LastUpdateTimestamp, syncMeta.LastUpdateHash)
			if err != nil {
				return err
			}
		}
		return nil

	})
	if err != nil {
		log.Println("[quics-client : NEEDCONTENT] ", err)
	}

}

// MustSyncMain is a function to handle MustSync transaction
func MustSyncMain() {

	err := QPClient.RecvTransactionHandleFunc(qstypes.MUSTSYNC, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		UUID := badger.GetUUID()
		mustSyncReq, err := qclient.MustSyncRecvHandler(stream)
		if err != nil {
			return err
		}
		// get paths
		afterPath := mustSyncReq.AfterPath
		beforePath := badger.GetBeforePathWithAfterPath(afterPath)

		// lock mutex for each path
		h := sha1.New()
		h.Write([]byte(afterPath))
		hash := h.Sum(nil)

		PSMut[uint8(hash[0]%PSMutModNum)].Lock()
		defer PSMut[uint8(hash[0]%PSMutModNum)].Unlock()

		path := filepath.Join(beforePath, afterPath)

		//If New File in coming, then make new sync meta in badger
		// e.g. new file created by other clients
		// if !badger.IsSyncMetadataExisted(path) {
		// 	syncMetadata := types.SyncMetadata{
		// 		BeforePath:          beforePath,
		// 		AfterPath:           afterPath,
		// 		LastUpdateTimestamp: 0,
		// 		LastUpdateHash:      "",
		// 		LastSyncTimestamp:   0,
		// 		LastSyncHash:        "",
		// 	}
		// 	err := badger.Update(path, syncMetadata.Encode())
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		syncMetadata := badger.GetSyncMetadata(path)

		if !CheckCanOverwrite(syncMetadata.LastSyncTimestamp, syncMetadata.LastUpdateTimestamp, mustSyncReq.LatestSyncTimestamp) {
			err := qclient.MustSyncHandler(stream, UUID, "", 0, "")
			if err != nil {
				return err
			}
			return fmt.Errorf("transaction fail, cannot overwrite")

		} else {
			err = qclient.MustSyncHandler(stream, UUID, afterPath, mustSyncReq.LatestSyncTimestamp, mustSyncReq.LatestHash)
			if err != nil {
				return err
			}
		}

		// Get File Contents
		IsRemoved := false

		// when "event remove" broadcasted, then do not save the file
		if mustSyncReq.LatestHash == "" {
			IsRemoved = true
		}
		_, err = qclient.GiveYouRecvHandler(stream, path, afterPath, mustSyncReq.LatestHash, IsRemoved)
		if err != nil {
			return err
		}

		if IsRemoved {
			badger.Delete(path)
		} else {
			//update syncmeta
			updatedSyncMeta := types.SyncMetadata{
				BeforePath:          beforePath,
				AfterPath:           afterPath,
				LastUpdateTimestamp: mustSyncReq.LatestSyncTimestamp,
				LastUpdateHash:      mustSyncReq.LatestHash,
				LastSyncTimestamp:   mustSyncReq.LatestSyncTimestamp,
				LastSyncHash:        mustSyncReq.LatestHash,
			}
			badger.Update(path, updatedSyncMeta.Encode())
		}

		err = qclient.GiveYouHandler(stream, UUID, afterPath, mustSyncReq.LatestSyncTimestamp, mustSyncReq.LatestHash)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("quics-client : [MUSTSYNC] ", err)
	}

}

// CheckCanOverwrite checks whether the file can be overwritten
func CheckCanOverwrite(LastSyncTimestamp uint64, LastUpdateTimestamp uint64, LastestSyncTimestamp uint64) bool {
	if LastSyncTimestamp == LastUpdateTimestamp && LastestSyncTimestamp > LastSyncTimestamp {
		return true
	}
	return false
}
