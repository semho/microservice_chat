package main

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	desc "github.com/semho/microservice_chat/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"log"
	"time"
)

const (
	address = "localhost:50051"
	userID  = 12
)

func getUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*desc.UserResponse, error) {
	response, err := client.Get(ctx, &desc.GetRequest{Id: userID})
	if err != nil {
		return nil, err
	}

	return response.GetUser(), nil
}

func createUser(ctx context.Context, client desc.AuthV1Client) (*desc.CreateResponse, error) {
	userDetail := &desc.UserDetail{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Role:  desc.Role_user.Enum(),
	}
	pass := gofakeit.Password(true, false, false, false, false, 32)
	passwordDetail := &desc.UserPassword{
		Password:        pass,
		PasswordConfirm: pass,
	}
	request := &desc.CreateRequest{
		Detail:   userDetail,
		Password: passwordDetail,
	}

	response, err := client.Create(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func updateUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*emptypb.Empty, error) {
	userUpdate := &desc.UpdateUserInfo{
		Name:  &wrapperspb.StringValue{Value: gofakeit.Name()},
		Email: &wrapperspb.StringValue{Value: gofakeit.Email()},
	}

	request := &desc.UpdateRequest{Id: userID, Info: userUpdate}
	response, err := client.Update(ctx, request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deleteUserByID(ctx context.Context, client desc.AuthV1Client, userID int64) (*emptypb.Empty, error) {
	request := &desc.DeleteRequest{Id: userID}
	response, err := client.Delete(ctx, request)
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

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r, err := getUserByID(ctx, c, userID)
	//r, err := createUser(ctx, c)
	//r, err := updateUserByID(ctx, c, userID)
	r, err := deleteUserByID(ctx, c, userID)
	if err != nil {
		log.Fatalf("failed to get user by id: %v", err)
	}

	log.Printf(color.RedString("Answer: \n"), color.GreenString("%+v", r))
}
