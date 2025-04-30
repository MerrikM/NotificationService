package main

import (
	"golang.org/x/net/context"
	"log"
	"net/http"
	"notificationservice/internal/WebSocket"
	"notificationservice/internal/config"
)

var brokers []string = []string{"localhost:9093"}
var topics []string = []string{"user-event"}
var topic = "user-event"
var groupID = "user-notifier-group"

func main() {

	go WebSocket.StartWebSocketServer("8089")

	// Создаём кастомную ConsumerGroup
	consumerGroup, err := config.NewCustomConsumerGroup(brokers, groupID, topics)
	if err != nil {
		log.Fatalf("Ошибка инициализации ConsumerGroup: %v", err)
	}

	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запуск Kafka consumer
	go func() {
		for {
			select {
			case <-context.Done():
				log.Println("Завершение работы consumer...")
				return
			default:
				log.Println("Начинаем потребление сообщений...")
				// Вызываем Consume, который в свою очередь вызывает ConsumeClaim
				err := consumerGroup.Group.Consume(context, consumerGroup.Topics, consumerGroup)
				if err != nil {
					log.Printf("Ошибка Consume(): %v", err)
				}
			}
		}
	}()

	select {}
}

func healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	log.Println("Сервер рабоатет!")
}

func notifyHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotImplemented)
	log.Println("Временная заглушка")
}
