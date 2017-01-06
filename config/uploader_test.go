package config

import
(
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"strings"
	"errors"
)

func TestExtractURLPath(t *testing.T) {

	Convey("Given a valid S3 URL to a bucket subdirectory.", t, func() {
		s3, _ := url.Parse("s3://test-bucket/munge")

		Convey("When the path is exratced.", func() {
			urlPath, err := extractURLPath(s3)

			Convey("Then the result is a path of the url.", func() {
				So(urlPath, ShouldEqual, "munge")
			})

			Convey("And the path does not contain a leading forward slash.", func() {
				So(strings.HasPrefix(urlPath, "/"), ShouldBeFalse)
			})

			Convey("And there are no errors.", func() {
				So(err == nil, ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid S3 URL to the root of a bucket.", t, func() {
		s3, _ := url.Parse("s3://test-bucket/")

		Convey("When the file path is exratced from the url.", func() {
			urlPath, err := extractURLPath(s3)

			Convey("Then the path should empty.", func() {
				So(urlPath, ShouldEqual, "")
			})

			Convey("And there are no errors.", func() {
				So(err == nil, ShouldBeTrue)
			})
		})
	})

	Convey("Given a nil S3 URL", t, func() {
		Convey("When the file path is exratced from the url.", func() {
			filePath, err := extractURLPath(nil)

			Convey("Then the appropriate error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, errors.New("Expected valid S3 URL but was nil."))
			})

			Convey("And the file path is empty.", func() {
				So(filePath, ShouldBeBlank)
			})
		})
	})
}
