package types

import (
	"bytes"
	"encoding/gob"
)

type RootDirList []RootDir

type RootDir struct {
	NickName     string `json:"nickName"`
	Path         string `json:"path"`
	BeforePath   string `json:"beforePath"`
	AfterPath    string `json:"afterPath"`
	IsRegistered bool   `json:"isRegistered"`
}

func (r *RootDirList) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(r); err != nil {

		panic(err)
	}
}

func (r *RootDirList) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(r); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (r *RootDir) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(r); err != nil {
		panic(err)
	}
}

func (r *RootDir) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(r); err != nil {

		panic(err)
	}

	return buffer.Bytes()
}
