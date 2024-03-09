#!/bin/bash

docker compose --env-file .env.local -f loc.docker-compose.yaml up --build
#go run <img> -config-path=.env.local
# <img> - название или id образа, -config-path=.env.local  - имя файла с переменными окружения