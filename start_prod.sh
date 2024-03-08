#!/bin/bash

docker compose --env-file .env.prod -f prod.docker-compose.yaml up --build
#go run main.go -config-path=../../../.prod.local