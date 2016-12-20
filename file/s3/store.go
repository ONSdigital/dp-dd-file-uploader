package s3

import (
	"github.com/ONSdigital/go-ns/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"io"
)

// NewFileStore factory method to initialise AWS S3 classes.
func NewFileStore(AWSRegion string, s3Bucket string) *FileStore {
	return &FileStore{
		Uploader:   s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(AWSRegion)})),
		BucketName: s3Bucket,
	}
}

// FileStore S3 implementation
type FileStore struct {
	Uploader   s3manageriface.UploaderAPI
	BucketName string
}

// SaveFile sends the file from the given reader to S3 under the given filename.
func (s3 FileStore) SaveFile(reader io.Reader, filename string) error {

	result, err := s3.Uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: aws.String(s3.BucketName),
		Key:    aws.String(filename),
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
