package app

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func Reboot() {
	log.Println("quics-client :\n\nrebooting ...")

	// str, err := os.Executable()

	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// newProcess := exec.Command(str)
	// newProcess.Stdout = os.Stdout
	// newProcess.Stderr = os.Stderr
	// newProcess.Stdin = os.Stdin
	// newProcess.Env = os.Environ()

	// err = newProcess.Start()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// os.Exit(0)

	self, err := os.Executable()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	args := os.Args // 인자를 그대로 전달합니다.
	env := os.Environ()

	err = syscall.Exec(self, args, env) // 현재 프로세스를 새로운 프로세스로 대체합니다.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
