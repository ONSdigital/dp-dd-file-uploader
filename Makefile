build:
	go build -o build/dp-dd-file-uploader

debug: build
	HUMAN_LOG=1 ./build/dp-dd-file-uploader

.PHONY: build debug
