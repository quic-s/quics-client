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
	// ch := make(chan int)
	// connection.InitWatcher()
	// rootdirlist := utils.GetRootDirs()
	// for _, value := range rootdirlist {
	// 	if value != "" {
	// 		connection.DirWatchAdd(value)

	// 	}

	// }
	// connection.DirWatchStart()
	// <-ch

}
