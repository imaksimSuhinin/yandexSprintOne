package handlers

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"yandexSprintOne/internal/converter"
	"yandexSprintOne/internal/data"
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

func ShowMetrics(w http.ResponseWriter, r *http.Request) {
	var stringMetricMap metric
	vars := mux.Vars(r)
	metricStringMap := make(map[string]metric)
	for k, v := range metricMap {
		if !v.isCounter {
			stringMetricMap.mtype = "gauge"
			stringMetricMap.value = vars["metricValue"]
			metricStringMap[k] = stringMetricMap
		} else {
			stringMetricMap.mtype = "counter"
			stringMetricMap.value = vars["metricValue"]
			metricStringMap[k] = stringMetricMap
		}

	}
	log.Println(stringMetricMap.value)
	Templ, _ := template.ParseFiles("internal/html/index.html")
	w.WriteHeader(http.StatusOK)
	Templ.Execute(w, metricMap)
}

func PostMetricHandler(w http.ResponseWriter, r *http.Request, base *data.DataBase) {
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
		base.UpdateGaugeValue(vars["metricName"], f)
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
		m.val = converter.Int64ToBytes(lastCounterData)
		m.isCounter = true
		metricMap[vars["metricName"]] = m
		base.UpdateCounterValue(vars["metricType"], vars["metricValue"])
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

}

func ShowValue(w http.ResponseWriter, r *http.Request, base *data.DataBase) {
	vars := mux.Vars(r)
	switch vars["metricType"] {
	case "gauge":
		name := vars["metricName"]
		x := base.ReadValue(name)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(x))
		r.Body.Close()
	case "counter":
		x := base.ReadValue(("PollCount"))
		name := vars["metricName"] + x
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
		r.Body.Close()
	default:
		w.Write([]byte("Unknown statName"))
		r.Body.Close()
	}
}
