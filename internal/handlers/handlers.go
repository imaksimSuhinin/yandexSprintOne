package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"yandexSprintOne/internal/converter"
)

type metric struct {
	mtype, value string
}
type metricValue struct {
	val       [8]byte
	isCounter bool
}

var (
	metricMap       = make(map[string]metricValue)
	lastCounterData int64
)

func SaveMetricToFile(m map[string]metricValue) {
	options := os.O_WRONLY | os.O_TRUNC | os.O_CREATE
	file, err := os.OpenFile("metrics.txt", options, os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, m)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func ShowMetrics(w http.ResponseWriter, _ *http.Request) {
	var stringMetricMap metric
	metricStringMap := make(map[string]metric)
	for k, v := range metricMap {
		if !v.isCounter {
			stringMetricMap.mtype = "gauge"
			stringMetricMap.value = strconv.FormatFloat((converter.Float64FromBytes([]byte(v.val[:]))), 'f', -1, 64)
			metricStringMap[k] = stringMetricMap
		} else {
			stringMetricMap.mtype = "counter"
			stringMetricMap.value = strconv.FormatInt((converter.Int64FromBytes([]byte(v.val[:]))), 10)
			metricStringMap[k] = stringMetricMap
		}

	}
	log.Println(stringMetricMap.value)
	Templ, _ := template.ParseFiles("internal/html/index.html")
	w.WriteHeader(http.StatusOK)
	Templ.Execute(w, metricMap)
}

func PostMetricHandler(w http.ResponseWriter, r *http.Request) {
	reqMethod := r.Method
	if reqMethod != "POST" {
		outputMessage := "Only POST method is alload"
		log.Println("Wrong method for the handler. " + outputMessage)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(outputMessage))
		return
	}
	var m metricValue
	vars := mux.Vars(r)
	switch vars["metricType"] {
	case "gauge":
		f, err := strconv.ParseFloat(vars["metricValue"], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		m.val = converter.Float64ToBytes(f)
		m.isCounter = false
		metricMap[vars["metricName"]] = m
		w.WriteHeader(http.StatusOK)
		r.Body.Close()
	case "counter":
		c, err := strconv.ParseInt(vars["metricValue"], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		lastCounterData = lastCounterData + c // Change naming...
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		m.val = converter.Int64ToBytes(lastCounterData)
		m.isCounter = true
		metricMap[vars["metricName"]] = m
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
		r.Body.Close()
	default:
		log.Println("Type", vars["metricType"], "wrong")
		outputMessage := "Type " + vars["metricType"] + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(outputMessage))
		r.Body.Close()
	}
	log.Println(metricMap)
	SaveMetricToFile(metricMap)
}
