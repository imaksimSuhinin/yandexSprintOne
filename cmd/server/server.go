package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	updateOsSignal()
}

func updateOsSignal() {
	sigChanel := make(chan os.Signal, 1)
	signal.Notify(sigChanel)
	exitChanel := make(chan int)
	s := <-sigChanel
	handleOsSignal(s)
	exitCode := <-exitChanel
	os.Exit(exitCode)
}

func handleOsSignal(signal os.Signal) {
	if signal == syscall.SIGTERM {
		fmt.Println("Got kill signal. ")
		fmt.Println("Program will terminate now.")
		os.Exit(0)
	} else if signal == syscall.SIGINT {
		fmt.Println("Got CTRL+C signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else if signal == syscall.SIGQUIT {
		fmt.Println("Got Quit signal")
		fmt.Println("Closing.")
		os.Exit(0)
	} else {
		fmt.Println("Ignoring signal: ", signal)
	}
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
