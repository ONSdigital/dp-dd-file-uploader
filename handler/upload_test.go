package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ONSdigital/dp-dd-file-uploader/event/eventtest"
	"github.com/ONSdigital/dp-dd-file-uploader/file/filetest"
	"github.com/ONSdigital/dp-dd-file-uploader/handler"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	. "github.com/smartystreets/goconvey/convey"
	unrolled "github.com/unrolled/render"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	render.Renderer = unrolled.New(unrolled.Options{
		Directory: "../templates",
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

	Convey("Handler returns 200 status code response when request body is a valid file", t, func() {
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
		So(recorder.Code, ShouldEqual, 200)
		So(fileStore.Invocations, ShouldEqual, 1)
	})

	Convey("Handler returns 500 status code response when file save fails.", t, func() {
		fileStore := filetest.NewDummyFileStore()
		handlers.FileStore = fileStore
		eventProducer := eventtest.NewDummyEventProducer()
		handlers.EventProducer = eventProducer

		recorder := httptest.NewRecorder()

		// set a known filename allowing the file save error to be returned.
		requestBody := strings.Replace(exampleMultipartBody, "AF001EW.csv", "fileSaveError.csv", 1)
		requestBodyReader := bytes.NewReader([]byte(requestBody))
		request, err := http.NewRequest("POST", "/", requestBodyReader)
		request.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryezYpRsrGowIiw0K4")
		So(err, ShouldBeNil)

		handlers.Upload(recorder, request)

		var response = &handlers.Response{}
		json.Unmarshal([]byte(recorder.Body.String()), response)

		fmt.Println(recorder.Body)
		So(recorder.Code, ShouldEqual, 500)
		So(fileStore.Invocations, ShouldEqual, 1)
		So(response.Message, ShouldEqual, handlers.FailedToSaveFile)
	})

	Convey("Handler returns 500 status code response when event send fails.", t, func() {
		fileStore := filetest.NewDummyFileStore()
		handlers.FileStore = fileStore
		eventProducer := eventtest.NewDummyEventProducer()
		handlers.EventProducer = eventProducer

		recorder := httptest.NewRecorder()

		// set a known filename allowing the event error to be returned.
		requestBody := strings.Replace(exampleMultipartBody, "AF001EW.csv", "EventError.csv", 1)
		requestBodyReader := bytes.NewReader([]byte(requestBody))
		request, err := http.NewRequest("POST", "/", requestBodyReader)
		request.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundaryezYpRsrGowIiw0K4")
		So(err, ShouldBeNil)

		handlers.Upload(recorder, request)

		var response = &handlers.Response{}
		json.Unmarshal([]byte(recorder.Body.String()), response)

		fmt.Println(recorder.Body)
		So(recorder.Code, ShouldEqual, 500)
		So(fileStore.Invocations, ShouldEqual, 1)
		So(response.Message, ShouldEqual, handlers.FailedToSendEvent)
	})
}
