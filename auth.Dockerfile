FROM golang:1.21.8-alpine AS builder

COPY ./auth /github.com/semho/microservice_chat/auth
WORKDIR /github.com/semho/microservice_chat/auth

RUN go mod download
RUN go build -o ./bin/auth_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/
COPY --from=builder /github.com/semho/microservice_chat/auth/bin/auth_server .

ENTRYPOINT ["./auth_server"]