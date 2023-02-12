package main

import (
	"github.com/imaksimSuhinin/yandexSprintOne/internal/config"
	loc_metric "github.com/imaksimSuhinin/yandexSprintOne/internal/metrics"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/os"
	"net/http"
	"time"
)

const (
	//delayRefresh      time.Duration = 2
	//delayUpload       time.Duration = 10
	httpClientTimeOut time.Duration = 20
)

func main() {

	var metrics loc_metric.Metrics
	var client = startClient()

	upload := time.NewTicker(config.AppConfig.ReportInterval * time.Second)
	refresh := time.NewTicker(config.AppConfig.PollInterval * time.Second)

	for {
		select {
		case <-upload.C:
			metrics.PostMetricsJSON(client)
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
