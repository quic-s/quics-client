package types

import (
	"bytes"
	"encoding/gob"
)

type SyncMetadata struct { // Per file
	Path                string // key, Local Absolute Path
	LastUpdateTimestamp uint64 // Local File changed time
	LastUpdateHash      string
	LastSyncTimestamp   uint64 // Sync Success Time
	LastSyncHash        string
}

type PleaseSync struct {
	Uuid  string
	Event string
	// e.g., /home/ubuntu/rootDir/file
	BeforePath          string // /home/ubuntu
	AfterPath           string // /rootDir/file
	LastUpdateTimestamp uint64
	LastUpdateHash      string
}

type MustSync struct {
	LatestHash          string // depends on server
	LatestSyncTimestamp uint64 // depends on server
	BeforePath          string
	AfterPath           string
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

func (pleaseSync *PleaseSync) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseSync); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
func (pleaseSync *PleaseSync) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(pleaseSync); err != nil {
		panic(err)
	}

}

func (mustSync *MustSync) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(mustSync); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (mustSync *MustSync) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(mustSync); err != nil {
		panic(err)
	}

}
