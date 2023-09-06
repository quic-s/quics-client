package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Reboot() {
	log.Println("quics-client :\n\trebooting ...")

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
