package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func FullScanMain() {

	err := QPClient.RecvTransactionHandleFunc(qstypes.FULLSCAN, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		askAllMetaReq, err := qclient.AskAllMetaRecvHandler(stream)
		if err != nil {
			return err
		}
		if askAllMetaReq.UUID == "" {
			return fmt.Errorf(" AskAllMetaRecvHandler : UUID is empty ")
		}
		UUID := badger.GetUUID()
		if askAllMetaReq.UUID != UUID {
			return fmt.Errorf(" AskAllMetaRecvHandler : UUID is not same ")
		}

		rawlist, err := badger.GetAllSyncMetadataAmongRoot()
		if err != nil {
			return err
		}

		// Compare to file.Info and go ps thread to update syncMetadata as file.Info says
		resultList := []qstypes.SyncMetadata{}
		for _, item := range rawlist {
			path := filepath.Join(item.BeforePath, item.AfterPath)
			info, err := os.Stat(path)
			if err != nil {
				continue
			}
			hashtocompare := utils.MakeHash(item.AfterPath, info)
			if item.LastUpdateHash != hashtocompare {
				go PSwhenWrite(path)
			}
			convertedItem := qstypes.SyncMetadata{
				BeforePath:          item.BeforePath,
				AfterPath:           item.AfterPath,
				LastUpdateHash:      item.LastUpdateHash,
				LastUpdateTimestamp: item.LastUpdateTimestamp,
				LastSyncTimestamp:   item.LastSyncTimestamp,
				LastSyncHash:        item.LastSyncHash,
			}

			resultList = append(resultList, convertedItem)
		}

		// return askAllMetaHandler
		err = qclient.AskAllMetaHandler(stream, askAllMetaReq.UUID, resultList)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println("[FULLSCAN]: ", err)
	}
}
