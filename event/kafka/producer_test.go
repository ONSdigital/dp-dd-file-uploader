package kafka_test

import (
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"github.com/ONSdigital/dp-dd-file-uploader/event/kafka"
	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	kafkaConfig.Producer.Return.Successes = true

	kafkaTopicName := "fileUploaded"

	filename := "exampleFilename.csv"
	time := time.Now()

	mockProducer := mocks.NewSyncProducer(t, kafkaConfig)
	mockProducer.ExpectSendMessageWithCheckerFunctionAndSucceed(func(val []byte) error {

		var event event.FileUploaded
		json.Unmarshal(val, &event)

		if event.Filename != filename {
			return errors.New("Filename was not added to the message.")
		}
		if event.Time != time.UTC().Unix() {
			return errors.New("Time was not added to the message.")
		}

		return nil
	})

	Convey("Given a mock producer with a single expected input that succeeds", t, func() {

		var eventProducer event.Producer = kafka.Producer{
			Producer:  mockProducer,
			TopicName: kafkaTopicName,
		}

		Convey("When the producer is called", func() {
			eventProducer.FileUploaded(event.FileUploaded{
				Filename: filename,
				Time:     time.UTC().Unix(),
			})
		})
	})
}
