package s3_test

import (
	"github.com/ONSdigital/dp-dd-file-uploader/file/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"net/url"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
)

type mockUploader struct {
	invocations int
}

func (mockUploader *mockUploader) Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {

	mockUploader.invocations++

	return &s3manager.UploadOutput{
		Location: "",
		UploadID: "",
	}, nil
}

func TestUpload(t *testing.T) {

	Convey("Given a s3FileStore instance with a mock s3 client and valid aws config", t, func() {

		s3URL, _ := url.Parse("s3://dp-csv-splitter/smooosh")

		uploader := mockUploader{}
		awsCFG := config.NewAWSConfig("region1", s3URL)

		s3FileStore := s3.FileStore{
			Uploader:   &uploader,
			AWSConfig: awsCFG,
		}

		Convey("Given a reader with some test data", func() {
			reader := strings.NewReader("this is data")

			Convey("When SaveFile is called", func() {
				s3FileStore.SaveFile(reader, "filename")

				So(uploader.invocations, ShouldEqual, 1)
			})
		})
	})
}

/*func TestFileStore_GetUploadInputKey(t *testing.T) {
	Convey("Given a s3FileStore instance with a valid s3URL that has an empty path", t, func() {
		s3URL, _ := url.Parse("s3://dp-csv-splitter/")

		s3FileStore := s3.FileStore{
			Uploader:   &mockUploader{},
			BucketName: "the bucketname",
			S3URL: s3URL,
		}

		Convey("When GetUploadInputKey in called with a valid filename.", func() {
			filename := "test-file"
			result := s3FileStore.GetUploadInputKey(filename)

			Convey("Then the result should equal the file name.", func() {
				So(result, ShouldEqual, filename)
			})
		})
	})

	Convey("Given a s3FileStore instance with a valid s3URL that has path.", t, func() {
		s3URL, _ := url.Parse("s3://dp-csv-splitter/smoosh")

		s3FileStore := s3.FileStore{
			Uploader:   &mockUploader{},
			BucketName: "the bucketname",
			S3URL: s3URL,
		}

		Convey("When GetUploadInputKey in called with a valid filename.", func() {
			filename := "test-file"
			result := s3FileStore.GetUploadInputKey(filename)

			Convey("Then the result equals $path/$filename.", func() {
				So(result, ShouldEqual, "smoosh/"+ filename)
			})
		})
	})

	Convey("Given a s3FileStore instance with a nil s3URL value.", t, func() {
		s3URL, _ := url.Parse("s3://dp-csv-splitter/smoosh")

		s3FileStore := s3.FileStore{
			Uploader:   &mockUploader{},
			BucketName: "the bucketname",
			S3URL: s3URL,
		}

		Convey("When GetUploadInputKey in called with any value.", func() {
			filename := "test-file"
			result := s3FileStore.GetUploadInputKey(filename)

			Convey("Then an error is returned.", func() {
				So(result, ShouldEqual, "smoosh/"+ filename)
			})
		})
	})
}*/
