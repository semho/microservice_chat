
FROM golang:1.21.8-alpine AS builder

COPY ./auth /github.com/semho/microservice_chat/auth
WORKDIR /github.com/semho/microservice_chat/auth
ARG ENV_FILE_CONTENTS
# Создаем .env файл на основе переменной окружения ENV_FILE_CONTENTS
RUN echo "$ENV_FILE_CONTENTS" > auth.env
# Проверяем содержимое файла .env
RUN cat auth.env

RUN go mod download
RUN go build -o ./bin/auth_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/semho/microservice_chat/auth/bin/auth_server .
COPY --from=builder /github.com/semho/microservice_chat/auth/auth.env .env

# Проверяем содержимое файла .env
RUN cat .env

ENTRYPOINT ["./auth_server"]