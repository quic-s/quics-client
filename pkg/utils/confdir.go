package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ViperConfig struct {
	Key   string
	Value string
}

// CreateDirIfNotExisted creates the quics folder if it does not exist
func CreateDirIfNotExisted() {
	quicsDir := GetQuicsDirPath()
	log.Println(quicsDir)

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

// GetQuicsDirPath returns the path of the quics folder
func GetQuicsDirPath() string {

	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(tempDir, ".quics")
}

// get new emptyfile in quicsDir,  if not exist then create new emptyfile
func GetEmptyFilePath() string {
	quicsDir := GetQuicsDirPath()
	emptyFilePath := filepath.Join(quicsDir, "emptyfile")
	_, err := os.Stat(emptyFilePath)
	if os.IsNotExist(err) {
		emptyFile, err := os.Create(emptyFilePath)
		if err != nil {
			log.Fatal(err)
		}
		emptyFile.Close()
	}
	return emptyFilePath
}

func ReadEnvFile() map[string]string {
	envPath := filepath.Join(GetQuicsDirPath(), "qic.env")
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

	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")

	// 줄마다 키와 값으로 분리하고 리스트에 추가
	kvMap := make(map[string]string)
	for _, line := range lines {
		// 줄이 비어있으면 건너뛰기
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		key := parts[0]
		value := strings.Join(parts[1:], " ")
		kvMap[key] = value
	}

	//log.Println(kvMap)
	return kvMap
}
