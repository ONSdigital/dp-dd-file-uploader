package handlers

import (
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"net/http"
)

func Home(w http.ResponseWriter, req *http.Request) {
	render.Home(w)
}
