package eventtest

import (
	"errors"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"strings"
)

func NewDummyEventProducer() *DummyEventProducer {
	return &DummyEventProducer{}
}

type DummyEventProducer struct {
	Invocations int
}

func (eventProducer *DummyEventProducer) FileUploaded(event event.FileUploaded) error {

	eventProducer.Invocations++

	if strings.Contains(event.Filename, "EventError") {
		return errors.New("Error sending event")
	}

	return nil
}
