package Kafka

import "notificationservice/internal/DTO"

type ProducerInterface interface {
	ProduceMessage(data interface{}) error
	Close() error
}

type ConsumerInterface interface {
	ConsumeMessage() (string, error)
	Close() error
}

type UserNotifier interface {
	NotifyUserCreated(user *DTO.User) error
	NotifyUserUpdated(user *DTO.User) error
}
