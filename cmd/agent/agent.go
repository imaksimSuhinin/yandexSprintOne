package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/xlab/closer"
	"log"
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
	strtClient()

}

func strtClient() {
	client := resty.New()
	var m loc_metric.Metrics

	refresh := time.NewTicker(2 * time.Second)
	upload := time.NewTicker(10 * time.Second)

	for {
		<-refresh.C
		var z = m.UpdateMetrics()

		<-upload.C
		z.PostMetrics(client)

	}
}

func Exit() {
	log.Println("Exit...")
	os.Exit(0)

}
