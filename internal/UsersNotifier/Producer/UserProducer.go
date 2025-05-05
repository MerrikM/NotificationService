package Producer

import (
	"notificationservice/internal/DTO"
	"notificationservice/internal/Kafka"
)

type UserProducer struct {
	producer Kafka.ProducerInterface
	topic    string
}

func NewUserProducer(prodcuer Kafka.ProducerInterface, topic string) *UserProducer {
	return &UserProducer{
		producer: prodcuer,
		topic:    topic,
	}
}

func (userProducer *UserProducer) NotifyUserCreated(user *DTO.User) error {
	return nil
}

func (userProducer *UserProducer) NotifyUserUpdated(user *DTO.User) error {
	return nil
}
