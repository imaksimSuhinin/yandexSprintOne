package main

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

func main() {
	client := resty.New()

	client.
		SetRetryCount(3).
		SetRetryWaitTime(10 * time.Second)

	resp, err := client.R().
		SetPathParams(map[string]string{
			"host":  "127.0.0.1",
			"port":  strconv.Itoa(8080),
			"type":  "type",
			"name":  "name",
			"value": "value",
		}).
		SetHeader("Content-Type", "text/plain").
		Post("http://{host}:{port}/")

	if err != nil {

	}
	if resp.StatusCode() != 200 {
		errors.New("HTTP Status != 200")
	}
}
