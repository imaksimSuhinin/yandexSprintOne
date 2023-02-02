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
	database = data.InitDatabase()
	go startServer(database, getTemplate)
	os.UpdateOsSignal()
}

func startServer(database data.DataStorage, template *template.Template) {
	r := chi.NewRouter()

	r.MethodFunc(http.MethodGet, "/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowMetrics(writer, request, template)
	})

	r.MethodFunc(http.MethodGet, "/value/{metricType}/{metricName}",
		func(writer http.ResponseWriter, request *http.Request) {
			handlers.ShowValue(writer, request, database)
		})

	r.Route("/update", func(router chi.Router) {

		r.MethodFunc(http.MethodPost, "/update/{metricType}/{metricName}/{metricValue}",
			func(writer http.ResponseWriter, request *http.Request) {
				handlers.PostMetricHandler(writer, request, &database)
			})
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
