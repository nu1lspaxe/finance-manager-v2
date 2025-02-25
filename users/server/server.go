package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"users/pkg"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	echoServer *echo.Echo
	grpcServer *grpc.Server
	logger     *zap.Logger
}

func NewServer(certFile, keyFile string) (*Server, error) {
	logger, err := NewZapLogger()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	tlsConfig, err := LoadTLSConfig(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	echoServer := NewEchoServer(tlsConfig, logger)

	grpcServer := NewGrpcServer(tlsConfig, logger)

	return &Server{
		echoServer: echoServer,
		grpcServer: grpcServer,
		logger:     logger,
	}, nil
}

func (s *Server) Run(certFile, keyFile string) error {
	pool, err := SetPGConn(context.Background())
	if err != nil {
		s.logger.Fatal("failed to set up PostgreSQL connection", zap.Error(err))
	}
	defer pool.Close()

	repo := pkg.NewUserRepository(pool)
	service := pkg.NewUserService(repo)
	pkg.NewUserController(s.grpcServer, service)

	go func() {
		if err := s.echoServer.Server.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
			s.logger.Fatal("failed to listen and serve Echo", zap.Error(err))
		}
	}()

	lis, err := net.Listen("tcp", ":9443")
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
	}
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			s.logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	s.shutdown()
	return nil

}

func (s *Server) shutdown() {
	s.logger.Info("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echoServer.Shutdown(ctx); err != nil {
		s.logger.Error("Error shutting down Echo server", zap.Error(err))
	}

	s.grpcServer.GracefulStop()

	s.logger.Info("Servers shut down successfully")
}

func NewGrpcServer(tlsConfig *tls.Config, logger *zap.Logger) *grpc.Server {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	creds := credentials.NewTLS(tlsConfig)

	g := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.Creds(creds),
	)

	reflection.Register(g)

	return g
}

func NewEchoServer(tlsConfig *tls.Config, logger *zap.Logger) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	e.Use(middleware.Recover())

	if tlsConfig != nil {
		e.Server = &http.Server{
			Addr:      ":8443",
			Handler:   e,
			TLSConfig: tlsConfig,
		}
	}

	return e
}
