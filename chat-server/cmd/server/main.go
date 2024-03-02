package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	desc "github.com/semho/microservice_chat/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatServerV1Server
}

func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	if req.GetUsernames() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided")
	}

	log.Printf("User names: %+v", req.GetUsernames())

	newChat := &desc.Chat{
		Id:        int64(gofakeit.Uint64()),
		Usernames: req.GetUsernames(),
	}
	log.Printf("New chat: %+v", newChat)

	response := &desc.CreateChatResponse{
		Id: newChat.Id,
	}

	return response, nil
}

func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	log.Printf("delete chat by id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.GetFrom() == "" || req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided")
	}

	newMessage := &desc.SendMessageRequest{
		From: req.GetFrom(),
		Text: req.GetText(),
	}

	log.Printf("New message: %+v", newMessage)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatServerV1Server(s, &server{})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
