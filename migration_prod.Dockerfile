FROM alpine:3.19.1

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

ADD migrations/migrations-auth migrations/migrations-auth
ADD migrations/migrations-chat-server migrations/migrations-chat-server
ADD migration_prod.sh .
ADD .env.prod .

RUN chmod +x migration_prod.sh

ENTRYPOINT ["bash", "migration_prod.sh"]