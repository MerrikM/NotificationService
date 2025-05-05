package Consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"notificationservice/internal/DTO"
	"notificationservice/internal/Kafka"
	"notificationservice/internal/WebSocket"
)

type Handler struct{}

func (c *Handler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (c *Handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		log.Printf("Получено сообщение: %s", string(msg.Value))

		var user DTO.User
		err := json.Unmarshal(msg.Value, &user)
		if err != nil {
			fmt.Errorf("ошибка десериализации сообщения: %v", err)
			continue
		}

		notification := WebSocket.Notification{
			UserID:  user.ID,
			Message: user.Event,
		}

		log.Printf("Отправка уведомления в канал для userID=%d: %+v", notification.UserID, notification)
		WebSocket.Broadcast <- notification

		session.MarkMessage(msg, "")
	}
	return nil
}

func ConsumeMessage(consumer Kafka.ConsumerInterface) (*DTO.User, error) {
	message, err := consumer.ConsumeMessage()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения топика: %v", err)
	}

	// Десереализуем юзера
	var user DTO.User
	err = json.Unmarshal([]byte(message), &user)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации пользователя: %v", err)
	}

	if user.ID == 0 {
		return nil, errors.New("получен пользователь с нулевым ID")
	}

	return &user, nil
}
