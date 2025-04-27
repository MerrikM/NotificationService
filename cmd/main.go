package main

import (
	"fmt"
	"log"
	"net/http"
	"notificationservice/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/notify", notifyHandler)

	log.Println("Сервер запущен на порту: " + cfg.Server.Port)
	err = http.ListenAndServe(":"+cfg.Server.Port, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера: ", err)
	}
}

func healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	log.Println("Сервер рабоатет!")
}

func notifyHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	log.Println("Временная заглушка")
}
