FROM golang:1.21.8-alpine AS builder

COPY ./auth /github.com/semho/microservice_chat/auth
WORKDIR /github.com/semho/microservice_chat/auth

# берет из аргумента пайплайна
ARG ENV_FILE
COPY $ENV_FILE .env

RUN go mod download
RUN go build -o ./bin/auth_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/semho/microservice_chat/auth/bin/auth_server .

CMD ["./auth_server", "-config-path", ".env"]