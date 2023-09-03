package sharing

import (
	"bytes"
	"encoding/gob"
)

type FileDownloadRequest struct {
	Uuid       string
	BeforePath string // /home/ubuntu
	AfterPath  string // /rootDir/file
	MaxCnt     uint32
}

type FileDownloadResponse struct {
	Link  string
	Count uint32
}

//Encoding, Decoding

func (fileDownloadRequest *FileDownloadRequest) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileDownloadRequest); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}
func (fileDownloadRequest *FileDownloadRequest) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(fileDownloadRequest); err != nil {
		panic(err)
	}

}

func (fileDownloadResponse *FileDownloadResponse) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fileDownloadResponse); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (fileDownloadResponse *FileDownloadResponse) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(fileDownloadResponse); err != nil {
		panic(err)
	}

}
