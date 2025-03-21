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
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *GrpcServer
	gateway    *GatewayServer
	logger     *zap.Logger
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

	kafkaAddr := viper.GetString("kafka.addr")
	kafkaTopic := viper.GetString("kafka.topic")
	kafkaWriter := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaAddr),
		Topic:                  kafkaTopic,
		Balancer:               &kafka.LeastBytes{},
		WriteTimeout:           utils.TIMEOUT,
		AllowAutoTopicCreation: true,
	}

	repo := pkg.NewUserRepository(pool)
	service := pkg.NewUserService(repo, tlsConfig, kafkaWriter)
	controller := pkg.NewUserController(service)

	server := &Server{
		grpcServer: grpcServer,
		gateway:    gateway,
		logger:     logger,
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

	errGroup, errCtx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		s.StartCronJob(errCtx)
		return nil
	})

	grpcPort := viper.GetString("ports.grpc")
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return err
	}
	errGroup.Go(func() error {
		return s.grpcServer.Start(grpcLis)
	})

	httpPort := viper.GetString("ports.http")
	errGroup.Go(func() error {
		return s.gateway.Start(ctx, httpPort, grpcPort)
	})

	s.logger.Info("server started", zap.String("grpc", grpcPort), zap.String("http", httpPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-errCtx.Done():
	case <-quit:
		cancel()
	}

	s.logger.Info("server is shutting down...")
	s.Shutdown(errCtx)

	if err := errGroup.Wait(); err != nil && err != grpc.ErrServerStopped && err != http.ErrServerClosed {
		return err
	}

	s.logger.Info("server shutdown successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	defer s.logger.Sync()
	if err := s.gateway.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to shutdown HTTP server:", zap.Error(err))
	}
	s.grpcServer.server.GracefulStop()
	s.pool.Close()
}
