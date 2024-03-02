LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-auth-api
	make generate-chat-api

generate-auth-api:
	mkdir -p auth/pkg/auth_v1
	protoc --proto_path=auth/api/auth_v1 \
    	--go_out=auth/pkg/auth_v1 --go_opt=paths=source_relative \
    	--plugin=protoc-gen-go=./bin/protoc-gen-go \
    	--go-grpc_out=auth/pkg/auth_v1 --go-grpc_opt=paths=source_relative \
    	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
    	auth/api/auth_v1/auth.proto

generate-chat-server-api:
	mkdir -p chat-server/pkg/chat-server_v1
	protoc --proto_path=chat-server/api/chat-server_v1 \
    	--go_out=chat-server/pkg/chat-server_v1 --go_opt=paths=source_relative \
    	--plugin=protoc-gen-go=./bin/protoc-gen-go \
    	--go-grpc_out=chat-server/pkg/chat-server_v1 --go-grpc_opt=paths=source_relative \
    	--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
    	chat-server/api/chat-server_v1/chat-server.proto
