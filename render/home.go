package render

import (
	"github.com/unrolled/render"
	"io"
	"net/http"
)

type renderer interface {
	HTML(w io.Writer, status int, name string, binding interface{}, htmlOpt ...render.HTMLOptions) error
}

var Renderer renderer

func Home(w io.Writer) error {
	return Renderer.HTML(w, http.StatusAccepted, "index", "")
}
