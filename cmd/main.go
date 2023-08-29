package main

import (
	"github.com/quic-s/quics-client/pkg/quic"
)

func main() {
	//os.Exit(Run())
	// utils.CreateDirIfNotExisted()
	// http.RestServerStart()
	quic.ClientMessage("get", []byte("test"))

}
