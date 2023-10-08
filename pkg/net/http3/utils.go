package http3

import (
	"fmt"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/db/badger"
)

func PrintListOfTwoOptions() {
	conflictFileList := badger.GetConflictFileList()
	for i, conflictFile := range conflictFileList {
		fmt.Printf("%d. %s\n", i+1, filepath.Join(conflictFile.BeforePath, conflictFile.AfterPath))
		fmt.Printf("\tServer ModDate: %s\n", conflictFile.Conflict.ServerModDate)
		fmt.Printf("\tLocal ModDate: %s\n", conflictFile.Conflict.LocalModDate)
	}
}
