package types

import (
	"bytes"
	"encoding/gob"
)

type PleaseFile struct {
	Uuid          string
	BeforePath    string
	AfterPath     string
	SyncTimestamp uint64
}

func (pleaseFile *PleaseFile) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(pleaseFile); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
func (pleaseFile *PleaseFile) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(pleaseFile); err != nil {
		panic(err)
	}

}
