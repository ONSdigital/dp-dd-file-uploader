package s3

import (
	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"io"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
)

// NewFileStore factory method to initialise AWS S3 classes.
func NewFileStore(awsCfg *config.AWSConfig) *FileStore {
	return &FileStore{
		Uploader:   s3manager.NewUploader(session.New(&aws.Config{Region: awsCfg.GetRegion()})),
		AWSConfig: awsCfg,
	}
}

// FileStore S3 implementation
type FileStore struct {
	Uploader  s3manageriface.UploaderAPI
	AWSConfig *config.AWSConfig
}

// SaveFile sends the file from the given reader to S3 under the given filename.
func (fs FileStore) SaveFile(reader io.Reader, filename string) error {
	result, err := fs.Uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: fs.AWSConfig.GetBucketName(),
		Key:    fs.AWSConfig.GetFilePath(filename),
	})
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to upload"})
		return err
	}

	log.Debug("Upload successful", log.Data{
		"uploadLocation": result.Location,
	})
	return nil
}
