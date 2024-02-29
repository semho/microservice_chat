package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	desc "github.com/semho/microservice_chat/auth/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	return &desc.GetResponse{
		User: &desc.UserResponse{
			Id: req.GetId(),
			Detail: &desc.UserDetail{
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Role:  desc.Role_admin.Enum(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if req.GetDetail() == nil || req.GetPassword() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Detail and Password must be provided")
	}

	if req.GetPassword().Password != req.GetPassword().PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "Password and Password Confirm do not match")
	}

	log.Printf("User name: %v", req.Detail.Name)

	newUser := &desc.User{
		Id:       int64(gofakeit.Uint64()),
		Detail:   req.GetDetail(),
		Password: req.GetPassword(),
	}
	fmt.Printf("user %v", newUser)

	response := &desc.CreateResponse{
		Id: newUser.Id,
	}

	return response, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	updateUser := &desc.UpdateUserInfo{
		Name:  req.GetInfo().GetName(),
		Email: req.GetInfo().GetEmail(),
	}

	log.Printf("update: %v", updateUser)

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	log.Printf("delete user by id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
