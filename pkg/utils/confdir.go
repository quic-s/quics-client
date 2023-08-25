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

// read .qis.env file if it is existed
func ReadEnvFile() {
	envPath := filepath.Join(GetDirPath(), ".qic.env")
	file, err := os.Open(envPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// 파일 전체 읽기
	data, err := os.ReadFile(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	// 파일 내용 출력
	log.Println(string(data))
}
