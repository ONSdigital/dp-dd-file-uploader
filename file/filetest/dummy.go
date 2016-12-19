package filetest

import (
	"errors"
	"github.com/ONSdigital/go-ns/log"
	"io"
	"strings"
)

func NewDummyFileStore() *DummyFileStore {
	return &DummyFileStore{}
}

type DummyFileStore struct {
	Invocations int
}

func (fileStore *DummyFileStore) SaveFile(reader io.Reader, filename string) error {

	fileStore.Invocations++

	log.Debug("Save file called.", log.Data{})

	if strings.Contains(filename, "fileSaveError") {
		return errors.New("Error saving file")
	}

	return nil
}
