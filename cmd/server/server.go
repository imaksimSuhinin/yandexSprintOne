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
