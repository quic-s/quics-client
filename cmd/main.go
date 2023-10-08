package main

import (
	"os"

	"github.com/quic-s/quics-client/pkg/utils"
	"github.com/quic-s/quics-client/pkg/viper"
)

func main() {

	utils.CreateDirIfNotExisted()
	viper.InitViper()

	os.Exit(Run())

}
