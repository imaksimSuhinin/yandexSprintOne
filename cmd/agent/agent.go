package main

import (
	"github.com/imaksimSuhinin/yandexSprintOne/internal/config"
	loc_metric "github.com/imaksimSuhinin/yandexSprintOne/internal/metrics"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/os"
	"net/http"
	"time"
)

//const (
//	delayRefresh      time.Duration = 2
//	delayUpload       time.Duration = 10
//	httpClientTimeOut time.Duration = 20
//)

func main() {
	conf := config.New()

	var metrics loc_metric.Metrics
	var client = startClient(*conf)

	upload := time.NewTicker(conf.AgentConfig.ReportInterval * time.Second)
	refresh := time.NewTicker(conf.AgentConfig.PollInterval * time.Second)

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

func startClient(conf config.Config) *http.Client {
	client := &http.Client{
		Timeout: conf.AgentConfig.HTTPClientTimeOut * time.Second,
	}

	return client
}
