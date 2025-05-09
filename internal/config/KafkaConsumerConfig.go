package config

import (
	"fmt"
	"github.com/IBM/sarama"
	"notificationservice/internal/Kafka"
)

type Consumer struct {
	Client    sarama.Client
	Consumer  sarama.Consumer
	Topic     string
	Partition int32
}

func NewConsumer(brokers []string, topic string, partition int32) (Kafka.ConsumerInterface, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к брокеру: %v", err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания консьюмера: %v", err)
	}

	return &Consumer{
		Client:    client,
		Consumer:  consumer,
		Topic:     topic,
		Partition: partition,
	}, nil
}

func (consumer *Consumer) Close() error {
	if err := consumer.Consumer.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия консьюмера: %v", err)
	}
	if err := consumer.Client.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия клиента: %v", err)
	}
	return nil
}

func (consumer *Consumer) ConsumeMessage() (string, error) {
	partitionConsumer, err := consumer.Consumer.ConsumePartition(consumer.Topic, consumer.Partition, sarama.OffsetNewest)
	if err != nil {
		return "", fmt.Errorf("ошибка в партициях: %v", err)
	}
	message := <-partitionConsumer.Messages()
	if message == nil {
		return "", fmt.Errorf("канал сообщений закрыт")
	}
	fmt.Println("Message: ", string(message.Value))
	return string(message.Value), nil
}
