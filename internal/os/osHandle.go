package os

import (
	"fmt"
	"os"
	"os/signal"
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

	signal.Notify(SigChanel)
	exitChanel := make(chan int)
	s := <-SigChanel
	HandleOsSignal(s)
	exitCode := <-exitChanel
	os.Exit(exitCode)
}
