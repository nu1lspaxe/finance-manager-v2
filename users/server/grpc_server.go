package server

import (
	"crypto/tls"
	"net"
	"users/utils"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server *grpc.Server
	logger *zap.Logger
}

func NewGrpcServer(tlsConfig *tls.Config, logger *zap.Logger) *GrpcServer {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
	creds := credentials.NewTLS(tlsConfig)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.Creds(creds),
		grpc.ConnectionTimeout(utils.TIMEOUT),
	)

	reflection.Register(server)
	return &GrpcServer{
		server: server,
		logger: logger,
	}
}

func (g *GrpcServer) Start(lis net.Listener) error {
	return g.server.Serve(lis)
}
