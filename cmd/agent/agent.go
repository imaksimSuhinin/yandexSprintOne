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
	upload := time.NewTicker(delayUpload * time.Second)
	refresh := time.NewTicker(delayRefresh * time.Second)
	for {
		select {
		case <-upload.C:
			metrics.PostMetrics(startClient())
		case <-refresh.C:
			metrics.UpdateMetrics()
		}
	}
	os.UpdateOsSignal()
}

func startClient() *http.Client {
	client := &http.Client{
		Timeout: httpClientTimeOut * time.Second,
	}

	return client
}
