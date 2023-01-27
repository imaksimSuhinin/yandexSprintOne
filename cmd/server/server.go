package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	m.Write("PollCount", vars["metricName"])
	var x = m.Read("PollCount")
	//x := m.GetData("PollCount")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(x)))
	//w.Write([]byte("Unknown statName"))
}
