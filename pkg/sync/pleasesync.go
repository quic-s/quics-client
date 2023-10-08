package sync

import (
	"fmt"
	"log"
	"os"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func PSwhenWrite(path string, info os.FileInfo) {
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "WRITE"
	localModTime := info.ModTime().String()

	// Get PrevSyncMetadata
	prevSyncMetaByte, err := badger.View(path)
	if err != nil {
		log.Println(err)
	}
	prevSyncMetadata := types.SyncMetadata{}
	prevSyncMetadata.Decode(prevSyncMetaByte)

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, info), // make new hash
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
		Conflict:            prevSyncMetadata.Conflict,
	}
	badger.Update(path, syncMetadata.Encode())

	// PleaseSync Transaction
	Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client : [PLEASESYNC] transaction start")

		// Get FileMeta before PleaseSync
		fileMetaRes, err := qclient.SendFileMeta(stream, UUID, AfterPath)
		if err != nil {
			log.Println("quics-client : ", err)

			return err
		}

		// Check condition conflict before pleaseSync
		if IsConfilcted(fileMetaRes.LatestSyncTimestamp, syncMetadata.LastUpdateTimestamp, fileMetaRes.LatestHash, syncMetadata.LastSyncHash) {
			err := badger.AddConflictAndConflictFileList(path, types.ConflictMetadata{
				ServerModDate:   fileMetaRes.ModifiedDate,
				ServerDevice:    "",
				ServerTimestamp: fileMetaRes.LatestSyncTimestamp,
				ServerHash:      fileMetaRes.LatestHash,

				LocalModDate:   localModTime,
				LocalDevice:    "",
				LocalTimestamp: syncMetadata.LastUpdateTimestamp,
				LocalHash:      syncMetadata.LastUpdateHash,
			})
			if err != nil {
				return err
			}
			return fmt.Errorf("quics-client : [PLEASESYNC] transaction fail")
		}

		// Send PleaseSync
		_, err = qclient.SendPleaseSync(stream, UUID, event, BeforePath, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash)
		if err != nil {
			return err
		}

		// Update Sync Timestamp and hash as same as update Timestamp and hash
		syncMetadata = types.SyncMetadata{
			BeforePath:          BeforePath,
			AfterPath:           AfterPath,
			LastUpdateTimestamp: syncMetadata.LastUpdateTimestamp,
			LastUpdateHash:      syncMetadata.LastUpdateHash,
			LastSyncTimestamp:   syncMetadata.LastUpdateTimestamp,
			LastSyncHash:        syncMetadata.LastUpdateHash,
			Conflict:            syncMetadata.Conflict,
		}
		badger.Update(path, syncMetadata.Encode())

		// Send FileData
		_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
		if err != nil {
			return err
		}

		log.Println("quics-client : [PLEASESYNC] transaction success")
		return nil
	})
}

func PSwhenCreate(path string, info os.FileInfo) {
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "CREATE"
	modDate := info.ModTime().String()
	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, info), // make new hash
		LastSyncTimestamp:   0,
		LastSyncHash:        "",
		Conflict:            types.ConflictMetadata{},
	}
	badger.Update(path, syncMetadata.Encode())
	Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] transaction start")

		// Get FileMeta before PleaseSync
		fileMetaRes, err := qclient.SendFileMeta(stream, UUID, AfterPath)
		if err != nil {
			log.Println("quics-client : ", err)
			return err
		}

		// Check condition to overwrite
		if IsConfilcted(fileMetaRes.LatestSyncTimestamp, syncMetadata.LastUpdateTimestamp, fileMetaRes.LatestHash, syncMetadata.LastSyncHash) {

			badger.AddConflictAndConflictFileList(path, types.ConflictMetadata{
				ServerModDate:   fileMetaRes.ModifiedDate,
				ServerDevice:    "",
				ServerTimestamp: fileMetaRes.LatestSyncTimestamp,
				ServerHash:      fileMetaRes.LatestHash,

				LocalModDate:   modDate,
				LocalDevice:    "",
				LocalTimestamp: syncMetadata.LastUpdateTimestamp,
				LocalHash:      syncMetadata.LastUpdateHash,
			})
			return fmt.Errorf("quics-client : [PLEASESYNC] transaction fail")
		}

		// Send PleaseSync
		_, err = qclient.SendPleaseSync(stream, UUID, event, BeforePath, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash)
		if err != nil {
			return err
		}

		// Update Sync Timestamp and hash as same as update Timestamp and hash
		syncMetadata = types.SyncMetadata{
			BeforePath:          BeforePath,
			AfterPath:           AfterPath,
			LastUpdateTimestamp: syncMetadata.LastUpdateTimestamp,
			LastUpdateHash:      syncMetadata.LastUpdateHash,
			LastSyncTimestamp:   syncMetadata.LastUpdateTimestamp,
			LastSyncHash:        syncMetadata.LastUpdateHash,
			Conflict:            syncMetadata.Conflict,
		}
		badger.Update(path, syncMetadata.Encode())

		// Send FileData
		_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
		if err != nil {
			return err
		}

		log.Println("quics-client : [PLEASESYNC] transaction success")
		return nil

	})
}

func PSwhenRemove(path string) {
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "REMOVE"

	prevSyncMetaByte, err := badger.View(path)
	if err != nil {
		log.Println(err)
	}
	prevSyncMetadata := types.SyncMetadata{}
	prevSyncMetadata.Decode(prevSyncMetaByte)

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      "",
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
		Conflict:            prevSyncMetadata.Conflict,
	}
	badger.Update(path, syncMetadata.Encode())

	err = Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] transaction start")

		// Get FileMeta before PleaseSync
		fileMetaRes, err := qclient.SendFileMeta(stream, UUID, AfterPath)
		if err != nil {
			return err
		}

		// Check condition conflict
		if !IsConfilcted(fileMetaRes.LatestSyncTimestamp, syncMetadata.LastUpdateTimestamp, fileMetaRes.LatestHash, syncMetadata.LastSyncHash) {

			badger.AddConflictAndConflictFileList(path, types.ConflictMetadata{
				ServerModDate:   fileMetaRes.ModifiedDate,
				ServerDevice:    "",
				ServerTimestamp: fileMetaRes.LatestSyncTimestamp,
				ServerHash:      fileMetaRes.LatestHash,

				LocalModDate:   "Removed, No ModDate",
				LocalDevice:    "",
				LocalTimestamp: syncMetadata.LastUpdateTimestamp,
				LocalHash:      syncMetadata.LastUpdateHash,
			})

			return fmt.Errorf("Cannot remove file because of conflict")
		}

		// Send PleaseSync
		_, err = qclient.SendPleaseSync(stream, UUID, event, BeforePath, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash)
		if err != nil {
			return err
		}

		// Delete SyncMetadata when event is REMOVE
		badger.Delete(path)

		log.Println("quics-client : [PLEASESYNC] transaction success")
		return nil

	})
	if err != nil {
		log.Println("quics-client : [PLEASESYNC] >> Remove >>", err)
	}

}

func IsConfilcted(LastestSyncTimestamp uint64, LastUpdateTimestamp uint64, LastestSyncHash string, LastSyncHash string) bool {
	if LastestSyncTimestamp < LastUpdateTimestamp && LastestSyncHash == LastSyncHash {
		return true
	}
	return false
}
