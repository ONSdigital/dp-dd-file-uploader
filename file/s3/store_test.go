package s3_test

import (
	"github.com/ONSdigital/dp-dd-file-uploader/file/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
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

	Convey("Given a s3FileStore instance with a mock s3 client", t, func() {

		uploader := mockUploader{}

		s3FileStore := s3.FileStore{
			Uploader:   &uploader,
			BucketName: "the bucketname",
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
