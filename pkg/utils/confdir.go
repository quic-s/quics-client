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
	return filepath.Join(tempDir, "quics")
}

// read .qis.env file if it is existed
func ReadEnvFile() []map[string]string {
	envPath := filepath.Join(GetQuicsDirPath(), ".qic.env")
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

	//Key-value 형태의 리스트로 만듦
	dataStr := string(data)

	// 문자열을 줄 단위로 분리
	lines := strings.Split(dataStr, "\n")

	// 줄마다 키와 값으로 분리하고 리스트에 추가
	var kvList []map[string]string
	for _, line := range lines {
		// 줄이 비어있으면 건너뛰기
		if line == "" {
			continue
		}
		// 줄을 공백으로 분리
		parts := strings.Split(line, "=")
		// 첫 번째 부분을 키로, 나머지 부분을 값으로 합쳐서 맵에 저장
		key := parts[0]
		value := strings.Join(parts[1:], " ")
		kvMap := map[string]string{key: value}
		// 맵을 리스트에 추가
		kvList = append(kvList, kvMap)
	}
	// 파일 내용 출력
	log.Println(kvList)
	return kvList
}

func GetRootDir() []map[string]string {
	rawList := ReadEnvFile()
	kvList := []map[string]string{}
	for _, kvMap := range rawList {
		for key, value := range kvMap {
			if len(key) > 5 && key[:5] == "ROOT." {
				kvList = append(kvList, map[string]string{key[5:]: value})
			}
		}
	}
	return kvList
}

func IsRootDir(rootpath string) bool {
	rootDirList := GetRootDir()
	for _, kvMap := range rootDirList {
		for key, _ := range kvMap {
			if key == rootpath {
				return true
			}
		}
	}
	return false
}
