package main

import (
	"github.com/go-resty/resty/v2"
	"time"
	loc_metric "yandexSprintOne/internal/metrics"
)

func main() {

	client := resty.New()
	var m loc_metric.Metrics

	refresh := time.NewTicker(1 * time.Second)
	upload := time.NewTicker(2 * time.Second)

	for {
		<-refresh.C
		var z = m.UpdateMetrics()

		<-upload.C
		z.PostMetrics(client)

	}
}
