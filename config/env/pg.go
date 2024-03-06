package env

import (
	"errors"
	"github.com/semho/microservice_chat/config"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const (
	DSNEnvAuth       = "AUTH_DSN"
	DSNEnvChatServer = "CHAT_SERVER_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPGConfig(dsnName string) (*pgConfig, error) {
	dsn := os.Getenv(dsnName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
