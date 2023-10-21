package sync

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"

	"github.com/quic-s/quics-client/pkg/types"
	qp "github.com/quic-s/quics-protocol"
)

// TODO
func NeedContentMain() {

	err := QPClient.RecvTransactionHandleFunc("NEEDCONTENT", func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
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

		if syncMeta.LastUpdateTimestamp == req.LatestUpdateTimestamp && syncMeta.LastUpdateHash == req.LatestUpdateHash {
			err := qclient.NeedContentHandler(stream, badger.GetUUID(), afterPath, syncMeta.LastUpdateTimestamp, syncMeta.LastUpdateHash)
			if err != nil {
				return err
			}
		}
		return nil

	})
	if err != nil {
		log.Println("[NEEDCONTENT] ", err)
	}

}

func MustSyncMain() {

	err := QPClient.RecvTransactionHandleFunc("MUSTSYNC", func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		UUID := badger.GetUUID()
		mustSyncReq, err := qclient.MustSyncRecvHandler(stream)
		if err != nil {
			return err
		}
		// get paths
		afterPath := mustSyncReq.AfterPath
		beforePath := badger.GetBeforePathWithAfterPath(afterPath)
		path := filepath.Join(beforePath, afterPath)
		log.Println("before path >> ", beforePath)
		log.Println("after path >> ", afterPath)
		log.Println("path >> ", path)

		//If New File in coming, then make new sync meta in badger
		// e.g. new file created by other clients
		if !badger.IsSyncMetadataExisted(path) {
			syncMetadata := types.SyncMetadata{
				BeforePath:          beforePath,
				AfterPath:           afterPath,
				LastUpdateTimestamp: 0,
				LastUpdateHash:      "",
				LastSyncTimestamp:   0,
				LastSyncHash:        "",
			}
			err := badger.Update(path, syncMetadata.Encode())
			if err != nil {
				return err
			}
		}

		syncMetadata := badger.GetSyncMetadata(path)
		log.Println("syncMetadata >> ", syncMetadata)

		if !CheckCanOverwrite(syncMetadata.LastSyncTimestamp, syncMetadata.LastUpdateTimestamp, mustSyncReq.LatestSyncTimestamp) {
			err := qclient.MustSyncHandler(stream, UUID, "", 0, "")
			if err != nil {
				log.Println("can not response mustsync >> ")
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
		_, err = qclient.GiveYouRecvHandler(stream, path, IsRemoved)
		if err != nil {
			return err
		}

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

		if IsRemoved {
			badger.Delete(path)

		}

		err = qclient.GiveYouHandler(stream, UUID, afterPath, mustSyncReq.LatestSyncTimestamp, mustSyncReq.LatestHash)
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
