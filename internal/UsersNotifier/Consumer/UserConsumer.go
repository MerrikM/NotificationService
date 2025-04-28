package Consumer

import (
	"fmt"
	"notificationservice/internal/config"
)

func ConsumeMessage(consumer config.ConsumerInterface) error {
	err := consumer.ConsumeMessage()
	if err != nil {
		return fmt.Errorf("ошибка чтения топика: %v", err)
	}
	return nil
}
