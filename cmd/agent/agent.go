package main

import (
	loc_metric "github.com/imaksimSuhinin/yandexSprintOne/internal/metrics"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/os"
	"net/http"
	"time"
)

const (
	delayRefresh      time.Duration = 2
	delayUpload       time.Duration = 10
	httpClientTimeOut time.Duration = 30
)

func main() {

	var metrics loc_metric.Metrics
	var client = startClient()

	upload := time.NewTicker(delayUpload * time.Second)
	refresh := time.NewTicker(delayRefresh * time.Second)

	for {
		select {
		case <-upload.C:
			metrics.PostMetrics(client)
		case <-refresh.C:
			metrics.UpdateMetrics()
		case <-os.SigChanel:
			os.UpdateOsSignal()
		}
	}
}

func startClient() *http.Client {
	client := &http.Client{
		Timeout: httpClientTimeOut * time.Second,
	}

	return client
}
