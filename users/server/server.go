package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"users/pkg"
	"users/proto"
	"users/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *GrpcServer
	gateway    *GatewayServer
	Logger     *zap.Logger
	repository pkg.UserRepository
	service    *pkg.UserService
	controller *pkg.UserController
	pool       *pgxpool.Pool
}

func NewServer() (*Server, error) {
	logger, err := NewZapLogger()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	tlsConfig, err := utils.LoadTLSConfig()
	if err != nil {
		return nil, err
	}

	grpcServer := NewGrpcServer(tlsConfig, logger)
	gateway := NewGatewayServer(tlsConfig, logger)

	ctx := context.Background()
	connStr := viper.GetString("postgres.connection_string")
	pool, err := SetPGConn(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	repo := pkg.NewUserRepository(pool)
	service := pkg.NewUserService(repo, tlsConfig)
	controller := pkg.NewUserController(service)

	server := &Server{
		grpcServer: grpcServer,
		gateway:    gateway,
		Logger:     logger,
		repository: repo,
		service:    service,
		controller: controller,
		pool:       pool,
	}

	proto.RegisterUserServiceServer(grpcServer.server, controller)

	return server, nil
}

func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcPort := viper.GetString("ports.grpc")
	httpPort := viper.GetString("ports.http")

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return err
	}
	errChan := make(chan error, 2)

	go func() {
		if err := s.grpcServer.Start(grpcLis); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}()

	go func() {
		if err := s.gateway.Start(ctx, httpPort, grpcPort); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	s.Logger.Info("server started", zap.String("grpc", grpcPort), zap.String("http", httpPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		s.Shutddown()
		return err
	case <-quit:
		s.Shutddown()
	}

	return nil
}

func (s *Server) Shutddown() {
	defer s.Logger.Sync()
	defer s.pool.Close()
	s.grpcServer.server.GracefulStop()
	s.gateway.server.Shutdown(context.Background())
	s.Logger.Info("server shutdown...")
}
