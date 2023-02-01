package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	os "yandexSprintOne/internal/os"

	"yandexSprintOne/internal/data"
	"yandexSprintOne/internal/handlers"
)

var (
	httpServer  http.Server
	database    = data.InitDatabase()
	getTemplate = handlers.ParseTemplate("internal/html/index.html")
)

func main() {
	go startServer(database, getTemplate)
	os.UpdateOsSignal()
}

func startServer(database data.DataBase, template *template.Template) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowMetrics(writer, request, template)
	}).Methods("GET")
	r.HandleFunc("/value/{metricType}/{metricName}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShowValue(writer, request, &database)
		})
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.PostMetricHandler(writer, request, &database)
		}).Methods("POST")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
