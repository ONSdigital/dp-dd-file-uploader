package main

import (
	"github.com/ONSdigital/dp-dd-file-uploader/assets"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
	"github.com/ONSdigital/dp-dd-file-uploader/event/kafka"
	"github.com/ONSdigital/dp-dd-file-uploader/file/s3"
	"github.com/ONSdigital/dp-dd-file-uploader/handler"
	"github.com/ONSdigital/dp-dd-file-uploader/render"
	"github.com/ONSdigital/go-ns/handlers/healthcheck"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/timeout"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	unrolled "github.com/unrolled/render"
	"html/template"
	"net/http"
	"os"
	"time"
)

func main() {

	config.Load()
	log.Namespace = "dp-dd-file-uploader"

	bindAddr := os.Getenv("BIND_ADDR")
	if len(bindAddr) == 0 {
		bindAddr = config.BindAddr
	}

	var err error
	render.Renderer = unrolled.New(unrolled.Options{
		Asset:      assets.Asset,
		AssetNames: assets.AssetNames,
		Funcs: []template.FuncMap{{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}},
	})

	handlers.FileStore = s3.NewFileStore(config.AWSRegion, config.S3Bucket)
	handlers.EventProducer, err = kafka.NewProducer(config.KafkaAddr, config.TopicName)
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	router := pat.New()
	alice := alice.New(
		timeout.Handler(10*time.Second),
		log.Handler,
		requestID.Handler(16),
	).Then(router)

	router.Get("/healthcheck", healthcheck.Handler)
	router.Get("/upload_credentials", handlers.GetUploadCredentials)
	router.NewRoute().PathPrefix("/upload/").Handler(http.StripPrefix("/upload/", http.FileServer(http.Dir("static"))))
	//router.Get("/", handlers.Home)
	//router.Post("/", handlers.Upload)

	log.Debug("Starting server", log.Data{"bind_addr": bindAddr})

	server := &http.Server{
		Addr:         bindAddr,
		Handler:      alice,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
}
