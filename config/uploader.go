package config

import (
	"github.com/ONSdigital/go-ns/log"
	"os"
)

const bindAddrKey = "BIND_ADDR"
const kafkaAddrKey = "KAFKA_ADDR"
const s3BucketKey = "S3_BUCKET"
const awsRegionKey = "AWS_REGION"
const topicNameKey = "TOPIC_NAME"

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

// AWSAccessKey the AWS access key
var AWSAccessKey = ""

// AWSSecretKey the AWS secret key
var AWSSecretKey = ""

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

	if accessKeyEnv := os.Getenv("ACCESS_KEY"); len(accessKeyEnv) > 0 {
		AWSAccessKey = accessKeyEnv
	}

	if secretKeyEnv := os.Getenv("SECRET_KEY"); len(secretKeyEnv) > 0 {
		AWSSecretKey = secretKeyEnv
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
	})
}
