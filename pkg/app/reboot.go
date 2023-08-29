package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func Reboot(cmd *cobra.Command, args []string) {
	log.Println("\n\trebooting ...")
	// 현재 실행 중인 프로세스의 경로와 인자를 얻습니다.
	str, err := os.Executable()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	newProcess := exec.Command(str)
	newProcess.Stdout = os.Stdout
	newProcess.Stderr = os.Stderr
	newProcess.Stdin = os.Stdin
	newProcess.Env = os.Environ()

	err = newProcess.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}
