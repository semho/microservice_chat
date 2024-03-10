#!/bin/bash

docker compose --env-file .env.prod -f prod.docker-compose.yaml up --build

#docker build -t <img> -f <name>.Dockerfile .
#go run <img> -config-path=.env.local
# <img> - название или id образа, -config-path=.env.local  - имя файла с переменными окружения