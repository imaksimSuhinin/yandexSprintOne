package main

import (
	"net/http"
	"time"
	loc_metric "yandexSprintOne/internal/metrics"
	os "yandexSprintOne/internal/os"
)

func main() {

	go startClient()
	os.UpdateOsSignal()
}

func startClient() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	var m loc_metric.Metrics
	refresh := time.NewTicker(2 * time.Second)
	upload := time.NewTicker(10 * time.Second)

	for {
		<-refresh.C
		m.UpdateMetrics()

		<-upload.C
		m.PostMetrics(client)
	}
}
