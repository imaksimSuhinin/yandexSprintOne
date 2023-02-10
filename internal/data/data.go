package data

import (
	"errors"
	"fmt"
	"strconv"
)

type metricStorage interface {
	Write(key, value string) error
	Read(key string) string
	ReadValue(key string) (string, error)
	UpdateGaugeValue(key string, value float64) error
	UpdateCounterValue(key string, value string) error
}

type DataBase struct {
	data map[string]string
}

type DataStorage struct {
	Data metricStorage
}

func NewDataBase() *DataBase {
	return &DataBase{
		data: make(map[string]string),
	}
}

func InitDatabase() DataStorage {
	var dataStorage DataStorage
	dataStorage.Data = NewDataBase()
	dataStorage.Data.Write("Alloc", "0")
	dataStorage.Data.Write("BuckHashSys", "0")
	dataStorage.Data.Write("Frees", "0")
	dataStorage.Data.Write("GCCPUFraction", "0")
	dataStorage.Data.Write("GCSys", "0")

	dataStorage.Data.Write("HeapAlloc", "0")
	dataStorage.Data.Write("HeapIdle", "0")
	dataStorage.Data.Write("HeapInuse", "0")
	dataStorage.Data.Write("HeapObjects", "0")
	dataStorage.Data.Write("HeapReleased", "0")

	dataStorage.Data.Write("HeapSys", "0")
	dataStorage.Data.Write("LastGC", "0")
	dataStorage.Data.Write("Lookups", "0")
	dataStorage.Data.Write("MCacheInuse", "0")
	dataStorage.Data.Write("MCacheSys", "0")

	dataStorage.Data.Write("MSpanInuse", "0")
	dataStorage.Data.Write("MSpanSys", "0")
	dataStorage.Data.Write("Mallocs", "0")
	dataStorage.Data.Write("NextGC", "0")
	dataStorage.Data.Write("NumForcedGC", "0")

	dataStorage.Data.Write("NumGC", "0")
	dataStorage.Data.Write("OtherSys", "0")
	dataStorage.Data.Write("PauseTotalNs", "0")
	dataStorage.Data.Write("StackInuse", "0")
	dataStorage.Data.Write("StackSys", "0")

	dataStorage.Data.Write("Sys", "0")
	dataStorage.Data.Write("TotalAlloc", "0")
	dataStorage.Data.Write("PollCount", "0")
	dataStorage.Data.Write("RandomValue", "0")

	return dataStorage
}

func (m DataBase) Write(key, value string) error {
	m.data[key] = value
	return nil
}

func (m DataBase) Read(key string) string {
	value, err := m.data[key]
	if !err {
		return "0"
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
