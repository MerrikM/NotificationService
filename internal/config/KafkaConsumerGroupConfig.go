package config

import (
	"fmt"
	"github.com/IBM/sarama"
)

// ConsumerGroupHandler - структура, реализующая обработчик для Consumer Group
type ConsumerGroupHandler struct {
}

type CustomConsumerGroup struct {
	Group  sarama.ConsumerGroup
	Topics []string
}

func NewCustomConsumerGroup(brokers []string, groupID string, topics []string) (*CustomConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к брокеру: %v", err)
	}

	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания ConsumerGroup: %v", err)
	}

	return &CustomConsumerGroup{
		Group:  consumerGroup,
		Topics: topics,
	}, nil
}
