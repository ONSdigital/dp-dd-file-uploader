package config

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"net/url"
	"os"
	"time"
)

const bindAddrKey = "BIND_ADDR"
const kafkaAddrKey = "KAFKA_ADDR"
const awsRegionKey = "AWS_REGION"
const topicNameKey = "TOPIC_NAME"
const timeoutKey = "UPLOAD_TIMEOUT"
const s3URLKey = "S3_URL"

const maxUploadTimeout = 1 * time.Hour

// BindAddr the address to bind to.
var BindAddr = ":20019"

// KafkaAddr the Kafka address to send messages to.
var KafkaAddr = "localhost:9092"

// AWSRegion the AWS region to use.
var AWSRegion = "eu-west-1"

// TopicName the name of the Kafka topic to send messages to.
var TopicName = "file-uploaded"

// UploadTimeout is the time to allow for an upload to complete. As per
// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/ this will
// be used to set the ReadTimeout, WriteTimeout and go-ns timeout.Handler timeout as the
// upload will encompass all three.
var UploadTimeout = 1 * time.Minute

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

	if timeoutEnv := os.Getenv(timeoutKey); len(timeoutEnv) > 0 {
		var err error
		UploadTimeout, err = time.ParseDuration(timeoutEnv)
		if err == nil && UploadTimeout > maxUploadTimeout {
			err = fmt.Errorf("Upload timeout too large: %v max allowed: %v", UploadTimeout, maxUploadTimeout)
		}
		if err != nil {
			log.Error(err, log.Data{
				"timeout": timeoutEnv,
			})
			os.Exit(1)
		}
	}
}

func Load() {
	// Will call init().
	log.Debug("dp-dd-file-uploader Configuration", log.Data{
		bindAddrKey:  BindAddr,
		kafkaAddrKey: KafkaAddr,
		topicNameKey: TopicName,
		awsRegionKey: AWSRegion,
		timeoutKey:   UploadTimeout,
		s3URLKey:     S3URL,
	})
}
