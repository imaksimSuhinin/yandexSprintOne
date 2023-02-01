package main

import (
	"github.com/xlab/closer"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	loc_metric "yandexSprintOne/internal/metrics"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	closer.DebugSignalSet = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}
}

func main() {

	closer.Bind(Exit)
	startClient()

}

func Exit() {
	log.Println("Exit...")
	os.Exit(0)

}

func startClient() {
	client := &http.Client{}
	var m loc_metric.Metrics
	refresh := time.NewTicker(2 * time.Second)
	upload := time.NewTicker(2 * time.Second)

	for {
		<-refresh.C
		m.UpdateMetrics()

		<-upload.C
		m.PostMetrics(client)

	}
}
