package viper

import (
	"log"

	"github.com/spf13/viper"
)

// use viper package to read .env file
// return the value of the key
func GetViperEnvVariables(key string) string {

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("quics-client :Error while reading config file : ", err)
		return ""

	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func WriteViperEnvVariables(key string, value string) {

	viper.Set(key, value)
	err := viper.WriteConfigAs(QicEnvPath)
	if err != nil {
		log.Fatalf("Error while writing config file :  %s", err)
	}
}

func DeleteViperVariablesByKey(key string) {

	// Future Plan : Update
	WriteViperEnvVariables(key, "")

}
