package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"archive/zip"
	"github.com/ONSdigital/dp-dd-file-uploader/aws"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
	"github.com/ONSdigital/dp-dd-file-uploader/event"
	"github.com/ONSdigital/dp-dd-file-uploader/file"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"github.com/ONSdigital/go-ns/handlers/response"
	"github.com/ONSdigital/go-ns/log"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"bufio"
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

var (
	TooManyFilesInZip = errors.New("More than one file in zip archive")
	InvalidFileInZip  = errors.New("Non-CSV file in zip archive")
)

func Upload(w http.ResponseWriter, req *http.Request) {

	if FileStore == nil {
		log.ErrorR(req, errors.New("The FileStore dependency has not been configured"), nil)
		return
	}

	if EventProducer == nil {
		log.ErrorR(req, errors.New("The EventProducer dependency has not been configured"), nil)
		return
	}

	multipartReader, err := req.MultipartReader()
	if err != nil {
		handleFileReadFailure(w, req, err, nil)
		return
	}

	var part *multipart.Part
	for {
		part, err = multipartReader.NextPart()
		if err != nil {
			handleFileReadFailure(w, req, err, nil)
			return
		}
		if part != nil && part.FormName() == "file" {
			break
		}
	}

	// NB: we will get an io.EOF error above if the part was not found, so part will not be nil here
	defer part.Close()

	tempFile, err := ioutil.TempFile(config.UploadTempDir, "file-upload-")
	if err != nil {
		handleFileReadFailure(w, req, err, tempFile)
		return
	}
	log.DebugR(req, "Writing file upload to temporary file", log.Data{
		"filename": tempFile.Name(),
	})

	bytesWritten, err := io.Copy(tempFile, part)
	if err != nil {
		handleFileReadFailure(w, req, err, tempFile)
		return
	}
	log.DebugR(req, "Successfully wrote file to temporary storage", log.Data{
		"size": bytesWritten,
	})

	// Rewind file to start ready to read and stream to S3
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		handleFileReadFailure(w, req, err, tempFile)
		return
	}

	// Continue upload to S3 in a separate goroutine
	go uploadFileToS3(tempFile, part.FileName(), log.Context(req))

	err = render.Home(w)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to render home page"})
	}
}

func uploadFileToS3(file *os.File, filename string, context string) {
	defer (func() {
		err := file.Close()
		if err != nil {
			log.ErrorC(context, err, log.Data{"filename": file.Name()})
		}

		err = os.Remove(file.Name())
		if err != nil {
			log.ErrorC(context, err, log.Data{"filename": file.Name()})
		}
	})()

	log.DebugC(context, "Streaming file to s3", log.Data{"filename": filename})

	var reader io.Reader = file

	if filepath.Ext(filename) == ".zip" {
		log.DebugC(context, "Zip file detected - decompressing during upload", nil)
		var err error
		reader, filename, err = decompressZipFile(file, context)
		if err != nil {
			log.ErrorC(context, err, log.Data{
				"message": "Unable to decode zip archive",
			})
		}
	}

	byteMarkers := []ByteMarker{}
	reader = CreateValidatingReader(reader, context)
	reader = ByteMarkRecorder(reader, byteMarkers, 2, context)
	err := FileStore.SaveFile(reader, filename)
	if err != nil {
		log.ErrorC(context, err, log.Data{"message": FailedToSaveFile})
		return
	}

	uploadedEvent := event.FileUploaded{
		Time:  time.Now().UTC().Unix(),
		S3URL: S3Config.GetS3FileURL(filename),
	}

	err = EventProducer.FileUploaded(uploadedEvent)
	if err != nil {
		log.ErrorC(context, err, log.Data{"message": FailedToSendEvent})
		return
	}
}

func handleFileReadFailure(w http.ResponseWriter, req *http.Request, err error, tempFile *os.File) {
	log.ErrorR(req, err, log.Data{"message": FailedToReadRequest})
	err = response.WriteJSON(w, Response{Message: FailedToReadRequest}, http.StatusBadRequest)
	if err != nil {
		log.ErrorR(req, err, log.Data{"message": "Failed to write JSON response"})
		w.WriteHeader(http.StatusBadRequest)
	}

	if tempFile != nil {
		err = os.Remove(tempFile.Name())
		if err != nil {
			log.ErrorR(req, err, log.Data{"message": "Unable to remove temporary file", "file": tempFile.Name()})
		}
	}
}

func decompressZipFile(file *os.File, context string) (reader io.Reader, filename string, err error) {
	stat, err := file.Stat()
	if err != nil {
		log.ErrorC(context, err, log.Data{"message": "Unable to determine file size"})
		return
	}
	zipReader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return
	}

	if len(zipReader.File) != 1 {
		err = TooManyFilesInZip
		return
	}

	entry := zipReader.File[0]
	filename = entry.Name
	if filepath.Ext(filename) != ".csv" {
		err = InvalidFileInZip
		return
	}

	reader, err = entry.Open()
	return
}

// CreateValidatingReader creates a reader that will return an error if the stream being read does not represent a valid csv file.
func CreateValidatingReader(sourceReader io.Reader, context string) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	csvReader := csv.NewReader(sourceReader)
	csvWriter := csv.NewWriter(pipeWriter)
	// create a goroutine that will read from the csvReader and close the pipe if an error is returned by csvReader, or the number of fields isn't correct
	go func() {
		rowCount := 0
		for {
			row, err := csvReader.Read()
			if err != nil {
				csvWriter.Flush();
				pipeWriter.CloseWithError(err)
				log.DebugC(context, "Finished Reading file", log.Data{"rowCount": rowCount, "err": err})
				return
			}
			if len(row)%3 != 0 {
				message := fmt.Sprintf("Wrong number of fields in file at row %d - must be a multiple of 3, but was %d", rowCount, len(row))
				csvWriter.Flush();
				pipeWriter.CloseWithError(errors.New(message))
				return
			}
			if rowCount%50000 == 0 {
				log.DebugC(context, "Saving file to S3", log.Data{"rowCount": rowCount})
			}
			if (rowCount > 0 || row[0]!="Observation") {
				csvWriter.Write(row);
				rowCount++
			}
		}
	}()
	return pipeReader
}

type ByteMarker struct {
	BlockNumber int
	FirstRow int
	LastRow int
	FirstByte int
	LastByte int
}

func ByteMarkRecorder(reader io.Reader, byteMarkers []ByteMarker, blockSize int, context string) io.Reader {
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		rowCount := 0
		byteCount := 0
		blockNumber :=1
		currentBlock := ByteMarker{BlockNumber:1, FirstRow:1, FirstByte:0}
		byteReader := bufio.NewReader(reader)
		for {
			line, readErr := byteReader.ReadBytes('\n')
			_, writeErr := pipeWriter.Write(line)
			if writeErr != nil {
				log.ErrorC(context, writeErr, log.Data{"message": "Error writing to pipe"})
				return
			}
			if readErr != nil {
				pipeWriter.CloseWithError(readErr)
				break
			}
			byteCount += len(line)
			rowCount ++
			if rowCount % blockSize == 0 {
				currentBlock.LastRow = rowCount
				currentBlock.LastByte = byteCount-1
				byteMarkers = append(byteMarkers, currentBlock)
				log.DebugC(context, "Sending block", log.Data{"block": currentBlock})
				blockNumber++
				currentBlock = ByteMarker{BlockNumber:blockNumber, FirstRow: rowCount+1, FirstByte:byteCount}
			}
		}
		if (currentBlock.FirstRow <= rowCount) {
			currentBlock.LastRow = rowCount
			currentBlock.LastByte = byteCount
			byteMarkers = append(byteMarkers, currentBlock)
			log.DebugC(context, "Final block", log.Data{"block": currentBlock})
		}
	}()
	return pipeReader
}