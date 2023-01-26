package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"yandexSprintOne/internal/converter"
	"yandexSprintOne/internal/handlers"
	agent "yandexSprintOne/internal/metrics"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.ShowMetrics)
	r.HandleFunc("/value/{metricType}/{metricName}", ShowValue)
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlers.PostMetricHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowValue(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//x, err := vars["metricType"]
	//if !err {
	//	errors.New("Значение по ключу не найдено")
	//}
	w.WriteHeader(http.StatusOK)
	var m *agent.Metrics

	dst := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	i := converter.Int64ToBytes(int64(m.PollCount))
	copy(dst[:], i[:])
	w.Write(dst)
	log.Println(dst)
}
