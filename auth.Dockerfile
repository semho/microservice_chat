ARG ENV_FILE
# берет из аргумента пайплайна
COPY $ENV_FILE /github.com/semho/microservice_chat/auth/.env
# Проверяем содержимое файла .env
RUN cat /github.com/semho/microservice_chat/auth/.env
FROM golang:1.21.8-alpine AS builder

COPY ./auth /github.com/semho/microservice_chat/auth
WORKDIR /github.com/semho/microservice_chat/auth

RUN go mod download
RUN go build -o ./bin/auth_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/semho/microservice_chat/auth/bin/auth_server .
COPY --from=builder /github.com/semho/microservice_chat/auth/.env .

# Проверяем содержимое файла .env
RUN cat .env

ENTRYPOINT ["./auth_server"]