FROM onsdigital/dp-go

WORKDIR /app/

COPY ./build/dp-dd-file-uploader .

ENTRYPOINT ./dp-dd-file-uploader
