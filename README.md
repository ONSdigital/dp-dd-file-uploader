dp-dd-file-uploader
================

A web UI for uploading files into the data discovery project.

### Configuration

| Environment variable | Default | Description
| -------------------- | ------- | -----------
| BIND_ADDR            | :20019           | The host and port to bind to
| KAFKA_ADDR           | localhost:9092   | The address of the Kafka instance
| TOPIC_NAME           | dp-csv-splitter  | The name of the topic to send file uploaded events to
| AWS_REGION           | eu-west-1        | The AWS region the S3 bucket is hosted in
| S3_BUCKET            | file-uploaded    | The name of the S3 bucket to store files.
| UPLOAD_TIMEOUT       | 1m               | The time before an upload times out. Use 'm' for minutes, 's' for seconds etc. Maximum of 1h.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright ©‎ 2016, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
