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
	return filepath.Join(tempDir, ".quics")
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

	dataStr := string(data)
	lines := strings.Split(dataStr, "\n")

	// 줄마다 키와 값으로 분리하고 리스트에 추가
	var kvList []map[string]string
	for _, line := range lines {
		// 줄이 비어있으면 건너뛰기
		if line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		key := parts[0]
		value := strings.Join(parts[1:], " ")
		kvMap := map[string]string{key: value}
		kvList = append(kvList, kvMap)
	}

	log.Println(kvList)
	return kvList
}

func GetRootDirs() []map[string]string {
	// ex) GetRootDirs() -> [{"ROOT.b": "/home/user/a/b"}, {"ROOT.d": "/home/user/c/d"}]
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

func GetRootDir(key string) string {
	// ex) GetRootDir("b") -> "/home/user/a/b"
	key = "ROOT." + key
	rootDirList := GetRootDirs()
	for _, kvMap := range rootDirList {
		for k, v := range kvMap {
			if k == key {
				return v
			}
		}
	}
	return ""
}

func IsRootDir(rootpath string) bool {
	// ex) IsRootDir("/home/user/a/b") -> true
	rootDirList := GetRootDirs()
	for _, kvMap := range rootDirList {
		for _, value := range kvMap {
			if value == rootpath {
				return true
			}
		}
	}
	return false
}

func IsDuplicateKey(key string) bool {
	// ex) IsDuplicateKey("b") -> true
	rootDirList := GetRootDirs()
	for _, kvMap := range rootDirList {
		for k, _ := range kvMap {
			if k == "ROOT."+key {
				return true
			}
		}
	}
	return false
}

func IsDuplicateValue(value string) bool {
	// ex) IsDuplicateValue("/home/user/a/b") -> true
	rootDirList := GetRootDirs()
	for _, kvMap := range rootDirList {
		for _, v := range kvMap {
			if v == value {
				return true
			}
		}
	}
	return false
}
