package main

import (
	loc_metric "github.com/imaksimSuhinin/yandexSprintOne/internal/metrics"
	os "github.com/imaksimSuhinin/yandexSprintOne/internal/os"
	"net/http"
	"time"
)

const (
	delayRefresh      time.Duration = 2
	delayUpload       time.Duration = 10
	httpClientTimeOut time.Duration = 10
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
		Timeout: httpClientTimeOut * time.Second,
	}

	return client
}

func getUpload(m *loc_metric.Metrics, client *http.Client) {
	upload := time.NewTicker(delayUpload * time.Second)
	for {
		select {
		case <-upload.C:
			m.PostMetrics(client)
		}
	}
}

func getRefresh(m *loc_metric.Metrics) {
	refresh := time.NewTicker(delayRefresh * time.Second)

	for {
		select {
		case <-refresh.C:
			m.UpdateMetrics()
		}
	}
}
