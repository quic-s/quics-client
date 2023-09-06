package types

import (
	"bytes"
	"encoding/gob"
)

type MessageData interface {
	Encode() []byte
	Decode([]byte)
}

type RegisterClientRequest struct {
	Uuid           string
	ClientPassword string
}

type RegisterClientResponse struct {
	Uuid string
}

type RegisterRootDirRequest struct {
	Uuid            string
	RootDirPassword string
	// e.g., /home/ubuntu/rootDir/*
	BeforePath string // /home/ubuntu
	AfterPath  string // /rootDir/*
}

type NotClientAnymoreRequest struct {
	Uuid           string
	ClientPassword string
}

type NotRootDirAnymorRequest struct {
	Uuid            string
	RootPath        string
	RootDirPassword string
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

func (notClientAnymoreRequest *NotClientAnymoreRequest) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(notClientAnymoreRequest); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (notClientAnymoreRequest *NotClientAnymoreRequest) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(notClientAnymoreRequest); err != nil {
		panic(err)
	}

}

func (notRootDirAnymorRequest *NotRootDirAnymorRequest) Encode() []byte {
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(notRootDirAnymorRequest); err != nil {
		panic(err)
	}

	return buffer.Bytes()
}

func (notRootDirAnymorRequest *NotRootDirAnymorRequest) Decode(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(notRootDirAnymorRequest); err != nil {
		panic(err)
	}

}
