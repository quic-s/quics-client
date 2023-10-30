package sync

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
	"github.com/quic-s/quics-client/pkg/net/qclient"
	"github.com/quic-s/quics-client/pkg/types"
	"github.com/quic-s/quics-client/pkg/utils"
	qp "github.com/quic-s/quics-protocol"
	qstypes "github.com/quic-s/quics/pkg/types"
)

func FullScanMain() {

	err := QPClient.RecvTransactionHandleFunc(qstypes.FULLSCAN, func(conn *qp.Connection, stream *qp.Stream, transactionName string, transactionID []byte) error {
		log.Print("quics-client : [FULL SCAN] transaction start")
		askAllMetaReq, err := qclient.AskAllMetaRecvHandler(stream)
		if err != nil {
			return err
		}
		if askAllMetaReq.UUID == "" {
			return fmt.Errorf("UUID is empty")
		}
		UUID := badger.GetUUID()
		if askAllMetaReq.UUID != UUID {
			return fmt.Errorf("UUID is not same")
		}

		// make list to compare
		// If os.Stat of certain path is existed then set IsExisted true
		rawlist, err := badger.GetAllSyncMetadataAmongRoot()
		if err != nil {
			return err
		}
		comparelist := []*types.ComparingSyncMetadata{}
		for _, item := range rawlist {
			convertedItem := types.ComparingSyncMetadata{
				Path:      filepath.Join(item.BeforePath, item.AfterPath),
				IsExisted: false,
				Sync:      *item,
			}
			comparelist = append(comparelist, &convertedItem)
		}

		// Compare to file.Info and go ps thread to update syncMetadata as file.Info says
		resultList := []qstypes.SyncMetadata{}

		rootDirList := badger.GetRootDirList()
		for _, rootDir := range rootDirList {
			rootpath := rootDir.Path

			err = filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					info, err := os.Stat(path)
					if err != nil {
						return err
					}

					// OS : O, SyncMetadata : X
					if item := ChangeTrueInComparelistIfExisted(comparelist, path); item == nil {
						go PleaseSync(path)
						BeforePath, AfterPath := badger.SplitBeforeAfterRoot(path)

						// TODO: Consider making hash of file is needed
						resultList = append(resultList, qstypes.SyncMetadata{
							BeforePath:          BeforePath,
							AfterPath:           AfterPath,
							LastUpdateHash:      utils.MakeHash(AfterPath, info),
							LastUpdateTimestamp: 1,
							LastSyncHash:        "",
							LastSyncTimestamp:   0,
						})
						return nil

					} else {
						// OS : O, SyncMetadata : O
						hashtocompare := utils.MakeHash(item.Sync.AfterPath, info)
						if item.Sync.LastUpdateHash != hashtocompare {
							go PleaseSync(path)
							// OS : O, SyncMetadata : O, hash is not same -> MS
							convertedItem := qstypes.SyncMetadata{
								BeforePath:          item.Sync.BeforePath,
								AfterPath:           item.Sync.AfterPath,
								LastUpdateHash:      hashtocompare,
								LastUpdateTimestamp: item.Sync.LastUpdateTimestamp + 1,
								LastSyncTimestamp:   item.Sync.LastSyncTimestamp,
								LastSyncHash:        item.Sync.LastSyncHash,
							}
							resultList = append(resultList, convertedItem)
						} else {
							// OS : O, SyncMetadata : O, hash is same -> MS
							convertedItem := qstypes.SyncMetadata{
								BeforePath:          item.Sync.BeforePath,
								AfterPath:           item.Sync.AfterPath,
								LastUpdateHash:      item.Sync.LastUpdateHash,
								LastUpdateTimestamp: item.Sync.LastUpdateTimestamp,
								LastSyncTimestamp:   item.Sync.LastSyncTimestamp,
								LastSyncHash:        item.Sync.LastSyncHash,
							}
							resultList = append(resultList, convertedItem)

						}
						return nil

					}
				}

				return nil
			})
			if err != nil {
				continue
			}
		}
		// OS : X, SyncMetadata : O
		// When Case in Remove
		for _, item := range GetComparelistIfNotExisted(comparelist) {
			go PleaseSync(item.Path)
			convertedItem := qstypes.SyncMetadata{
				BeforePath:          item.Sync.BeforePath,
				AfterPath:           item.Sync.AfterPath,
				LastUpdateHash:      item.Sync.LastUpdateHash,
				LastUpdateTimestamp: item.Sync.LastUpdateTimestamp,
				LastSyncTimestamp:   item.Sync.LastSyncTimestamp,
				LastSyncHash:        item.Sync.LastSyncHash,
			}
			resultList = append(resultList, convertedItem)
		}

		// return askAllMetaHandler
		err = qclient.AskAllMetaHandler(stream, askAllMetaReq.UUID, resultList)
		if err != nil {
			return err
		}
		log.Println("quics-client : [FULL SCAN] transaction success")
		return nil
	})
	if err != nil {
		log.Println("quics-client : [FULLSCAN]: ", err)
	}
}

// ChangeTrueInComparelistIfExisted changes IsExisted true if path is existed in comparelist
func ChangeTrueInComparelistIfExisted(comparelist []*types.ComparingSyncMetadata, path string) *types.ComparingSyncMetadata {
	for _, item := range comparelist {
		if item.Path == path {
			item.IsExisted = true
			return item
		}
	}
	return nil
}

// GetComparelistIfNotExisted returns list of SyncMetadata which is not existed in OS
func GetComparelistIfNotExisted(comparelist []*types.ComparingSyncMetadata) []*types.ComparingSyncMetadata {
	resultList := []*types.ComparingSyncMetadata{}
	for _, item := range comparelist {
		if !item.IsExisted {
			resultList = append(resultList, item)
		}
	}
	return resultList
}
