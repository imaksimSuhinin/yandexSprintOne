package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"yandexSprintOne/internal/data"
	"yandexSprintOne/internal/handlers"
)

func main() {
	database := data.InitDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowMetrics(writer, request)
	})
	r.HandleFunc("/value/{metricType}/{metricName}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShowValue(writer, request, &database)
		})
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.PostMetricHandler(writer, request, &database)
		})
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
