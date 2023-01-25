package main

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	localConst "yandexSprintOne/internal/const"
	localMetrics "yandexSprintOne/internal/metrics"
)

func main() {
	StartClient()

	tickerStatisticsRefresh := time.NewTicker(2 * time.Second)
	tickerStatisticsUpload := time.NewTicker(10 * time.Second)
	var metrics = localMetrics.Metrics{}
	select {
	case timeTickerRefresh := <-tickerStatisticsRefresh.C:
		log.Println("Refresh")
		log.Println(timeTickerRefresh)
		updateMetrics(metrics)

	case timeTickerUpload := <-tickerStatisticsUpload.C:
		log.Println(timeTickerUpload)

		getMetrics(metrics, localConst.ClientEndpoint+localConst.Port)
		log.Println("Upload")

	}

}

func StartClient() {
	var metrics = localMetrics.Metrics{}
	data := setData()
	client := resty.New()

	client.
		SetRetryCount(1).
		SetRetryWaitTime(10 * time.Second).
		SetRetryMaxWaitTime(90 * time.Second)

	err, response := setHeader(localConst.ClientEndpoint+localConst.Port, data, client.GetClient())

	readDataFromResponse(err, response)
	getMetrics(metrics, localConst.ClientEndpoint+localConst.Port)
	defer response.Body.Close()
}

func updateMetrics(metrics localMetrics.Metrics) {
	metrics.UpdateMetrics(1)

}
func getMetrics(metrics localMetrics.Metrics, endpoint string) {
	metrics.PostMetrics(endpoint, 1)
}

func setHeader(endpoint string, data url.Values, client *http.Client) (error, *http.Response) {
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	// отправляем запрос и получаем ответ
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Статус-код ", response.Status)
	return err, response
}

func setData() url.Values {
	data := url.Values{}
	reader := localMetrics.Metrics{}.Mallocs
	s := fmt.Sprintf("%f", reader)
	data.Set("url", s)
	return data
}

func readDataFromResponse(err error, response *http.Response) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("data:", string(body))
}
