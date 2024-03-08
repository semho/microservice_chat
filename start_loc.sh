#!/bin/bash

docker compose --env-file .env.local -f loc.docker-compose.yaml up --build
#go run main.go -config-path=../../../.env.local