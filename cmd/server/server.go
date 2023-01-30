package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/xlab/closer"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
	"yandexSprintOne/internal/data"
	"yandexSprintOne/internal/handlers"
)

var httpServer http.Server

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	closer.DebugSignalSet = []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	}
}

func main() {
	closer.Bind(Exit)
	database := data.InitDatabase()
	startServer(database)
}

func startServer(database data.DataBase) {
	r := mux.NewRouter()
	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		handlers.ShowMetrics(writer, request)
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

func Exit() {
	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	httpServer.Shutdown(gracefulCtx)
	log.Println("Exit...")
	os.Exit(0)
}
