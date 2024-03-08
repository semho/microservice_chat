package test

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	desc "github.com/semho/microservice_chat/chat-server/pkg/chat-server_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

type Server struct {
	desc.UnimplementedChatServerV1Server
	Pool *pgxpool.Pool
}

var internalServerError = status.Error(codes.Internal, "Internal server error")

func checkError(msg string, err error) error {
	log.Printf("%s: %v", msg, err)
	return internalServerError
}

func (s *Server) userExists(ctx context.Context, userName string) (bool, error) {
	var exists bool
	err := s.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "chat-server".public.users  WHERE name = $1)`,
		userName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Server) chatExists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	err := s.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM "chat-server".public.chats  WHERE id = $1)`,
		chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	if req.GetUsernames() == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: UserNames must be provided")
	}
	log.Printf("User names: %+v", req.GetUsernames())

	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)

	usersID := make([]int64, len(req.GetUsernames()))
	for i, user := range req.GetUsernames() {
		exists, err := s.userExists(ctx, user.GetName())
		if err != nil {
			return nil, checkError("Error checking user existence", err)
		}
		if exists {
			err = tx.QueryRow(ctx, `SELECT id FROM "chat-server".public.users WHERE name = $1`,
				user.GetName()).Scan(&usersID[i])
			if err != nil {
				return nil, checkError("Failed to select user from the database", err)
			}
		} else {
			err = tx.QueryRow(ctx, `INSERT INTO "chat-server".public.users (name) VALUES ($1) RETURNING id`,
				user.GetName()).Scan(&usersID[i])
			if err != nil {
				return nil, checkError("Failed to insert user into the database", err)
			}
		}
	}

	var chatID int64
	err = tx.QueryRow(ctx, `INSERT INTO "chat-server".public.chats DEFAULT VALUES RETURNING id`).Scan(&chatID)
	if err != nil {
		return nil, checkError("failed to create chat the database", err)
	}

	for _, userID := range usersID {
		_, err = tx.Exec(ctx, `INSERT INTO "chat-server".public.user_chat (user_id, chat_id) VALUES ($1, $2)`,
			userID, chatID)
		if err != nil {
			return nil, checkError("failed to insert userID and chatID into the database", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	newChat := &desc.Chat{
		Id:        chatID,
		Usernames: req.GetUsernames(),
	}
	log.Printf("New chat: %+v", newChat)

	response := &desc.CreateChatResponse{
		Id: newChat.Id,
	}

	return response, nil
}

func (s *Server) DeleteChat(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: Id must be provided")
	}

	exists, err := s.chatExists(ctx, req.GetId())
	if err != nil {
		return nil, checkError("Error checking chat existence", err)
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "Chat not found")
	}

	log.Printf("delete chat by id: %d", req.GetId())

	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM "chat-server".public.user_chat WHERE chat_id = $1`, req.GetId())
	if err != nil {
		return nil, checkError("failed to delete user_chat records from the database", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM "chat-server".public.messages WHERE chat_id = $1`, req.GetId())
	if err != nil {
		return nil, checkError("failed to delete user_chat records from the database", err)
	}

	_, err = tx.Exec(ctx, `DELETE FROM "chat-server".public.chats WHERE id = $1`, req.GetId())
	if err != nil {
		return nil, checkError("failed to delete chat the database", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if req.GetFrom() == "" || req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request: From or Text must be provided")
	}
	tx, err := s.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, checkError("Error starting transaction", err)
	}
	defer tx.Rollback(ctx)
	var userID int64
	err = tx.QueryRow(ctx, `SELECT DISTINCT id FROM "chat-server".public.users WHERE name = $1`,
		req.GetFrom()).Scan(&userID)
	if err != nil {
		return nil, checkError("Failed to select user from the database", err)
	}
	if userID == 0 {
		return nil, status.Error(codes.NotFound, "User not found")
	}

	var chatID int64
	err = tx.QueryRow(ctx, `SELECT DISTINCT chat_id FROM "chat-server".public.user_chat WHERE user_id = $1`,
		userID).Scan(&chatID)
	if err != nil {
		return nil, checkError("Failed to select chatID from the database", err)
	}
	if chatID == 0 {
		return nil, status.Error(codes.NotFound, "Chat not found")
	}
	log.Printf("UserID: %d, ChatID: %d", userID, chatID)
	_, err = tx.Exec(ctx, `INSERT INTO "chat-server".public.messages (user_id, chat_id, text) VALUES ($1, $2, $3)`,
		userID, chatID, req.GetText())
	if err != nil {
		return nil, checkError("failed to create message in the database", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, checkError("Error committing transaction", err)
	}

	newMessage := &desc.SendMessageRequest{
		From: req.GetFrom(),
		Text: req.GetText(),
	}

	log.Printf("New message: %+v", newMessage)

	return &emptypb.Empty{}, nil
}
