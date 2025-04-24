package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"records_bank/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GatewayServer struct {
	server *http.Server
	logger *zap.Logger
}

func NewGatewayServer(tlsConfig *tls.Config, logger *zap.Logger) *GatewayServer {
	mux := runtime.NewServeMux()

	err := mux.HandlePath("GET", "/openapiv2/*", openAPIServer("proto/openapiv2"))
	if err != nil {
		logger.Fatal("Failed to register OpenAPI handler", zap.Error(err))
	}
	err = mux.HandlePath("GET", "/v1/records/user", userBankRecordsFilter())
	if err != nil {
		logger.Fatal("Failed to register user records filter", zap.Error(err))
	}

	return &GatewayServer{
		server: &http.Server{
			Handler:   mux,
			TLSConfig: tlsConfig,
		},
		logger: logger,
	}
}

func (g *GatewayServer) Start(ctx context.Context, httpAddr string, grpcAddr string) error {
	creds := credentials.NewTLS(g.server.TLSConfig)

	err := proto.RegisterBankRecordServiceHandlerFromEndpoint(
		ctx,
		g.server.Handler.(*runtime.ServeMux),
		fmt.Sprintf(":%s", grpcAddr),
		[]grpc.DialOption{grpc.WithTransportCredentials(creds)},
	)
	if err != nil {
		return err
	}

	g.server.Addr = fmt.Sprintf(":%s", httpAddr)

	return g.server.ListenAndServeTLS("", "")
}
