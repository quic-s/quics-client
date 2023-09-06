package sync

import (
	"github.com/quic-s/quics-client/pkg/badger"
	"github.com/quic-s/quics-client/pkg/types"
)

// @URL /api/v1/status/root/
// ex) ShowStatus("/home/rootDir/text.txt)"
func ShowStatus(filepath string) types.SyncMetadata {
	value, err := badger.View(filepath)
	if err != nil {
		return types.SyncMetadata{}
	}
	syncMetadata := types.SyncMetadata{}
	syncMetadata.Decode(value)
	return syncMetadata
}
