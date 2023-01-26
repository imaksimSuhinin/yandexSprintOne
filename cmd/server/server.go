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
	r.HandleFunc("/value/{metricType}/{metricName}", handlers.ShowValue)
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlers.PostMetricHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
