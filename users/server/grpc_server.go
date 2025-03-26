package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"users/utils"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	server *grpc.Server
	logger *zap.Logger
}

func NewGrpcServer(tlsConfig *tls.Config, logger *zap.Logger, jwtManager *utils.JWTManager) *GrpcServer {
	creds := credentials.NewTLS(tlsConfig)
	authInterceptor := NewAuthInterceptor(jwtManager)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryLoggerInterceptor(logger),
			authInterceptor.Unary(),
		),
		grpc.ChainStreamInterceptor(
			StreamLoggerInterceptor(logger),
			authInterceptor.Stream(),
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

func UnaryLoggerInterceptor(l *zap.Logger) grpc.UnaryServerInterceptor {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	loggerFunc := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})

	return logging.UnaryServerInterceptor(loggerFunc, opts...)
}

func StreamLoggerInterceptor(l *zap.Logger) grpc.StreamServerInterceptor {
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}

	loggerFunc := logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})

	return logging.StreamServerInterceptor(loggerFunc, opts...)
}

type AuthInterceptor struct {
	jwtManager   *utils.JWTManager
	serviceRoles map[string][]string
}

func NewAuthInterceptor(jwtManager *utils.JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, serviceRoles()}
}

func serviceRoles() map[string][]string {
	const servicePath = "/users.UserService/"

	return map[string][]string{
		servicePath + "GetUser":           {"admin", "user"},
		servicePath + "GetAllUsers":       {"admin"},
		servicePath + "UpdateUser":        {"admin", "user"},
		servicePath + "DeleteUser":        {"admin", "user"},
		servicePath + "AddUserAccount":    {"admin", "user"},
		servicePath + "GetUserAccounts":   {"admin", "user"},
		servicePath + "DeleteUserAccount": {"admin", "user"},
		servicePath + "Logout":            {"admin", "user"},
	}
}

func (a *AuthInterceptor) authorize(ctx context.Context, method string) error {
	accessibleRoles, ok := a.serviceRoles[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := a.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	for _, role := range accessibleRoles {
		if role == claims.Role {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		err := a.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := a.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}
