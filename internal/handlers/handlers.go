package handlers

import (
	"github.com/go-chi/chi"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/converter"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/data"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

func ShowMetrics(w http.ResponseWriter, r *http.Request, template *template.Template) {
	var stringMetricMap metric
	vars := chi.URLParam
	metricStringMap := make(map[string]metric)
	for k, v := range metricMap {
		if !v.isCounter {
			stringMetricMap.mtype = "gauge"
			stringMetricMap.value = vars(r, "metricValue")
			metricStringMap[k] = stringMetricMap
		} else {
			stringMetricMap.mtype = "counter"
			stringMetricMap.value = vars(r, "metricValue")
			metricStringMap[k] = stringMetricMap
		}

	}
	log.Println(stringMetricMap.value)
	w.WriteHeader(http.StatusOK)
	template.Execute(w, metricMap)
}

func ParseTemplate(path string) *template.Template {
	Temple, _ := template.ParseFiles(path)
	return Temple
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
	vars := chi.URLParam
	switch vars(r, "metricType") {
	case "gauge":
		f, err := strconv.ParseFloat(vars(r, "metricValue"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		m.val = converter.Float64ToBytes(f)
		m.isCounter = false
		metricMap[vars(r, "metricName")] = m

		w.WriteHeader(http.StatusOK)
		err = base.UpdateGaugeValue(vars(r, "metricName"), f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		r.Body.Close()
	case "counter":
		c, err := strconv.ParseInt(vars(r, "metricValue"), 10, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		lastCounterData = lastCounterData + c // Change naming...
		m.val = converter.Int64ToBytes(lastCounterData)
		m.isCounter = true
		metricMap[vars(r, "metricName")] = m

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
		err = base.UpdateCounterValue(vars(r, "metricName"), vars(r, "metricValue"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		r.Body.Close()
	default:
		log.Println("Type", vars(r, "metricType"), "wrong")
		outputMessage := "Type " + vars(r, "metricType") + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(outputMessage))

		r.Body.Close()
	}
	log.Println(metricMap)

}

func ShowValue(w http.ResponseWriter, r *http.Request, base *data.DataBase) {
	vars := chi.URLParam
	switch vars(r, "metricType") {
	case "gauge":
		name := vars(r, "metricName")
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

		x, err := base.ReadValue(vars(r, "metricName"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown statName"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(x))
		r.Body.Close()
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Unknown statName"))
		r.Body.Close()
	}
}
