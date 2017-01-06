#!/bin/bash

if [[ $(docker inspect --format="{{ .State.Running }}" dp-dd-file-uploader) == "false" ]]; then
  exit 1;
fi
