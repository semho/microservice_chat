ARG ENV_FILE
FROM golang:1.21.8-alpine AS builder



COPY ./auth /github.com/semho/microservice_chat/auth
WORKDIR /github.com/semho/microservice_chat/auth

# берет из аргумента пайплайна
COPY $ENV_FILE auth.env
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