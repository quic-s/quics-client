package app

import (
	"log"
	"os"
	"syscall"
)

// TODO
func Reboot() error {
	log.Println("quics-client :\n\nrebooting ...")

	self, err := os.Executable()

	if err != nil {
		return err
	}
	args := os.Args // 인자를 그대로 전달합니다.
	env := os.Environ()

	err = syscall.Exec(self, args, env) // 현재 프로세스를 새로운 프로세스로 대체합니다.
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
