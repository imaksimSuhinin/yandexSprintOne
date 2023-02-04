package os

import (
	"fmt"
	"os"
	"syscall"
)

var SigChanel = make(chan os.Signal, 1)

func HandleOsSignal(signal os.Signal) {
	switch signal {
	case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
		fmt.Printf("Got %s signal\n", signal)
		fmt.Println("Closing")
		os.Exit(0)
	default:
		fmt.Println("Ignoring signal: ", signal)
	}
}

func UpdateOsSignal() {
	s := <-SigChanel
	HandleOsSignal(s)
	os.Exit(0)
}
