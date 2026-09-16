package support

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func startTestGrpcServer(t *testing.T) (*GrpcServer, string) {
	t.Helper()
	s := &GrpcServer{server: grpc.NewServer()}
	healthpb.RegisterHealthServer(s.server, health.NewServer())
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() { _ = s.server.Serve(ln) }()
	return s, ln.Addr().String()
}

func TestGrpcServerShutdown_streamOutlivesDeadline_forcesStop(t *testing.T) {
	s, addr := startTestGrpcServer(t)
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer func() { _ = conn.Close() }()
	stream, err := healthpb.NewHealthClient(conn).Watch(context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("watch: %v", err)
	}
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("first watch response: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if err := s.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown error = %v, want %v", err, context.DeadlineExceeded)
	}

	recvErr := make(chan error, 1)
	go func() {
		_, err := stream.Recv()
		recvErr <- err
	}()
	select {
	case err := <-recvErr:
		if err == nil {
			t.Fatal("stream still delivering after forced stop")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("stream still open after forced stop")
	}
}

func TestGrpcServerShutdown_noActiveStreams_stopsGracefully(t *testing.T) {
	s, _ := startTestGrpcServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown error = %v, want nil", err)
	}
}

func TestGrpcKeepaliveParams_idlePeer_pingedWithBoundedTimeout(t *testing.T) {
	params := grpcKeepaliveParams()

	if params.Time != grpcKeepaliveTime || params.Timeout != grpcKeepaliveTimeout {
		t.Fatalf("keepalive params = time %v, timeout %v", params.Time, params.Timeout)
	}
	if params.Timeout >= params.Time {
		t.Fatalf("keepalive timeout %v must be shorter than the ping interval %v", params.Timeout, params.Time)
	}
	if grpcMaxConcurrentStreams == 0 {
		t.Fatal("max concurrent streams must be bounded")
	}
}

func TestGrpcKeepaliveEnforcement_clientPings_toleratedWithoutStreams(t *testing.T) {
	policy := grpcKeepaliveEnforcement()

	if !policy.PermitWithoutStream {
		t.Fatal("clients keep pinging while idle, so the server must permit pings without streams")
	}
	if policy.MinTime <= 0 || policy.MinTime > grpcKeepaliveTime {
		t.Fatalf("enforcement min time = %v, want within (0, %v]", policy.MinTime, grpcKeepaliveTime)
	}
}

func TestGrpcServer_serve_listenerBroken_signalsApplicationShutdown(t *testing.T) {
	stub := &stubShutdowner{}
	s := &GrpcServer{server: grpc.NewServer(), shutdowner: stub}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	s.serve(listener)

	if calls := stub.Calls(); calls != 1 {
		t.Fatalf("shutdown signals = %d, want 1: a dead listener must terminate the app", calls)
	}
}

func TestGrpcServer_serve_gracefulStop_appKeptRunning(t *testing.T) {
	stub := &stubShutdowner{}
	s := &GrpcServer{server: grpc.NewServer(), shutdowner: stub}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	done := make(chan struct{})
	go func() {
		s.serve(listener)
		close(done)
	}()

	s.server.GracefulStop()
	<-done

	if calls := stub.Calls(); calls != 0 {
		t.Fatalf("shutdown signals = %d, want 0 for a graceful stop", calls)
	}
}
