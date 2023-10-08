package types

import (
	"bytes"
	"encoding/gob"
)

type ConflictFileList []SyncMetadata

// key = path
type SyncMetadata struct { // Per file
	BeforePath          string
	AfterPath           string
	LastUpdateTimestamp uint64 // Local File changed time
	LastUpdateHash      string
	LastSyncTimestamp   uint64 // Sync Success Time
	LastSyncHash        string
	Conflict            ConflictMetadata
}

type ConflictMetadata struct {
	BeforePath string
	AfterPath  string

	ServerDevice    string
	ServerTimestamp uint64
	ServerHash      string
	ServerModDate   string

	LocalDevice    string
	LocalTimestamp uint64
	LocalHash      string
	LocalModDate   string
}

func (conflictFileList *ConflictFileList) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(conflictFileList); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
func (conflictFileList *ConflictFileList) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(conflictFileList); err != nil {
		panic(err)
	}

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
