package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"users/pkg"
	"users/proto"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server *grpc.Server
	logger *zap.Logger
}

type HttpGateway struct {
	mux    *runtime.ServeMux
	logger *zap.Logger
}

type Application struct {
	grpcServer *GrpcServer
	gateway    *HttpGateway
	Logger     *zap.Logger
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
	)

	reflection.Register(server)
	return &GrpcServer{
		server: server,
		logger: logger,
	}
}

func NewHttpGateway(logger *zap.Logger) *HttpGateway {
	return &HttpGateway{
		mux:    runtime.NewServeMux(),
		logger: logger,
	}
}

func NewApplication() (*Application, error) {
	logger, err := NewZapLogger()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	tlsConfig, err := LoadTLSConfig()
	if err != nil {
		return nil, err
	}

	grpcServer := NewGrpcServer(tlsConfig, logger)
	gateway := NewHttpGateway(logger)

	return &Application{
		grpcServer: grpcServer,
		gateway:    gateway,
		Logger:     logger,
	}, nil
}

func (g *GrpcServer) Start(lis net.Listener) error {
	return g.server.Serve(lis)
}

func (h *HttpGateway) Start(ctx context.Context, httpAddr string, grpcAddr string) error {
	tlsConfig, err := LoadTLSConfig()
	if err != nil {
		return err
	}

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            tlsConfig.RootCAs,
	})

	err = proto.RegisterUserServiceHandlerFromEndpoint(
		ctx,
		h.mux,
		fmt.Sprintf(":%s", grpcAddr),
		[]grpc.DialOption{grpc.WithTransportCredentials(creds)},
	)
	if err != nil {
		return err
	}

	h.mux.HandlePath("GET", "/openapiv2/*", openAPIServer("proto/openapiv2"))

	httpServer := &http.Server{
		Addr:      fmt.Sprintf(":%s", httpAddr),
		Handler:   h.mux,
		TLSConfig: tlsConfig,
	}

	return httpServer.ListenAndServeTLS("", "")
}

func (a *Application) Run() error {
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

	repo := pkg.NewUserRepository(pool)
	service := pkg.NewUserService(repo)
	controller := pkg.NewUserController(service)
	proto.RegisterUserServiceServer(a.grpcServer.server, controller)

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
