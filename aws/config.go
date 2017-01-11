package aws

import (
	"errors"
	"fmt"
	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go/aws"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	bucketName *string
	path       string
	url        *url.URL
	awsRegion  *string
}

func NewAWSConfig(awsRegion string, url *url.URL) *Config {
	urlPath, err := extractURLPath(url)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
	return &Config{
		bucketName: aws.String(url.Host),
		path:       urlPath,
		url:        url,
		awsRegion:  aws.String(awsRegion),
	}
}

func extractURLPath(url *url.URL) (string, error) {
	if url == nil {
		return "", errors.New("Expected valid S3 URL but was nil.")
	}
	return strings.TrimPrefix(url.Path, "/"), nil
}

func (cfg *Config) GetFilePath(filename string) *string {
	return aws.String(fmt.Sprintf("%s/%s", cfg.path, filename))
}

func (cfg *Config) GetBucketName() *string {
	return cfg.bucketName
}

func (cfg *Config) GetS3FileURL(filename string) string {
	if strings.HasSuffix(cfg.url.String(), "/") {
		return cfg.url.String() + filename
	}
	return cfg.url.String() + "/" + filename
}

func (cfg *Config) GetRegion() *string {
	return cfg.awsRegion
}

func (cfg *Config) ToString() string {
	return fmt.Sprintf("{bucketName=%s, path=%s, region=%s, url=%s}", *cfg.bucketName, cfg.path, *cfg.awsRegion, cfg.url)
}
