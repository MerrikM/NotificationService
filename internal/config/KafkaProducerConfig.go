package config

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"notificationservice/internal/DTO"
)

type Producer struct {
	Client   sarama.Client
	Producer sarama.SyncProducer
	Topic    string
	Brokers  []string
	Config   *sarama.Config
}

type ProducerInterface interface {
	ProduceMessage(data interface{}) error
	Close() error
}

type UserNotifier interface {
	NotifyUserCreated(user *DTO.User) error
	NotifyUserUpdated(user *DTO.User) error
}

func NewProducer(brokers []string, topic string) (ProducerInterface, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// "распределение по кругу" между партициями.
	// Чтобы равномерно нагружать все партиции
	// и избежать ситуации, когда одна партиция перегружена, а другие пустуют
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка в создании клиента: %v", err)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("ошибка с воздании продюсера: %v", err)
	}

	return &Producer{
		Producer: producer,
		Client:   client,
		Topic:    topic,
		Brokers:  brokers,
		Config:   config,
	}, nil
}

func (producer *Producer) Close() error {
	if err := producer.Producer.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия продюсера: %v", err)
	}
	if err := producer.Client.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия клиента: %v", err)
	}
	return nil
}

func (producer *Producer) ProduceMessage(data interface{}) error {
	var byteData []byte
	var err error

	switch value := data.(type) {
	case []byte:
		byteData = value
	case string:
		byteData = []byte(value)
	default:
		byteData, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("ошибка сериализации данных: %v", err)
		}
	}

	message := &sarama.ProducerMessage{
		Topic: producer.Topic,
		Key:   sarama.ByteEncoder(byteData),
		Value: sarama.ByteEncoder(byteData),
	}

	_, _, err = producer.Producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %v", err)
	}

	fmt.Println("Сообщение - ", message.Value, " отправлено в топик: "+producer.Topic)
	return nil
}
