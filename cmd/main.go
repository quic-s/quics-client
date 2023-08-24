package main

import (
	"log"

	"github.com/quic-s/quics-client/pkg/utils"
)

func main() {
	// os.Exit(Run())
	// utils.CreateDirIfNotExisted()
	// http.RestServerStart()
	utils.WriteViperEnvVariables("hello2", "hi")
	log.Println(utils.GetViperEnvVariables("HELLO3"))

}
