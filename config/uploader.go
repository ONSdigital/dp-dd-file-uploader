package config

import (
	"github.com/ONSdigital/go-ns/log"
	"os"
	"net/url"
	"strings"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"errors"
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
var awsRegion = "eu-west-1"

// TopicName the name of the Kafka topic to send messages to.
var TopicName = "file-uploaded"

var AWScfg *AWSConfig

// Default S3 URL value.
var s3URL, _ = url.Parse("s3://dp-csv-splitter")

type AWSConfig struct {
	bucketName *string
	path string
	url *url.URL
	awsRegion *string
}

func extractURLPath(url *url.URL) (string, error) {
	if url == nil {
		return "", errors.New("Expected valid S3 URL but was nil.")
	}
	return strings.TrimPrefix(url.Path, "/"), nil
}

func (cfg *AWSConfig) GetFilePath(filename string) *string {
	return aws.String(fmt.Sprintf("%s/%s", cfg.path, filename))
}

func (cfg *AWSConfig) GetBucketName() *string {
	return cfg.bucketName
}

func (cfg *AWSConfig) GetS3FileURL(filename string) string {
	if strings.HasSuffix(cfg.url.String(), "/") {
		return cfg.url.String() + filename
	}
	return cfg.url.String() + "/" + filename
}

func (cfg *AWSConfig) GetRegion() *string {
	return cfg.awsRegion
}

func (cfg *AWSConfig) ToString() string {
	return fmt.Sprintf("{bucketName=%s, path=%s, region=%s, url=%s}", *cfg.bucketName,  cfg.path,  *cfg.awsRegion,  cfg.url)
}

func NewAWSConfig(awsRegion string, url *url.URL) *AWSConfig {
	urlPath, err := extractURLPath(url)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	return &AWSConfig{
		bucketName: aws.String(url.Host),
		path: urlPath,
		url: url,
		awsRegion: aws.String(awsRegion),
	}
}

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
		awsRegion = awsRegionEnv
	}

	if s3URLEnv := os.Getenv(s3URLKey); len(s3URLEnv) > 0 {
		var err error
		if s3URL, err = url.Parse(s3URLEnv); err != nil {
			log.Error(err, log.Data{"Failed to parse S3URL env var, will use default.": s3URL})
		}
	}
	AWScfg = NewAWSConfig(awsRegion, s3URL)
}

func Load() {
	// Will call init().
	log.Debug("dp-dd-file-uploader Configuration", log.Data{
		bindAddrKey:  BindAddr,
		kafkaAddrKey: KafkaAddr,
		topicNameKey: TopicName,
		"awsConfig": AWScfg.ToString(),
	})
}