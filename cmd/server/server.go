package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"yandexSprintOne/internal/handlers"
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
	vars := mux.Vars(r)
	//var m *agent.Metrics
	x := vars["metricName"]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(x))
	//w.Write([]byte("Unknown statName"))
}
