package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func main() {
	client := resty.New()
	var m Metrics

	refresh := time.NewTicker(2 * time.Second)
	upload := time.NewTicker(10 * time.Second)

	_ = <-refresh.C
	m.UpdateMetrics()

	_ = <-upload.C
	m.PostMetrics(client)

}

type gauge float64
type count int64

type Metrics struct {
	Alloc,
	TotalAlloc,
	LiveObjects,
	BuckHashSys,
	Frees,
	GCCPUFraction,
	GCSys,
	HeapAlloc,
	HeapIdle,
	HeapInuse,
	HeapObjects,
	HeapReleased,
	HeapSys,
	LastGC,
	Lookups,
	MCacheInuse,
	MCacheSys,
	MSpanInuse,
	MSpanSys,
	Mallocs,
	NextGC,
	NumForcedGC,
	OtherSys,
	StackInuse,
	StackSys,
	Sys,
	PauseTotalNs,
	NumGC,
	NumGoroutine,
	RandomValue gauge

	PollCount count
}

func (m *Metrics) UpdateMetrics() *Metrics {
	var rtm runtime.MemStats

	var PollCount = 0
	m.PollCount = count(PollCount)
	rand.Seed(time.Now().Unix())
	m.RandomValue = gauge(rand.Intn(100) + 1)

	PollCount++
	runtime.ReadMemStats(&rtm)

	m.NumGoroutine = gauge(runtime.NumGoroutine())

	m.Alloc = gauge(rtm.Alloc)
	m.TotalAlloc = gauge(rtm.TotalAlloc)
	m.Sys = gauge(rtm.Sys)
	m.Mallocs = gauge(rtm.Mallocs)
	m.Frees = gauge(rtm.Frees)

	m.LiveObjects = m.Mallocs - m.Frees

	m.PauseTotalNs = gauge(rtm.PauseTotalNs)
	m.NumGC = gauge(rtm.NumGC)

	m.PollCount = count(PollCount)
	rand.Seed(time.Now().Unix())
	m.RandomValue = gauge(rand.Intn(10000) + 1)
	return m

}

func (mertics *Metrics) PostMetrics(httpClient *resty.Client) {

	b, _ := json.Marshal(mertics)
	var inInterface map[string]float64
	json.Unmarshal(b, &inInterface)

	for field, val := range inInterface {
		var uri, mtype, mval string
		if field != "PollCount" {
			mtype = "gauge"
			mval = strconv.FormatFloat(val, 'f', -1, 64)
		} else {
			mtype = "gauge"
			mval = strconv.FormatFloat(val, 'f', -1, 64)
		}

		fmt.Println(uri)
		httpClient.
			SetRetryCount(3).
			SetRetryWaitTime(10 * time.Second)
		resp, err := httpClient.R().
			SetPathParams(map[string]string{
				"host":        "127.0.0.1",
				"port":        strconv.Itoa(8080),
				"metricType":  mtype,
				"metricName":  field,
				"metricValue": mval,
			}).
			SetHeader("Content-Type", "text/plain").
			Post("http://{host}:{port}/update/{metricType}/{metricName}/{metricValue}")
		if err != nil {
		}
		if resp.StatusCode() != 200 {
			errors.New("HTTP Status != 200")
		}
	}
}
