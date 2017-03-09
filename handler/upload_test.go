package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/dp-dd-file-uploader/assets"
	"github.com/ONSdigital/dp-dd-file-uploader/aws"
	"github.com/ONSdigital/dp-dd-file-uploader/event/eventtest"
	"github.com/ONSdigital/dp-dd-file-uploader/file/filetest"
	"github.com/ONSdigital/dp-dd-file-uploader/handler"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	. "github.com/smartystreets/goconvey/convey"
	unrolled "github.com/unrolled/render"
	"time"
)

var exampleMultipartBody string = `

------WebKitFormBoundaryezYpRsrGowIiw0K4
Content-Disposition: form-data; name="file"; filename="AF001EW.csv"
Content-Type: text/csv

observation,data_marking,statistical_unit_eng,statistical_unit_cym,measure_type_eng,measure_type_cym,observation_type,empty,obs_type_value,unit_multiplier,unit_of_measure_eng,unit_of_measure_cym,confidentuality,empty1,geographic_area,empty2,empty3,time_dim_item_id,time_dim_item_label_eng,time_dim_item_label_cym,time_type,empty4,statistical_population_id,statistical_population_label_eng,statistical_population_label_cym,cdid,cdiddescrip,empty5,empty6,empty7,empty8,empty9,empty10,empty11,empty12,dim_id_1,dimension_label_eng_1,dimension_label_cym_1,dim_item_id_1,dimension_item_label_eng_1,dimension_item_label_cym_1,is_total_1,is_sub_total_1,dim_id_2,dimension_label_eng_2,dimension_label_cym_2,dim_item_id_2,dimension_item_label_eng_2,dimension_item_label_cym_2,is_total_2,is_sub_total_2,dim_id_3,dimension_label_eng_3,dimension_label_cym_3,dim_item_id_3,dimension_item_label_eng_3,dimension_item_label_cym_3,is_total_3,is_sub_total_3
153223,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,All categories: Residence Type,All categories: Residence Type,,,
118177,,Person,,Count,,,,,,,,,,K04000001,,,,,,,,,,,,,,,,,,,,,Sex,Sex,,All categories: Sex,All categories: Sex,,,,Age,Age,,All categories: Age 16 and over,All categories: Age 16 and over,,,,Residence Type,Residence Type,,Lives in a household,Lives in a household,,,

------WebKitFormBoundaryezYpRsrGowIiw0K4
Content-Disposition: form-data; name="submit"

Upload
------WebKitFormBoundaryezYpRsrGowIiw0K4--

`

func TestUploadHandler(t *testing.T) {
	url, _ := url.Parse("s3://bucket1/dir/test.csv")
	handlers.S3Config = aws.NewAWSConfig("region1", url)

	render.Renderer = unrolled.New(unrolled.Options{
		Asset:      assets.Asset,
		AssetNames: assets.AssetNames,
		Funcs: []template.FuncMap{{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}},
	})

	Convey("Handler returns 400 status code response when request body is empty", t, func() {
		fileStore := filetest.NewDummyFileStore()
		handlers.FileStore = fileStore
		eventProducer := eventtest.NewDummyEventProducer()
		handlers.EventProducer = eventProducer

		recorder := httptest.NewRecorder()
		rdr := bytes.NewReader([]byte(``))
		request, err := http.NewRequest("POST", "/", rdr)
		So(err, ShouldBeNil)

		handlers.Upload(recorder, request)

		var response = &handlers.Response{}
		json.Unmarshal([]byte(recorder.Body.String()), response)

		So(recorder.Code, ShouldEqual, 400)
		So(response.Message, ShouldEqual, handlers.FailedToReadRequest)
	})

	Convey("Handler returns 202 Accepted status code response when request body is a valid file", t, func() {
		fileStore := filetest.NewDummyFileStore()
		handlers.FileStore = fileStore
		eventProducer := eventtest.NewDummyEventProducer()
		handlers.EventProducer = eventProducer

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte(exampleMultipartBody))
		request, err := http.NewRequest("POST", "/", requestBodyReader)
		request.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryezYpRsrGowIiw0K4")
		So(err, ShouldBeNil)

		handlers.Upload(recorder, request)

		var response = &handlers.Response{}
		json.Unmarshal([]byte(recorder.Body.String()), response)

		fmt.Println(recorder.Body)
		So(recorder.Code, ShouldEqual, 202)
		time.Sleep(1 * time.Second)
		So(fileStore.Invocations, ShouldEqual, 1)
	})

}

func TestValidatingReader(t *testing.T) {

	Convey("validatingReader panics when we're given an invalid csv file with too few fields", t, func() {
		invalidCsvFile := "header_1,header_2\n" + "value_1,value_2\n" + "value_1,value_2"
		source := strings.NewReader(invalidCsvFile)

		reader := handlers.CreateValidatingReader(source, "cntxt")

		var buf bytes.Buffer

		_, err := io.Copy(&buf, reader)
		So(err, ShouldNotBeNil)
	})

	Convey("validatingReader returns error when we're given an invalid csv file with mismatched lines", t, func() {
		invalidCsvFile := "header_1,header_2,header_3\n" + "value_1,value_2\n" + "value_1,value_2,value_3"
		source := strings.NewReader(invalidCsvFile)

		reader := handlers.CreateValidatingReader(source, "cntxt")

		var buf bytes.Buffer

		_, err := io.Copy(&buf, reader)
		So(err, ShouldNotBeNil)
	})

	Convey("validatingReader should not panic when we're given a valid csv file", t, func() {
		tempFile := "/tmp/test.csv"
		writer, _ := os.Create(tempFile)
		defer func() {
			writer.Close()
			os.Remove(tempFile)
			r := recover()
			So(r, ShouldBeNil)
		}()
		csvFile := "header_1,header_2,header_3\n" + "value_1,value_2,value_3\n" + "value_1,value_2,value_3"
		source := strings.NewReader(csvFile)

		reader := handlers.CreateValidatingReader(source, "cntxt")

		var buf bytes.Buffer
		_, err := io.Copy(&buf, reader)

		So(buf.String(), ShouldEqual, csvFile)
		So(err, ShouldBeNil)
	})

}
