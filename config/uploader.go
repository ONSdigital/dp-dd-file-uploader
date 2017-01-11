package config

import (
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"os"
	"time"
)

const bindAddrKey = "BIND_ADDR"
const kafkaAddrKey = "KAFKA_ADDR"
const s3BucketKey = "S3_BUCKET"
const awsRegionKey = "AWS_REGION"
const topicNameKey = "TOPIC_NAME"
const timeoutKey = "UPLOAD_TIMEOUT"

const maxUploadTimeout = 1 * time.Hour

// BindAddr the address to bind to.
var BindAddr = ":20019"

// KafkaAddr the Kafka address to send messages to.
var KafkaAddr = "localhost:9092"

// S3Bucket the name of the AWS s3 bucket to get the CSV files from.
var S3Bucket = "dp-csv-splitter"

// AWSRegion the AWS region to use.
var AWSRegion = "eu-west-1"

// TopicName the name of the Kafka topic to send messages to.
var TopicName = "file-uploaded"

// UploadTimeout is the time to allow for an upload to complete. As per
// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/ this will
// be used to set the ReadTimeout, WriteTimeout and go-ns timeout.Handler timeout as the
// upload will encompass all three.
var UploadTimeout = 1 * time.Minute

func init() {
	if bindAddrEnv := os.Getenv(bindAddrKey); len(bindAddrEnv) > 0 {
		BindAddr = bindAddrEnv
	}

	if kafkaAddrEnv := os.Getenv(kafkaAddrKey); len(kafkaAddrEnv) > 0 {
		KafkaAddr = kafkaAddrEnv
	}

	if s3BucketEnv := os.Getenv(s3BucketKey); len(s3BucketEnv) > 0 {
		S3Bucket = s3BucketEnv
	}

	if awsRegionEnv := os.Getenv(awsRegionKey); len(awsRegionEnv) > 0 {
		AWSRegion = awsRegionEnv
	}

	if topicNameEnv := os.Getenv(topicNameKey); len(topicNameEnv) > 0 {
		TopicName = topicNameEnv
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
		s3BucketKey:  S3Bucket,
		awsRegionKey: AWSRegion,
		timeoutKey:   UploadTimeout,
	})
}
