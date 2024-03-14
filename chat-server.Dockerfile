FROM golang:1.21.8-alpine AS builder

COPY ./chat-server /github.com/semho/microservice_chat/chat-server
WORKDIR /github.com/semho/microservice_chat/chat-server

# берет из аргумента пайплайна
ARG ENV_FILE
COPY $ENV_FILE .env

RUN go mod download
RUN go build -o ./bin/chat_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/semho/microservice_chat/chat-server/bin/chat_server .

CMD ["./chat_server", "-config-path", ".env"]