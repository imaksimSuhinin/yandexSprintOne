package runtime_loc

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

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

func (m *Metrics) UpdateMetrics(duration int) {
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	var PollCount = 0
	m.PollCount = count(PollCount)
	rand.Seed(time.Now().Unix())
	m.RandomValue = gauge(rand.Intn(100) + 1)
	for {
		<-time.After(interval)
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

	}
}

func (mertics *Metrics) PostMetrics(serverAddr string, duration int) {

	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		b, _ := json.Marshal(mertics)
		var inInterface map[string]float64
		json.Unmarshal(b, &inInterface)

		for field, val := range inInterface {
			var uri string
			if field != "PollCount" {
				uri = "update/gauge/" + field + "/" + strconv.FormatFloat(val, 'f', -1, 64)

			} else {
				uri = "update/counter/" + field + "/" + strconv.FormatFloat(val, 'f', -1, 64)
			}
			//request, err := http.Post(serverAddr+uri, "text/plain", bytes.NewReader([]byte(strconv.FormatFloat(val, 'f', -1, 64))))
			//
			//if err != nil {
			//	log.Fatal(err)
			//}
			//request.Body.Close()

			fmt.Println(uri)
		}
	}
}
