package data

type DataBase struct {
	data map[string]string
}

func NewDataBase() *DataBase {
	return &DataBase{
		data: make(map[string]string),
	}
}
func (m DataBase) Write(key, value string) error {
	m.data[key] = value
	return nil
}

func (m DataBase) Read(key string) string {
	value, err := m.data[key]
	if !err {
		return ""
	}
	return value
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

	return *metricData
}
