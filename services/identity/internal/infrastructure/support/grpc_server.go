package support

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/codejsha/bookstore-microservices/identity/generated/application/port/pb/userpb"
	"github.com/codejsha/bookstore-microservices/identity/internal/application/port/security"
	"github.com/codejsha/bookstore-microservices/identity/internal/config"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/keycloak"
	"github.com/codejsha/bookstore-microservices/identity/internal/infrastructure/adapter/protosvc"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 10 * time.Second
	grpcKeepaliveMinTime     = 15 * time.Second
	grpcMaxConcurrentStreams = 256
)

type GrpcServer struct {
	server         *grpc.Server
	shutdowner     fx.Shutdowner
	grpcCfg        *config.GrpcConfig
	userGrpcServer *protosvc.UserGrpcServer
}

func NewGrpcServer(
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	grpcCfg *config.GrpcConfig,
	userGrpcServer *protosvc.UserGrpcServer,
	introspector keycloak.Introspector,
	revocation security.RevocationChecker,
	risk security.RiskChecker,
	telemetryManager *TelemetryManager,
) *GrpcServer {
	s := &GrpcServer{
		shutdowner:     shutdowner,
		grpcCfg:        grpcCfg,
		userGrpcServer: userGrpcServer,
	}

	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(telemetryManager.TraceProvider),
			otelgrpc.WithMeterProvider(telemetryManager.MeterProvider),
		)),
		grpc.KeepaliveParams(grpcKeepaliveParams()),
		grpc.KeepaliveEnforcementPolicy(grpcKeepaliveEnforcement()),
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
	}
	if grpcCfg.AuthEnabled {
		authorizer := keycloak.NewTokenAuthorizer(introspector, revocation, risk)
		opts = append(opts, grpc.UnaryInterceptor(newAuthUnaryInterceptor(authorizer)))
		logrus.Info("grpc: application-level bearer authorization enabled")
	}

	s.server = grpc.NewServer(opts...)
	s.registerGrpcServers()
	if grpcCfg.ReflectionEnabled {
		reflection.Register(s.server)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			serverCfg := s.grpcCfg.Server
			addr := net.JoinHostPort(serverCfg.Host, serverCfg.Port)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("listen grpc on %s: %w", addr, err)
			}
			go s.serve(listener)
			return nil
		},
	})

	return s
}

func (s *GrpcServer) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
		return nil
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	}
}

func grpcKeepaliveParams() keepalive.ServerParameters {
	return keepalive.ServerParameters{
		Time:    grpcKeepaliveTime,
		Timeout: grpcKeepaliveTimeout,
	}
}

func grpcKeepaliveEnforcement() keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             grpcKeepaliveMinTime,
		PermitWithoutStream: true,
	}
}

func (s *GrpcServer) serve(listener net.Listener) {
	if err := s.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		logrus.Errorf("grpc server stopped: %v", err)
		s.signalShutdown()
	}
}

func (s *GrpcServer) signalShutdown() {
	if s.shutdowner == nil {
		return
	}
	if err := s.shutdowner.Shutdown(fx.ExitCode(1)); err != nil {
		logrus.Errorf("failed to signal application shutdown: %v", err)
	}
}

func (s *GrpcServer) registerGrpcServers() {
	userpb.RegisterUserServiceServer(s.server, s.userGrpcServer)
}

func newAuthUnaryInterceptor(authorizer *keycloak.TokenAuthorizer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		token, err := bearerFromMetadata(ctx)
		if err != nil {
			return nil, err
		}
		switch _, reason, aerr := authorizer.Authorize(ctx, token, false); reason {
		case keycloak.AuthzOK, keycloak.AuthzRiskRestricted:
			return handler(ctx, req)
		case keycloak.AuthzInactive:
			return nil, status.Error(codes.Unauthenticated, "token is not active")
		case keycloak.AuthzIntrospectUnavailable:
			logrus.WithContext(ctx).WithError(aerr).
				Warn("grpc authz denied: introspection unavailable")
			return nil, status.Error(codes.PermissionDenied, "authorization unavailable")
		case keycloak.AuthzRevoked:
			return nil, status.Error(codes.PermissionDenied, "token has been revoked")
		case keycloak.AuthzRiskBlocked:
			return nil, status.Error(codes.PermissionDenied, "unauthorized")
		case keycloak.AuthzRiskUnavailable:
			logrus.WithContext(ctx).WithError(aerr).
				Warn("grpc authz denied: risk store unavailable")
			return nil, status.Error(codes.PermissionDenied, "authorization unavailable")
		case keycloak.AuthzRevocationUnavailable:
			logrus.WithContext(ctx).WithError(aerr).
				Warn("grpc authz denied: revocation denylist unavailable")
			return nil, status.Error(codes.PermissionDenied, "authorization unavailable")
		default:
			return nil, status.Error(codes.PermissionDenied, "unauthorized")
		}
	}
}

func bearerFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing request metadata")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization token")
	}
	token := keycloak.BearerToken(values[0])
	if token == "" {
		return "", status.Error(codes.Unauthenticated, "malformed authorization token")
	}
	return token, nil
}
