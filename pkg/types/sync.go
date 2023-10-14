package types

import (
	"bytes"
	"encoding/gob"
)

// key = path
type SyncMetadata struct { // Per file
	BeforePath          string
	AfterPath           string
	LastUpdateTimestamp uint64 // Local File changed time
	LastUpdateHash      string
	LastSyncTimestamp   uint64 // Sync Success Time
	LastSyncHash        string
}

func (syncMetadata *SyncMetadata) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(syncMetadata); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
func (syncMetadata *SyncMetadata) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(syncMetadata); err != nil {
		panic(err)
	}

}
