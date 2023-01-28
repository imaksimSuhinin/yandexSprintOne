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
	//	log.Println(vars["metricType"])
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
		err = base.UpdateGaugeValue(vars["metricName"], f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
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

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))

		err = base.UpdateCounterValue(vars["metricName"], vars["metricValue"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			log.Println("dsfsdfsdfsdfsdfsdfsdfsdfsd" + vars["metricName"])
			return
		}
		r.Body.Close()
	default:
		log.Println("Type", vars["metricType"], "wrong")
		outputMessage := "Type " + vars["metricType"] + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(outputMessage))

		r.Body.Close()
	}
	//	log.Println(metricMap)

}

func ShowValue(w http.ResponseWriter, r *http.Request, base *data.DataBase) {
	vars := mux.Vars(r)
	switch vars["metricType"] {
	case "gauge":
		name := vars["metricName"]
		x, err := base.ReadValue(name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown statName"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(x))
		r.Body.Close()
	case "counter":
		x, err := base.ReadValue("PollCount")
		//x, err := base.ReadValue(vars["metricName"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown statName"))
			return
		}
		//	prevValInt, err := strconv.ParseInt(vars["metricName"], 10, 64)
		//if err != nil {
		//	errors.New(" value is not int64")
		//}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(x))
		r.Body.Close()
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Unknown statName"))
		r.Body.Close()
	}
}
