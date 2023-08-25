package utils

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// use viper package to read .env file
// return the value of the key
func GetViperEnvVariables(key string) string {

	envPath := filepath.Join(GetDirPath(), ".qic.env")
	_, err := os.Stat(envPath)
	if err != nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
	} else {

		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file : ", err)
		return ""

	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value

}

func WriteViperEnvVariables(key string, value string) {
	envPath := filepath.Join(GetDirPath(), ".qic.env") // force to write
	_, err := os.Stat(envPath)
	if os.IsNotExist(err) {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		err = viper.WriteConfigAs(envPath)
		if err != nil {
			log.Fatalf("Error while writing config file :  %s", err)
		}
	} else {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading config file : ", err)

	}
	viper.Set(key, value)
	err = viper.WriteConfigAs(envPath)
	if err != nil {
		log.Fatalf("Error while writing config file :  %s", err)
	}
}
