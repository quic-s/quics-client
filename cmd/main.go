package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/quic-s/quics-client/pkg/http"
)

func main() {
	// os.Exit(Run())
	// temp 폴더의 경로를 얻습니다.
	tempDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// quics라는 폴더의 경로를 만듭니다.
	quicsDir := filepath.Join(tempDir, "quics")

	// quics 폴더가 존재하는지 확인합니다.
	_, err = os.Stat(quicsDir)

	// quics 폴더가 존재하지 않으면
	if os.IsNotExist(err) {
		// quics 폴더를 생성합니다.
		err = os.Mkdir(quicsDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created quics folder:", quicsDir)
	} else {
		// quics 폴더가 존재하면
		log.Println("Using existing quics folder:", quicsDir)
	}

	http.RestServerStart()

}
