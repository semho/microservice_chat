package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	desc "github.com/semho/microservice_chat/auth/pkg/auth_v1"
	"github.com/semho/microservice_chat/config"
	"github.com/semho/microservice_chat/config/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "../../../.env", "path to config file")
}

type server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

type UserDB struct {
	ID        int64
	Name      string
	Email     string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateUserRequest struct {
	ID    int64
	Name  string
	Email string
}

func (s *server) userExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func buildUpdateQuery(req UpdateUserRequest) (string, []any) {
	var setStatements []string
	var values []any

	if req.Name != "" {
		setStatements = append(setStatements, "name = $"+strconv.Itoa(len(values)+1))
		values = append(values, req.Name)
	}

	if req.Email != "" {
		setStatements = append(setStatements, "email = $"+strconv.Itoa(len(values)+1))
		values = append(values, req.Email)
	}

	query := "UPDATE users SET " + strings.Join(setStatements, ", ") + " WHERE id = $" + strconv.Itoa(len(values)+1)
	return query, append(values, req.ID)
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	var user UserDB
	err = s.pool.QueryRow(ctx, "SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1",
		req.GetId()).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &desc.GetResponse{
		User: &desc.UserResponse{
			Id: req.GetId(),
			Detail: &desc.UserDetail{
				Name:  user.Name,
				Email: user.Email,
				Role:  desc.Role(user.Role - 1).Enum(),
			},
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
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

	log.Printf("User name: %v", req.GetDetail().Name)
	//приводим к строке, иначе не верно будет браться id в enum
	role := fmt.Sprint(req.GetDetail().GetRole())
	roleValue, ok := desc.Role_value[role]
	//костыль для записи в БД, т.к. enum c 0, а в БД с 1
	if !ok {
		roleValue = int32(desc.Role_user) + 1
	} else {
		roleValue++
	}

	row := s.pool.QueryRow(ctx, `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) 
		RETURNING id`, req.GetDetail().Name, req.GetDetail().Email, req.GetPassword().Password, roleValue)
	var userID int64
	if err := row.Scan(&userID); err != nil {
		log.Printf("failed to insert user into the database: %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	newUser := &desc.User{
		Id:       userID,
		Detail:   req.GetDetail(),
		Password: req.GetPassword(),
	}
	log.Printf("user %v", newUser)

	response := &desc.CreateResponse{
		Id: newUser.Id,
	}

	return response, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	updateRequest := UpdateUserRequest{
		ID:    req.GetId(),
		Name:  req.GetInfo().GetName().GetValue(),
		Email: req.GetInfo().GetEmail().GetValue(),
	}
	query, values := buildUpdateQuery(updateRequest)

	res, err := s.pool.Exec(ctx, query, values...)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("Обновлено строк: %d", rowCount)

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

	exists, err := s.userExists(ctx, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	res, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, req.GetId())
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	rowCount := res.RowsAffected()
	log.Printf("удалено строк: %d", rowCount)

	log.Printf("delete user by id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	grpcConfig, err := env.NewGRPCConfig(env.GrpcPortEnvAuth)
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig(env.DSNEnvAuth)
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{pool: pool})

	log.Printf("server listening at: %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
