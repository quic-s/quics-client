package sync

import (
	"reflect"
	"strconv"

	"github.com/quic-s/quics-client/pkg/db/badger"
)

// @URL /api/v1/status/root/
// ex) ShowStatus("/home/rootDir/text.txt)"
func ShowStatus(filepath string) (string, error) {
	value := badger.GetSyncMetadata(filepath)
	if reflect.ValueOf(value).IsZero() {
		return "", nil
	}
	result := "\n\n=== Status ===\n"
	result += "path : " + filepath + "\n"
	result += "LastUpdateTimestamp : " + strconv.Itoa(int(value.LastUpdateTimestamp)) + "\n"
	result += "LastUpdateHash : " + value.LastUpdateHash + "\n"
	result += "LastSyncTimestamp : " + strconv.Itoa(int(value.LastSyncTimestamp)) + "\n"
	result += "LastSyncHash : " + value.LastSyncHash + "\n"
	result += "===============\n"

	return result, nil

}
