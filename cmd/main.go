package main

import (
	"golang.org/x/net/context"
	"log"
	"notificationservice/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var brokers []string = []string{"localhost:9093"}
var topics []string = []string{"user-event"}
var topic = "user-event"
var groupID = "user-notifier-group"

func main() {
	context, cancel := context.WithCancel(context.Background())
	defer cancel()

	var waitGroup sync.WaitGroup

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		if err := config.StartWebSocketServer("8089"); err != nil {
			log.Printf("веб-сокет не стартовал: %v", err)
			cancel()
		}
	}()

	consumerGroup, err := config.NewCustomConsumerGroup(brokers, groupID, topics)
	if err != nil {
		log.Fatalf("Ошибка инициализации ConsumerGroup: %v", err)
	}

	// Запуск Kafka consumer
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		defer func() {
			err := consumerGroup.Group.Close()
			if err != nil {
				log.Printf("Не получилось закрыть ConsumerGroup: %v", err)
			}
		}()
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
					select {
					case <-context.Done():
						return
					case <-time.After(5 * time.Second):
					}
				}
			}
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChannel
		cancel()
	}()

	waitGroup.Wait()
	log.Println("все горутины остановились")
}
