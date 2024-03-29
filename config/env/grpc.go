package env

import (
	"errors"
	"github.com/semho/microservice_chat/config"
	"net"
	"os"
)

var _ config.GRPCConfig = (*grpcConfig)(nil)

const (
	grpcHostEnvName       = "HOST"
	GrpcPortEnvAuth       = "PORT_AUTH"
	GrpcPortEnvChatServer = "PORT_CHAT_SERVER"
)

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig(portEnvName string) (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}
	port := os.Getenv(portEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
