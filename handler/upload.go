package handlers

import (
	"errors"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"github.com/ONSdigital/dp-dd-file-uploader/file"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"github.com/ONSdigital/go-ns/handlers/response"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
	"time"
)

var FileStore file.Store
var EventProducer event.Producer

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

	err = FileStore.SaveFile(file, header.Filename)
	if err != nil {
		log.Error(err, log.Data{"message": FailedToSaveFile})
		err = response.WriteJSON(w, Response{Message: FailedToSaveFile}, http.StatusInternalServerError)
		if err != nil {
			log.Error(err, log.Data{"message": "Failed to write JSON response"})
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	event := event.FileUploaded{
		Time:  time.Now().UTC().Unix(),
		S3URL: config.AWScfg.GetS3FileURL(header.Filename),
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
