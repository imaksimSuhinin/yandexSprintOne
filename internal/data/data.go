package data

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type DataBase struct {
	data map[string]string
	*sync.RWMutex
}

func NewDataBase() DataBase {
	return DataBase{
		data:    make(map[string]string),
		RWMutex: &sync.RWMutex{},
	}
}

func InitDatabase() DataBase {
	var metricData = NewDataBase()

	metricData.Write("Alloc", "0")
	metricData.Write("BuckHashSys", "0")
	metricData.Write("Frees", "0")
	metricData.Write("GCCPUFraction", "0")
	metricData.Write("GCSys", "0")

	metricData.Write("HeapAlloc", "0")
	metricData.Write("HeapIdle", "0")
	metricData.Write("HeapInuse", "0")
	metricData.Write("HeapObjects", "0")
	metricData.Write("HeapReleased", "0")

	metricData.Write("HeapSys", "0")
	metricData.Write("LastGC", "0")
	metricData.Write("Lookups", "0")
	metricData.Write("MCacheInuse", "0")
	metricData.Write("MCacheSys", "0")

	metricData.Write("MSpanInuse", "0")
	metricData.Write("MSpanSys", "0")
	metricData.Write("Mallocs", "0")
	metricData.Write("NextGC", "0")
	metricData.Write("NumForcedGC", "0")

	metricData.Write("NumGC", "0")
	metricData.Write("OtherSys", "0")
	metricData.Write("PauseTotalNs", "0")
	metricData.Write("StackInuse", "0")
	metricData.Write("StackSys", "0")

	metricData.Write("Sys", "0")
	metricData.Write("TotalAlloc", "0")
	metricData.Write("PollCount", "0")
	metricData.Write("RandomValue", "0")

	return metricData
}

func (m DataBase) Write(key, value string) error {
	m.data[key] = value
	return nil
}

func (m DataBase) Read(key string) string {
	value, err := m.data[key]
	if !err {

		return key
	}
	return value
}

func (m DataBase) UpdateGaugeValue(key string, value float64) error {
	return m.Write(key, fmt.Sprintf("%v", value))
}

func (m DataBase) UpdateCounterValue(key string, value string) error {
	prevVal := m.Read(key)
	prevValInt, err := strconv.ParseInt(prevVal, 10, 64)
	if err != nil {
		return errors.New(" value is not int64")
	}
	lastValInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return errors.New(" value is not int64")
	}
	res := prevValInt + lastValInt
	newValue := fmt.Sprintf("%v", res)
	m.Write(key, newValue)
	return nil
}

func (m DataBase) ReadValue(key string) (string, error) {
	value, err := m.data[key]
	if !err {
		return "", errors.New("Значение не найдено" + key)
	}
	return value, nil
}
