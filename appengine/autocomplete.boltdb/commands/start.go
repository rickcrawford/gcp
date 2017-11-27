package commands

import (
	"log"
	"os"
	"syscall"
)

func start(sig <-chan os.Signal) bool {
	var exit bool
	switch <-sig {
	case syscall.SIGINT, syscall.SIGTERM:
		exit = true

	// case syscall.SIGHUP:
	default:
		log.Println("reload")
	}
	return exit
}
