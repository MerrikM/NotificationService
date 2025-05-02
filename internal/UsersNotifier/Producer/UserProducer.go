package Producer

import (
	"notificationservice/internal/DTO"
	"notificationservice/internal/Event"
	"notificationservice/internal/config"
)

type UserProducer struct {
	producer config.ProducerInterface
	topic    string
}

func NewUserProducer(prodcuer config.ProducerInterface, topic string) *UserProducer {
	return &UserProducer{
		producer: prodcuer,
		topic:    topic,
	}
}

func (userProducer *UserProducer) NotifyUserCreated(user *DTO.User) error {
	event := &Event.UserEvent{
		EventType: "user_created",
		User:      user,
	}
	return userProducer.producer.ProduceMessage(event)
}

func (userProducer *UserProducer) NotifyUserUpdated(user *DTO.User) error {
	event := &Event.UserEvent{
		EventType: "user_updated",
		User:      user,
	}
	return userProducer.producer.ProduceMessage(event)
}
