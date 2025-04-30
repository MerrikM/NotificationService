package Consumer

import (
	"encoding/json"
	"errors"
	"fmt"
	"notificationservice/internal/DTO"
	"notificationservice/internal/config"
)

func ConsumeMessage(consumer config.ConsumerInterface) (*DTO.User, error) {
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
