package utils

import (
	"log"
	"os"
	"path/filepath"
)

// CreateDirIfNotExisted creates the quics folder if it does not exist
func CreateDirIfNotExisted() {
	quicsDir := GetDirPath()

	_, err := os.Stat(quicsDir)

	if os.IsNotExist(err) {
		err = os.Mkdir(quicsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created quics folder:", quicsDir)
	} else {
		log.Println("Using existing quics folder:", quicsDir)
	}

}

// GetDirPath returns the path of the quics folder
func GetDirPath() string {

	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(tempDir, "quics")
}
