package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ONSdigital/dp-dd-file-uploader/aws"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"github.com/ONSdigital/dp-dd-file-uploader/file"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"github.com/ONSdigital/go-ns/handlers/response"
	"github.com/ONSdigital/go-ns/log"
)

var FileStore file.Store
var EventProducer event.Producer
var S3Config *aws.Config

type Response struct {
	Message string `json:"message,omitempty"`
}

var FailedToReadRequest string = "Failed to read upload file from the request."
var FailedToSaveFile string = "Failed to save the given file."
var FailedToSendEvent string = "Failed to send file uploaded event."

func Upload(w http.ResponseWriter, req *http.Request) {

	if FileStore == nil {
		log.Error(errors.New("The FileStore dependency has not been configured"), nil)
		return
	}

	if EventProducer == nil {
		log.Error(errors.New("The EventProducer dependency has not been configured"), nil)
		return
	}

	file, header, err := req.FormFile("file")
	if err != nil {
		log.Error(err, log.Data{"message": FailedToReadRequest})
		err = response.WriteJSON(w, Response{Message: FailedToReadRequest}, http.StatusBadRequest)
		if err != nil {
			log.Error(err, log.Data{"message": "Failed to write JSON response"})
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Error(err, nil)
		}
	}()

	log.Debug("Attempting to read file from request", log.Data{"filename": header.Filename})

	reader := CreateValidatingReader(file)

	err = FileStore.SaveFile(reader, header.Filename)
	if err != nil {
		log.Error(err, log.Data{"message": FailedToSaveFile})
		err = response.WriteJSON(w, Response{Message: fmt.Sprintf("%s %v", FailedToSaveFile, err)}, http.StatusInternalServerError)
		if err != nil {
			log.Error(err, log.Data{"message": "Failed to write JSON response"})
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	event := event.FileUploaded{
		Time:  time.Now().UTC().Unix(),
		S3URL: S3Config.GetS3FileURL(header.Filename),
	}

	log.Debug("Sending file to S3 ", log.Data{"url": event.S3URL})

	err = EventProducer.FileUploaded(event)
	if err != nil {
		log.Error(err, log.Data{"message": FailedToSendEvent})
		response.WriteJSON(w, Response{Message: FailedToSendEvent}, http.StatusInternalServerError)
		if err != nil {
			log.Error(err, log.Data{"message": "Failed to write JSON response"})
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = render.Home(w)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to render home page"})
	}

}

// CreateValidatingReader creates a reader that will return an error if the stream being read does not represent a valid csv file.
func CreateValidatingReader(sourceReader io.Reader) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	tee := io.TeeReader(sourceReader, pipeWriter)
	csvReader := csv.NewReader(tee)
	// create a goroutine that will read from the csvReader and close the pipe if an error is returned by csvReader, or the number of fields isn't correct
	go func() {
		stop := false
		i := 0
		for !stop {
			row, err := csvReader.Read()
			if err != nil {
				stop = true
				pipeWriter.CloseWithError(err)
				return
			}
			if len(row)%3 != 0 {
				stop = true
				message := fmt.Sprintf("Wrong number of fields in file - must be a multiple of 3, but was %d", len(row))
				pipeWriter.CloseWithError(errors.New(message))
				return
			}
			i++
		}
		log.Debug("", log.Data{"message": fmt.Sprintf("Valid CSV file has %d rows", i)})
	}()
	return pipeReader
}
