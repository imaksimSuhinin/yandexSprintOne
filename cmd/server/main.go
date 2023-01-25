package main

import "net/http"

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", HelloWorld)
	// конструируем свой сервер
	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	server.ListenAndServe()
}
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello, World</h1>"))
}
