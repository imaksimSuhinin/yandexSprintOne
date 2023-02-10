package main

import (
	"github.com/go-chi/chi"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/data"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/handlers"
	os "github.com/imaksimSuhinin/yandexSprintOne/internal/os"
	"html/template"
	"log"
	"net/http"
)

const (
	httpServerAddress string = ":8080"
)

var (
	httpServer  http.Server
	database    = data.DataStorage{}
	getTemplate = handlers.ParseTemplate("internal/html/index.html")
)

func main() {
	go startServer(getTemplate)
	os.UpdateOsSignal()
}

func startServer(template *template.Template) {
	r := chi.NewRouter()

	r.MethodFunc(http.MethodGet, "/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowMetrics(writer, request, template)
	})

	r.MethodFunc(http.MethodPost, "/value/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowJSONValue(writer, request)
	})

	r.MethodFunc(http.MethodGet, "/value/{metricType}/{metricName}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShowValue(writer, request)
		})

	r.MethodFunc(http.MethodPost, "/update/{metricType}/{metricName}/{metricValue}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.PostMetricHandler(writer, request)
		})
	
	httpServer := &http.Server{
		Addr:    httpServerAddress,
		Handler: r,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
