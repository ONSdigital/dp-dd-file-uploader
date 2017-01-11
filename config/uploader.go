package config

import (
	"github.com/ONSdigital/go-ns/log"
	"net/url"
	"os"
)

const bindAddrKey = "BIND_ADDR"
const kafkaAddrKey = "KAFKA_ADDR"
const awsRegionKey = "AWS_REGION"
const topicNameKey = "TOPIC_NAME"
const s3URLKey = "S3_URL"

// BindAddr the address to bind to.
var BindAddr = ":20019"

// KafkaAddr the Kafka address to send messages to.
var KafkaAddr = "localhost:9092"

// AWSRegion the AWS region to use.
var AWSRegion = "eu-west-1"

// TopicName the name of the Kafka topic to send messages to.
var TopicName = "file-uploaded"

// Default S3 URL value.
var S3URL, _ = url.Parse("s3://dp-csv-splitter")

func init() {
	if bindAddrEnv := os.Getenv(bindAddrKey); len(bindAddrEnv) > 0 {
		BindAddr = bindAddrEnv
	}

	if kafkaAddrEnv := os.Getenv(kafkaAddrKey); len(kafkaAddrEnv) > 0 {
		KafkaAddr = kafkaAddrEnv
	}

	if topicNameEnv := os.Getenv(topicNameKey); len(topicNameEnv) > 0 {
		TopicName = topicNameEnv
	}

	if awsRegionEnv := os.Getenv(awsRegionKey); len(awsRegionEnv) > 0 {
		AWSRegion = awsRegionEnv
	}

	if s3URLEnv := os.Getenv(s3URLKey); len(s3URLEnv) > 0 {
		var err error
		if S3URL, err = url.Parse(s3URLEnv); err != nil {
			log.Error(err, log.Data{"Failed to parse S3URL env var, will use default.": S3URL})
		}
	}
}

func Load() {
	// Will call init().
	log.Debug("dp-dd-file-uploader Configuration", log.Data{
		bindAddrKey:  BindAddr,
		kafkaAddrKey: KafkaAddr,
		topicNameKey: TopicName,
		s3URLKey:     S3URL,
	})
}
