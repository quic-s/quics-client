package sync

import (
	"os"
	"strconv"

	"github.com/quic-s/quics-client/pkg/db/badger"
)

// @URL /api/v1/sync/status
// ex) ShowStatus("/home/rootDir/text.txt)"
func ShowStatus(filepath string) (string, error) {
	pathInfo, err := os.Stat(filepath)
	if err != nil {
		return "", err
	}
	if pathInfo.IsDir() {
		return "", nil
	}

	if !badger.IsSyncMetadataExisted(filepath) {
		return "", nil
	}

	value := badger.GetSyncMetadata(filepath)
	result := "\n\n=== Status ===\n"
	result += "path : " + filepath + "\n"
	result += "LastUpdateTimestamp : " + strconv.Itoa(int(value.LastUpdateTimestamp)) + "\n"
	result += "LastUpdateHash : " + value.LastUpdateHash + "\n"
	result += "---------------\n"
	result += "LastSyncTimestamp : " + strconv.Itoa(int(value.LastSyncTimestamp)) + "\n"
	result += "LastSyncHash : " + value.LastSyncHash + "\n"
	result += "===============\n"

	return result, nil

}
