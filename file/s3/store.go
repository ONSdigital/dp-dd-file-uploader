package s3

import (
	"github.com/ONSdigital/dp-dd-file-uploader/aws"
	"github.com/ONSdigital/go-ns/log"
	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"io"
)

// NewFileStore factory method to initialise AWS S3 classes.
func NewFileStore(s3Config *aws.Config) *FileStore {
	return &FileStore{
		Uploader: s3manager.NewUploader(session.New(&awsSDK.Config{Region: s3Config.GetRegion()})),
		S3Config: s3Config,
	}
}

// FileStore S3 implementation
type FileStore struct {
	Uploader s3manageriface.UploaderAPI
	S3Config *aws.Config
}

// SaveFile sends the file from the given reader to S3 under the given filename.
func (fs FileStore) SaveFile(reader io.Reader, filename string) error {
	result, err := fs.Uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: fs.S3Config.GetBucketName(),
		Key:    fs.S3Config.GetFilePath(filename),
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
