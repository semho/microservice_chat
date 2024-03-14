#!/bin/bash
source .env

sleep 2 && goose -dir "${MIGRATION_AUTH}" postgres "${MIGRATION_AUTH_DSN}" up -v && \
goose -dir "${MIGRATION_CHAT_SERVER}" postgres "${MIGRATION_CHAT_SERVER_DSN}" up -v