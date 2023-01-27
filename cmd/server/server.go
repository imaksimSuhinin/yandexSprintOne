package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	data "yandexSprintOne/internal/data"
	"yandexSprintOne/internal/handlers"
)

func main() {
	metricData := data.InitDatabase()
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.ShowMetrics)
	r.HandleFunc("/value/{metricType}/{metricName}",
		func(writer http.ResponseWriter, request *http.Request) {
			ShowValue(writer, request, &metricData)
		})
	r.HandleFunc("/update/{metricType}/{metricName}/{metricValue}", handlers.PostMetricHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func ShowValue(w http.ResponseWriter, r *http.Request, m *data.DataBase) {
	vars := mux.Vars(r)

	oldValueInt, err := strconv.ParseInt(m.Read("PollCount"), 10, 64)
	if err != nil {
		errors.New("MemStats value is not int64")
	}
	value, err := strconv.ParseInt(vars["metricName"], 10, 64)
	newValue := fmt.Sprintf("%v", oldValueInt+value)
	m.Write("PollCount", newValue)

	//m.Write("PollCount", vars["metricName"])
	var x = m.Read("PollCount")
	//x := m.GetData("PollCount")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(x))
	//w.Write([]byte("Unknown statName"))
}
