package config

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"notificationservice/internal/DTO"
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

func (c *CustomConsumerGroup) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *CustomConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *CustomConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		log.Printf("Получено сообщение: %s", string(msg.Value))

		var user DTO.User
		err := json.Unmarshal(msg.Value, &user)
		if err != nil {
			fmt.Errorf("ошибка десериализации сообщения: %v", err)
			continue
		}
		
		notification := Notification{
			UserID:  user.ID,
			Message: user.Event,
		}

		log.Printf("Отправка уведомления в канал для userID=%d: %+v", notification.UserID, notification)
		Broadcast <- notification

		session.MarkMessage(msg, "")
	}
	return nil
}
