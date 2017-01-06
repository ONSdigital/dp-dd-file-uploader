FROM ubuntu:16.04

WORKDIR /app/

COPY ./build/dp-dd-file-uploader .

ENTRYPOINT ./dp-dd-file-uploader
