package event

// Producer interface for sending events.
type Producer interface {
	FileUploaded(event FileUploaded) (err error)
}

// FileUploaded event
type FileUploaded struct {
	Time     int64  `json:"time"`
	S3URL    string `json:"s3URL"`
}
