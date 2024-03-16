#!/bin/sh

wait_for_port() {
    host="$1"
    port="$2"
    timeout="${3:-15}"

    echo "Waiting for $host:$port to be available..."
    timeout $timeout sh -c "while ! nc -z $host $port; do sleep 1; done"
}

# Ожидание доступности порта
wait_for_port "pg-chat-server" "5432"
./chat_server -config-path=/root/.env


