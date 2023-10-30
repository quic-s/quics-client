package sync

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func PleaseSync(path string) {
	prevInfo, _ := os.Stat(path)
	prevSize := int64(0)
	if prevInfo != nil {
		prevSize = prevInfo.Size()
	}
	time.Sleep(100 * time.Millisecond)
	//log.Println("quics-client : [PLEASESYNC] PleaseSync start", path)

	h := sha1.New()
	h.Write([]byte(path))
	hash := h.Sum(nil)

	PSMut[uint8(hash[0]%PSMutModNum)].Lock()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
		PSMut[uint8(hash[0]%PSMutModNum)].Unlock()

	}()

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		// when file is removed
		PSwhenRemove(path)
		return
	} else if err != nil {
		log.Println("quics-client : ", err)
		return
	}

	if prevSize != fileInfo.Size() {
		// when file is under writing
		return
	}
	if fileInfo.IsDir() {
		return
	}

	if !badger.IsSyncMetadataExisted(path) {
		// when file is created
		PSwhenCreate(path, fileInfo)
		return
	}

	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)

	// Get PrevSyncMetadata
	prevSyncMetadata := badger.GetSyncMetadata(path)

	// update syncMeta for events happened
	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, fileInfo), // make new hash
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
	}

	// if LastUpdateHash is same, then return
	if CanReturnPSByMS(&prevSyncMetadata, &syncMetadata) {
		return
	}

	PSwhenWrite(path, fileInfo, syncMetadata)
}

func CanReturnPSByMS(prevSyncMeta *types.SyncMetadata, currSyncMeta *types.SyncMetadata) bool {
	return prevSyncMeta.LastUpdateHash == currSyncMeta.LastUpdateHash
}

func PSwhenWrite(path string, info os.FileInfo, syncMetadata types.SyncMetadata) {
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "WRITE"

	badger.Update(path, syncMetadata.Encode())

	// PleaseSync Transaction
	err := Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client : [PLEASESYNC] write transaction start")

		serverfilemeta := qstypes.FileMetadata{}
		serverfilemeta.DecodeFromOSFileInfo(info)

		// Send PleaseSync
		res, err := qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
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
		}
		badger.Update(path, syncMetadata.Encode())

		if res.Status == "GIVEME" {
			// Send FileData
			_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("file already exited in server")
		}

		return nil
	})
	if err != nil {
		log.Println("quics-client : [PLEASESYNC] write transaction failed :", err)
	}
	log.Println("quics-client : [PLEASESYNC] write transaction success")
}

func PSwhenCreate(path string, info os.FileInfo) {
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	log.Println(AfterPath)
	UUID := badger.GetUUID()
	event := "CREATE"

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, info), // make new hash
		LastSyncTimestamp:   0,
		LastSyncHash:        "",
	}
	badger.Update(path, syncMetadata.Encode())
	err := Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] create transaction start")

		// Send PleaseSync

		serverfilemeta := qstypes.FileMetadata{}

		serverfilemeta.DecodeFromOSFileInfo(info)
		// Send PleaseSync

		res, err := qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
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
		}
		badger.Update(path, syncMetadata.Encode())

		// Send FileData
		if res.Status == "GIVEME" {
			// Send FileData
			_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("file already exited in server")
		}

		return nil

	})
	if err != nil {
		log.Println("quics-client : [PLEASESYNC] create transaction failed", err)
	}
	log.Println("quics-client : [PLEASESYNC] create transaction success")
}

func PSwhenRemove(path string) {
	//pre-request
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "REMOVE"

	// Update Sync Timestamp and hash as same as update Timestamp and hash
	prevSyncMetadata := badger.GetSyncMetadata(path)
	if reflect.ValueOf(prevSyncMetadata).IsZero() {
		return
	}

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      "",
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
	}
	badger.Update(path, syncMetadata.Encode())

	err := Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] remove transaction start")

		// Send PleaseSync

		serverfilemeta := qstypes.FileMetadata{
			ModTime: time.Now(),
		}

		// Send PleaseSync
		res, err := qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
		if err != nil {
			return err
		}
		// Delete SyncMetadata when event is REMOVE
		badger.Delete(path)
		// Send FileData
		if res.Status == "GIVEME" {

			_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, utils.GetEmptyFilePath())
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("file already exited in server")
		}

		log.Println("quics-client : [PLEASESYNC] remove transaction success")
		return nil

	})
	if err != nil {
		log.Println("quics-client : [PLEASESYNC] remove transaction failed", err)
	}

}
