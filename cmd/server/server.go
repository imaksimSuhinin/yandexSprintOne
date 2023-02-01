package main

import (
	"github.com/go-chi/chi"
	"html/template"
	"log"
	"net/http"
	"yandexSprintOne/internal/data"
	"yandexSprintOne/internal/handlers"
	os "yandexSprintOne/internal/os"
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
	r := chi.NewRouter()

	r.Route("/update", func(router chi.Router) {
		r.MethodFunc(http.MethodGet, "/", func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShowMetrics(writer, request, template)
		})

		r.MethodFunc(http.MethodGet, "/value/{metricType}/{metricName}",
			func(writer http.ResponseWriter, request *http.Request) {
				handlers.ShowValue(writer, request, &database)
			})
		r.MethodFunc(http.MethodPost, "/update/{metricType}/{metricName}/{metricValue}",
			func(writer http.ResponseWriter, request *http.Request) {
				handlers.PostMetricHandler(writer, request, &database)
			})
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
