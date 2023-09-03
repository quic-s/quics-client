package connection

import (
	"bytes"
	"encoding/gob"
)

type RegisterClientRequest struct {
	Ip string
}

type RegisterClientResponse struct {
	Uuid string
}

type RegisterRootDirRequest struct {
	Uuid     string
	Password string
	// e.g., /home/ubuntu/rootDir/*
	BeforePath string // /home/ubuntu
	AfterPath  string // /rootDir/*
}

func (registerRootDirRequest *RegisterRootDirRequest) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerRootDirRequest); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (registerRootDirRequest *RegisterRootDirRequest) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(registerRootDirRequest); err != nil {
		panic(err)
	}

}

func (registerClientRequest *RegisterClientRequest) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerClientRequest); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (registerClientRequest *RegisterClientRequest) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(registerClientRequest); err != nil {
		panic(err)
	}

}

func (registerClientResponse *RegisterClientResponse) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(registerClientResponse); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (registerClientResponse *RegisterClientResponse) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(registerClientResponse); err != nil {
		panic(err)
	}

}
