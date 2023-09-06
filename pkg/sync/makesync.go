package sync

import (
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/types"
)

func CanOverWrite(lastUpdate uint64, lastSync uint64, lastestSync uint64) bool {
	if lastUpdate != lastSync {
		return false
	}
	if lastestSync <= lastSync {
		return false
	}
	return true
}

func UpdateSyncMetadata(syncMetadata types.SyncMetadata) error {
	error := badger.Update(syncMetadata.Path, syncMetadata.Encode())
	if error != nil {
		return error
	}
	return nil
}
