package os

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var SigChanel = make(chan os.Signal, 1)

func HandleOsSignal(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else if signal == syscall.SIGQUIT {
		fmt.Println("Got Quit signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else {
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
