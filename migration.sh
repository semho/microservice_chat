#!/bin/bash
source .env

# Ожидание доступности PostgreSQL
while ! nc -z pg-auth 5432; do
  >&2 echo "PostgreSQL недоступен - ожидание..."
  sleep 2
done

#while ! nc -z pg-chat-server 5432; do
#  >&2 echo "PostgreSQL недоступен - ожидание..."
#  sleep 2
#done


sleep 2 && goose -dir "${MIGRATION_AUTH}" postgres "${MIGRATION_AUTH_DSN}" up -v
#&& \
#goose -dir "${MIGRATION_CHAT_SERVER}" postgres "${MIGRATION_CHAT_SERVER_DSN}" up -v