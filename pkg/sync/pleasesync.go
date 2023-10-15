package sync

import (
	"crypto/sha1"
	"log"
	"os"
	"time"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-protocol/pkg/stream"

	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func CanReturnPSByMS(prevSyncMeta *types.SyncMetadata, currSyncMeta *types.SyncMetadata) bool {

	return prevSyncMeta.LastUpdateHash == currSyncMeta.LastUpdateHash

}

func PSwhenWrite(path string) {
	h := sha1.New()
	h.Write([]byte(path))
	hash := h.Sum(nil)

	PSMut[uint8(hash[0]%16)].Lock()
	defer PSMut[uint8(hash[0]%16)].Unlock()

	// pre request
	time.Sleep(50 * time.Millisecond)
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "WRITE"

	// Get PrevSyncMetadata
	prevSyncMetadata := badger.GetSyncMetadata(path)
	if prevSyncMetadata.LastSyncTimestamp == 0 {
		return
	}

	// update syncMeta for events happened
	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, info), // make new hash
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
	}

	if CanReturnPSByMS(&prevSyncMetadata, &syncMetadata) {
		return
	}

	badger.Update(path, syncMetadata.Encode())

	// PleaseSync Transaction
	Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *qp.Stream, transactionName string, transactionID []byte) error {
		log.Println("quics-client : [PLEASESYNC] transaction start")

		serverfilemeta := qstypes.FileMetadata{}
		serverfilemeta.DecodeFromOSFileInfo(info)

		// Send PleaseSync
		_, err = qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
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
		_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
		if err != nil {
			return err
		}

		log.Println("quics-client : [PLEASESYNC] transaction success")
		return nil
	})
}

func PSwhenCreate(path string) {
	h := sha1.New()
	h.Write([]byte(path))
	hash := h.Sum(nil)

	PSMut[uint8(hash[0]%16)].Lock()
	defer PSMut[uint8(hash[0]%16)].Unlock()

	//pre-requests
	time.Sleep(50 * time.Millisecond)
	info, err := os.Stat(path)
	if err != nil {
		log.Println("quics-client : ", err)
		return
	}
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	log.Println(AfterPath)
	UUID := badger.GetUUID()
	event := "CREATE"

	if badger.IsSyncMetadataExisted(path) {
		return
	}

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: 1,
		LastUpdateHash:      utils.MakeHash(AfterPath, info), // make new hash
		LastSyncTimestamp:   0,
		LastSyncHash:        "",
	}
	badger.Update(path, syncMetadata.Encode())
	Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] transaction start")

		// Send PleaseSync

		serverfilemeta := qstypes.FileMetadata{}
		serverfilemeta.DecodeFromOSFileInfo(info)
		// Send PleaseSync
		_, err := qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
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
		_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, path)
		if err != nil {
			return err
		}

		log.Println("quics-client : [PLEASESYNC] transaction success")
		return nil

	})
}

func PSwhenRemove(path string) {
	h := sha1.New()
	h.Write([]byte(path))
	hash := h.Sum(nil)

	PSMut[uint8(hash[0]%16)].Lock()
	defer PSMut[uint8(hash[0]%16)].Unlock()

	//pre-request
	BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)
	UUID := badger.GetUUID()
	event := "REMOVE"
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		return
	}

	if !badger.IsSyncMetadataExisted(path) {
		return
	}

	// Update Sync Timestamp and hash as same as update Timestamp and hash
	prevSyncMetadata := badger.GetSyncMetadata(path)

	syncMetadata := types.SyncMetadata{
		BeforePath:          BeforePath,
		AfterPath:           AfterPath,
		LastUpdateTimestamp: prevSyncMetadata.LastUpdateTimestamp + 1,
		LastUpdateHash:      "",
		LastSyncTimestamp:   prevSyncMetadata.LastSyncTimestamp,
		LastSyncHash:        prevSyncMetadata.LastSyncHash,
	}
	badger.Update(path, syncMetadata.Encode())

	err = Conn.OpenTransaction(qstypes.PLEASESYNC, func(stream *stream.Stream, transactionName string, transactionID []byte) error {

		log.Println("quics-client : [PLEASESYNC] transaction start")

		// Send PleaseSync

		serverfilemeta := qstypes.FileMetadata{
			ModTime: time.Now(),
		}

		// Send PleaseSync
		_, err := qclient.SendPleaseSync(stream, UUID, event, AfterPath, syncMetadata.LastUpdateTimestamp, syncMetadata.LastUpdateHash, syncMetadata.LastSyncHash, serverfilemeta)
		if err != nil {
			return err
		}

		// Send FileData
		_, err = qclient.SendPleaseTake(stream, UUID, AfterPath, utils.GetEmptyFilePath())
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
