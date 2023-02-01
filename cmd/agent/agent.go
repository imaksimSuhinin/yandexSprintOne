package main

import (
	"net/http"
	"time"
	loc_metric "yandexSprintOne/internal/metrics"
	os "yandexSprintOne/internal/os"
)

func main() {

	var metrics loc_metric.Metrics
	go getRefresh(&metrics)
	go startClient()
	go getUpload(&metrics, startClient())
	os.UpdateOsSignal()
}

func startClient() *http.Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return client
}

func getUpload(m *loc_metric.Metrics, client *http.Client) {
	upload := time.NewTicker(10 * time.Second)
	for {
		<-upload.C
		m.PostMetrics(client)
	}
}

func getRefresh(m *loc_metric.Metrics) {
	refresh := time.NewTicker(2 * time.Second)

	for {
		<-refresh.C
		m.UpdateMetrics()
	}
}
