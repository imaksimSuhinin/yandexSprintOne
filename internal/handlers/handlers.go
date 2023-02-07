package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/converter"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/data"
	"github.com/imaksimSuhinin/yandexSprintOne/internal/metrics"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	MetricType  string = "metricType"
	MetricName  string = "metricName"
	MetricValue string = "metricValue"
)

var (
	metricMap       = make(map[string]metricValue)
	lastCounterData int64
	database        = data.InitDatabase()
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type metric struct {
	mtype, value string
}

type metricValue struct {
	val       [8]byte
	isCounter bool
}

func ShowMetrics(w http.ResponseWriter, r *http.Request, template *template.Template) {
	var stringMetricMap metric
	vars := chi.URLParam
	metricStringMap := make(map[string]metric)
	for k, v := range metricMap {
		if !v.isCounter {
			stringMetricMap.mtype = metrics.MetricTypeGauge
			stringMetricMap.value = vars(r, MetricValue)
			metricStringMap[k] = stringMetricMap
		} else {
			stringMetricMap.mtype = metrics.MetricTypeCounter
			stringMetricMap.value = vars(r, MetricValue)
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
	vars := chi.URLParam
	switch vars(r, MetricType) {
	case metrics.MetricTypeGauge:
		f, err := strconv.ParseFloat(vars(r, MetricValue), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		m.val = converter.Float64ToBytes(f)
		m.isCounter = false
		metricMap[vars(r, MetricName)] = m

		w.WriteHeader(http.StatusOK)
		err = database.Data.UpdateGaugeValue(vars(r, MetricName), f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		r.Body.Close()
	case metrics.MetricTypeCounter:
		c, err := strconv.ParseInt(vars(r, MetricValue), 10, 64)
		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		lastCounterData = lastCounterData + c // Change naming...
		m.val = converter.Int64ToBytes(lastCounterData)
		m.isCounter = true
		metricMap[vars(r, MetricName)] = m

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
		z := vars(r, MetricValue)
		err = database.Data.UpdateCounterValue(vars(r, MetricName), z)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		r.Body.Close()
	default:
		log.Println("Type", vars(r, MetricType), "wrong")
		outputMessage := "Type " + vars(r, MetricType) + " not supported, only [counter/gauge]"
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(outputMessage))

		r.Body.Close()
	}
	log.Println(metricMap)

}

func ShowValue(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam
	switch vars(r, MetricType) {
	case metrics.MetricTypeGauge:
		name := vars(r, MetricName)
		x, err := database.Data.ReadValue(name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Unknown statName"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(x))
		r.Body.Close()
	case metrics.MetricTypeCounter:

		x, err := database.Data.ReadValue(vars(r, MetricName))
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

func PostJsonMetricHandler(w http.ResponseWriter, r *http.Request) {
	var requestMetric Metrics

	err := json.NewDecoder(r.Body).Decode(&requestMetric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch requestMetric.MType {
	case metrics.MetricTypeGauge:
		w.WriteHeader(http.StatusOK)
		err = database.Data.UpdateGaugeValue(requestMetric.ID, *requestMetric.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Server error"))
			return
		}
		r.Body.Close()
	case metrics.MetricTypeCounter:
		c := math.Round(*requestMetric.Value)
		lastCounterData = lastCounterData + int64(c)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ok"))
		str := string(lastCounterData)
		err = database.Data.UpdateCounterValue(requestMetric.ID, str)

		r.Body.Close()
	default:
		w.WriteHeader(http.StatusNotImplemented)
		r.Body.Close()
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))

}

func ShowJsonValue(w http.ResponseWriter, r *http.Request) {
	var requestJson struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	getValue, err := database.Data.ReadValue(requestJson.ID)

	responseJson := Metrics{
		ID:    requestJson.ID,
		MType: requestJson.Type,
	}
	switch responseJson.MType {
	case metrics.MetricTypeGauge:
		v, err := strconv.ParseFloat(getValue, 64)
		if err != nil {
		}
		responseJson.Value = &v
	case metrics.MetricTypeCounter:
		var counterVal int64
		counterVal, err = strconv.ParseInt(getValue, 1, 64)
		responseJson.Delta = &counterVal
	}
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
