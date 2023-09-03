package viper

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

var (
	tem        = ""
	QicEnvPath = path.Join(tem + ".qic.env")
)

func init() {
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	tem = filepath.Join(tempDir, ".quics")

	_, err = os.Stat(QicEnvPath)
	if os.IsNotExist(err) {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		err = viper.WriteConfigAs(QicEnvPath)
		if err != nil {
			log.Fatalf("Error while writing config file :  %s", err)
		}
	} else {
		viper.SetConfigFile(QicEnvPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file : ", err)

	}
}
