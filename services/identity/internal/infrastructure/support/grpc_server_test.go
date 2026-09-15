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
