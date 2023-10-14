package sync

import (
	"fmt"
)

const (
	SERVER = "SERVER"
	LOCAL  = "LOCAL"
)

func PrintTwoOptions(path string, serverModDate string, localModDate string) {
	fmt.Println("---- FILE CONFLICTED ----")
	fmt.Println(" path >> ", path)
	fmt.Printf("Server ModDate: %s\n", serverModDate)
	fmt.Printf("Local ModDate: %s\n", localModDate)
	fmt.Println("-------------------------\n")
	fmt.Println(" Choose one between two options, 1 or 2 ")
	fmt.Println(" 1. File At Server")
	fmt.Println(" 2. File At Local")
}

// func ChooseOne(path string, Side string) error {

// 	UUID := badger.GetUUID()
// 	prevSyncMetadata := badger.GetSyncMetadata(path)

// 	err := Conn.OpenTransaction("CONFLICT", func(stream *stream.Stream, transactionName string, transactionID []byte) error {
// 		if Side == SERVER {
// 			newTimestamp := prevSyncMetadata.Conflict.ServerTimestamp + 1
// 			newHash := prevSyncMetadata.Conflict.ServerHash
// 			_, err := qclient.SendPleaseServerFile(stream, path, UUID, prevSyncMetadata.AfterPath, prevSyncMetadata.Conflict.ServerTimestamp, newTimestamp, newHash)
// 			if err != nil {
// 				return err
// 			}

// 			//Update Sync Metadata
// 			syncMetadata := types.SyncMetadata{
// 				BeforePath:          prevSyncMetadata.BeforePath,
// 				AfterPath:           prevSyncMetadata.AfterPath,
// 				LastUpdateTimestamp: newTimestamp,
// 				LastUpdateHash:      newHash,
// 				LastSyncTimestamp:   newTimestamp,
// 				LastSyncHash:        newHash,
// 			}
// 			badger.Update(path, syncMetadata.Encode())

// 		} else if Side == LOCAL {
// 			newTimestamp := prevSyncMetadata.Conflict.LocalTimestamp + 1
// 			newHash := prevSyncMetadata.Conflict.LocalHash
// 			_, err := qclient.SendPleaseLocalFile(stream, path, UUID, prevSyncMetadata.AfterPath, prevSyncMetadata.Conflict.LocalTimestamp, newTimestamp, newHash)
// 			if err != nil {
// 				return err
// 			}

// 			//Update Sync Metadata
// 			syncMetadata := types.SyncMetadata{
// 				BeforePath:          prevSyncMetadata.BeforePath,
// 				AfterPath:           prevSyncMetadata.AfterPath,
// 				LastUpdateTimestamp: newTimestamp,
// 				LastUpdateHash:      newHash,
// 				LastSyncTimestamp:   newTimestamp,
// 				LastSyncHash:        newHash,
// 			}
// 			badger.Update(path, syncMetadata.Encode())

// 		} else {
// 			fmt.Errorf("chosenOne is not valid")
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return fmt.Errorf("[CONFLICT] ", err)
// 	}
// 	return nil
// }
