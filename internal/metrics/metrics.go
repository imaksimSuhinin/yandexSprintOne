package metrics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"
)

const (
	MetricTypeGauge   string = "gauge"
	MetricTypeCounter string = "counter"
)

const (
	scheme          string = "http"
	host            string = "127.0.0.1"
	port            string = "8080"
	contentTypeKey  string = "Content-Type"
	contentTypeText string = "text/plain"
	contentTypeJSON string = "application/json"
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

	var resp *http.Response

	for field, val := range inInterface {
		var mkey, mtype, mval string

		if field != "PollCount" {
			mtype = MetricTypeGauge
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
		} else {
			mtype = MetricTypeCounter
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
		}

		url := NewURLConnectionString(scheme, host+":"+port, "/update/"+mtype+"/"+mkey+"/"+mval)
		var req, err = http.NewRequest(http.MethodPost, url, nil)
		req.Header.Add(contentTypeKey, contentTypeText) // добавляем заголовок Accept

		resp, err = httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			defer resp.Body.Close()
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New("HTTP Status != 200")
		}
	}
	defer resp.Body.Close()

	log.Println("Post...")

	return nil
}
func (m *Metrics) PostMetricsJSON(httpClient *http.Client) error {

	OneMetrics := struct {
		ID    string  `json:"id"`
		MType string  `json:"type"`
		Delta int64   `json:"delta"`
		Value float64 `json:"value"`
	}{}

	b, _ := json.Marshal(&m)
	var inInterface map[string]float64
	json.Unmarshal(b, &inInterface)

	var resp *http.Response

	for field, val := range inInterface {
		var mkey, mtype, mval string

		if field != "PollCount" {
			mtype = MetricTypeGauge
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
			OneMetrics.MType = mtype
			OneMetrics.ID = mval

			OneMetrics.Value, _ = strconv.ParseFloat(mval, 64)
			//log.Println("PollCount:" + string(OneMetrics.MType))

		} else {
			mtype = MetricTypeCounter
			mval = strconv.FormatFloat(val, 'f', -1, 64)
			mkey = field
			OneMetrics.MType = mtype
			OneMetrics.ID = mkey
			OneMetrics.Delta, _ = strconv.ParseInt(mval, 10, 64)
		}
		log.Println("PollCount:" + string(OneMetrics.Delta))

		statJSON, _ := json.Marshal(OneMetrics)
		fmt.Println(mkey)
		url := NewURLConnectionString(scheme, host+":"+port, "/update/")
		var req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(statJSON))
		req.Header.Add(contentTypeKey, contentTypeJSON) // добавляем заголовок Accept
		log.Println(statJSON)
		resp, err = httpClient.Do(req)
		if err != nil {
			fmt.Println(err)
			resp.Body.Close()
		}
		if resp.StatusCode != http.StatusOK {
			return errors.New("HTTP Status != 200")
		}
	}
	defer resp.Body.Close()

	log.Println("Post...")

	return nil
}

func NewURLConnectionString(proto, host, path string) string {
	var v = make(url.Values)

	var u = url.URL{
		Scheme:   proto,
		Host:     host,
		Path:     path,
		RawQuery: v.Encode(),
	}
	return u.String()
}
