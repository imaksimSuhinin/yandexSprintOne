package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func main() {
	//client := resty.New()
	//
	//client.
	//	SetRetryCount(3).
	//	SetRetryWaitTime(10 * time.Second)
	//
	//resp, err := client.R().
	//	SetPathParams(map[string]string{
	//		"host":  "127.0.0.1",
	//		"port":  strconv.Itoa(8080),
	//		"type":  "type",
	//		"name":  "name",
	//		"value": "value",
	//	}).
	//	SetHeader("Content-Type", "text/plain").
	//	Post("http://{host}:{port}/")
	//
	//if err != nil {
	//
	//}
	//if resp.StatusCode() != 200 {
	//	errors.New("HTTP Status != 200")
	//}

	//client := &http.Client{}
	//var body = []byte(`{"message":"Hello"}`)
	//request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/", bytes.NewBuffer(body))
	//if err != nil {
	//	errors.New("HTTP Status != 200")
	//}
	//
	//request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//
	//// отправляем запрос
	//response, err := client.Do(request)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//defer response.Body.Close()
	//
	//_body, _err := io.ReadAll(response.Body)
	//if _err != nil {
	//	fmt.Println(_err)
	//	os.Exit(1)
	//}
	//// и печатаем его
	//fmt.Println(string(_body))

	//client := http.Client{
	//	Timeout: 6 * time.Second,
	//}
	//
	//req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/", nil)
	//
	//req.Header.Add("Content-Type", "text/plain")
	//_, err := client.Do(req)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	client := resty.New()

	resp, err := client.R().
		EnableTrace().
		Get("http://127.0.0.1:8080")

	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}
