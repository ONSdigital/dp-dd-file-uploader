package main

import (
	"github.com/ONSdigital/dp-dd-file-uploader/assets"
	"github.com/ONSdigital/dp-dd-file-uploader/aws"
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
)

func main() {

	config.Load()
	s3Config := aws.NewAWSConfig(config.AWSRegion, config.S3URL)
	log.Namespace = "dp-dd-file-uploader"

	var err error
	render.Renderer = unrolled.New()
	handlers.FileStore = s3.NewFileStore(s3Config)

	render.Renderer = unrolled.New(unrolled.Options{
		Asset:      assets.Asset,
		AssetNames: assets.AssetNames,
		Funcs: []template.FuncMap{{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		}},
	})

	handlers.EventProducer, err = kafka.NewProducer(config.KafkaAddr, config.TopicName)
	handlers.S3Config = s3Config
	if err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}

	router := pat.New()
	alice := alice.New(
		timeout.Handler(config.UploadTimeout),
		log.Handler,
		requestID.Handler(16),
	).Then(router)

	router.Get("/healthcheck", healthcheck.Handler)
	router.Get("/", handlers.Home)
	router.Post("/", handlers.Upload)

	log.Debug("Starting server", log.Data{"bind_addr": config.BindAddr})

	server := &http.Server{
		Addr:         config.BindAddr,
		Handler:      alice,
		ReadTimeout:  config.UploadTimeout,
		WriteTimeout: config.UploadTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error(err, nil)
		os.Exit(1)
	}
}
