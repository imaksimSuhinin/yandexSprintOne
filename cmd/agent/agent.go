package main

import (
	"github.com/go-resty/resty/v2"
	"yandexSprintOne/internal/metrics"
)

func main() {
	client := resty.New()
	var metrics runtime_loc.Metrics
	metrics.UpdateMetrics(2)
	metrics.PostMetrics(client, 4)

}
