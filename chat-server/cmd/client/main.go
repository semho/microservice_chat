package main

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	desc "github.com/semho/microservice_chat/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

const (
	address = "localhost:50052"
	chatID  = 12
)

func createChat(ctx context.Context, client desc.ChatServerV1Client) (*desc.CreateChatResponse, error) {
	users := make([]*desc.User, 3)

	for i := 0; i < 3; i++ {
		user := &desc.User{
			Name: gofakeit.Name(),
		}
		users[i] = user
	}

	request := &desc.CreateChatRequest{
		Usernames: users,
	}

	response, err := client.CreateChat(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deleteChatByID(ctx context.Context, client desc.ChatServerV1Client, chatID int64) (*emptypb.Empty, error) {
	request := &desc.DeleteRequest{Id: chatID}
	response, err := client.DeleteChat(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func sendMessage(ctx context.Context, client desc.ChatServerV1Client) (*emptypb.Empty, error) {
	request := &desc.SendMessageRequest{
		Text: gofakeit.BeerName(),
		From: gofakeit.Name(),
	}
	response, err := client.SendMessage(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewChatServerV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r, err := createChat(ctx, c)
	//r, err := deleteChatByID(ctx, c, chatID)
	r, err := sendMessage(ctx, c)
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("Answer: \n"), color.GreenString("%+v", r))
}
