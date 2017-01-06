package event

// Producer interface for sending events.
type Producer interface {
	FileUploaded(event FileUploaded) (err error)
}

// FileUploaded event
type FileUploaded struct {
	Filename string `json:"filename"`
	Time     int64  `json:"time"`
	S3Path string `json:"s3Path"`
}
