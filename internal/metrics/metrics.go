package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	//"github.com/go-resty/resty/v2"

	//"github.com/go-resty/resty/v2"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

type gauge float64

type counter int64

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

	PollCount counter
}

var PollCount = 0

func (m *Metrics) UpdateMetrics() *Metrics {
	var rtm runtime.MemStats
	PollCount++
	m.PollCount = counter(PollCount)
	rand.Seed(time.Now().Unix())
	m.RandomValue = gauge(rand.Intn(100) + 1)
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
	m.PollCount = counter(PollCount)
	rand.Seed(time.Now().Unix())
	m.RandomValue = gauge(rand.Intn(10000) + 1)
	log.Println("refresh...")
	return m
}

func (m *Metrics) PostMetrics(httpClient *http.Client) error {
	b, _ := json.Marshal(&m)
	var inInterface map[string]float64
	json.Unmarshal(b, &inInterface)

	for field, val := range inInterface {
		var uri, mkey, mtype, mval string

		if field != "PollCount" {
			mtype = "gauge"
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
		} else {
			mtype = "counter"
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
		}
		fmt.Println(uri, mtype, mval)
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		host := "127.0.0.1"
		port := "8080"

		// добавляем заголовки
		var req, err = http.NewRequest("POST", "http://"+host+":"+port+"/update/"+mtype+"/"+mkey+"/"+mval, nil)

		// добавляем заголовки
		req.Header.Add("Content-Type", "text/plain") // добавляем заголовок Accept

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)

		}
		if resp.StatusCode != 200 {
			return errors.New("HTTP Status != 200")
		}
	}
	log.Println("Post...")
	return nil
}
