package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"records/pkg"
	"records/proto"
	"records/utils"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *GrpcServer
	gateway    *GatewayServer
	Logger     *zap.Logger
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

	jwtSecret := []byte(viper.GetString("jwt.secret"))
	jwtManager := utils.NewJWTManager(jwtSecret)

	grpcServer := NewGrpcServer(tlsConfig, logger, jwtManager.Clone())
	gateway := NewGatewayServer(tlsConfig, logger)

	return &Server{
		grpcServer: grpcServer,
		gateway:    gateway,
		Logger:     logger,
	}, nil
}

func (a *Server) Run() error {
	defer a.Logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connStr := viper.GetString("postgres.connection_string")
	grpcPort := viper.GetString("ports.grpc")
	httpPort := viper.GetString("ports.http")

	pool, err := SetPGConn(ctx, connStr)
	if err != nil {
		return err
	}
	defer pool.Close()

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		return err
	}
	errChan := make(chan error, 2)

	go func() {
		if err := a.grpcServer.Start(grpcLis); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}()

	go func() {
		if err := a.gateway.Start(ctx, httpPort, grpcPort); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	repo := pkg.NewRecordRepository(pool)
	service := pkg.NewRecordService(repo)
	controller := pkg.NewRecordController(service)
	proto.RegisterRecordServiceServer(a.grpcServer.server, controller)

	a.Logger.Info("server started", zap.String("grpc", grpcPort), zap.String("http", httpPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		a.grpcServer.server.GracefulStop()
		return err
	case <-quit:
		a.grpcServer.server.GracefulStop()
	}

	return nil
}
