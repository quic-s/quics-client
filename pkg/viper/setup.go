package viper

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var QicEnvPath string

func InitViper() {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	tem := filepath.Join(tempDir, ".quics")
	QicEnvPath = filepath.Join(tem, "qic.env")

	_, err = os.Stat(QicEnvPath)
	if os.IsNotExist(err) {

		viper.SetConfigFile(".env")
		viper.SetConfigType("env")

		err = viper.ReadInConfig()
		if err != nil {
			log.Println("quics-client : Error while initial reading config file : ", err)
		}
		err = viper.WriteConfigAs(QicEnvPath)
		if err != nil {
			log.Fatalf("quics-client : Error while writing config file   %s", err)
		}
	} else {
		viper.SetConfigFile(QicEnvPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("quics-client : Error while reading config file : ", err)

	}
}
