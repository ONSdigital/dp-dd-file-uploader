package handlers

import (
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

func Home(w http.ResponseWriter, _ *http.Request) {
	err := render.Home(w)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to render home page"})
	}
}
