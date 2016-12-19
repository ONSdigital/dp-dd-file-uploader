package file

import (
	"io"
)

type Store interface {
	SaveFile(reader io.Reader, filename string) (err error)
}
