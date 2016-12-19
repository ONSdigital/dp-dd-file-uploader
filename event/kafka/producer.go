package kafka

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"github.com/Shopify/sarama"
)

func NewProducer(kafkaAddress string, topicName string) (*Producer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddress}, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer:  producer,
		TopicName: topicName,
	}, nil
}

// Producer wraps an internal kafka producer
type Producer struct {
	Producer  sarama.SyncProducer
	TopicName string
}

// FileUploaded sends a new event.
func (kafka Producer) FileUploaded(event event.FileUploaded) error {

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: kafka.TopicName,
		Key:   sarama.StringEncoder(event.Filename),
		Value: sarama.ByteEncoder(eventJSON),
	}

	_, _, err = kafka.Producer.SendMessage(producerMsg)
	if err != nil {
		return err
	}

	return nil
}
