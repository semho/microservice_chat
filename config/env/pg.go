package env

import (
	"errors"
	"github.com/semho/microservice_chat/config"
	"os"
)

var _ config.PGConfig = (*pgConfig)(nil)

const dsnEnvName = "AUTH_DSN"

type pgConfig struct {
	dsn string
}

func NewPGConfig() (*pgConfig, error) {
	dsn := os.Getenv(dsnEnvName)
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
