package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {

	//os.Exit(Run())
	// utils.CreateDirIfNotExisted()
	//http.RestServerStart()

	log.Println(viper.AllSettings())
}
