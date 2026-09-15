package support

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx/fxtest"

	"github.com/codejsha/shared-library-go/pkg/config"
	"github.com/codejsha/shared-library-go/pkg/logging"
)

func newTestGinServer(t *testing.T) *GinServer {
	t.Helper()
	appCfg := &config.AppConfig{Logging: config.LoggingConfig{Level: "error"}}
	serverCfg := &config.ServerConfig{Host: "127.0.0.1", Port: "0", Mode: gin.TestMode}
	return NewGinServer(
		fxtest.NewLifecycle(t),
		serverCfg,
		logging.NewLogHelper(appCfg),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func getHealthPath(s *GinServer, path string) int {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	s.engine.ServeHTTP(rec, req)
	return rec.Code
}

func TestGinServer_healthRoute_returns200(t *testing.T) {
	s := newTestGinServer(t)

	if code := getHealthPath(s, "/health"); code != http.StatusOK {
		t.Fatalf("GET /health = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_readyRoute_nilDataSource_returns503(t *testing.T) {
	s := newTestGinServer(t)

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinServer_readyRoute_checkSucceeds_returns200(t *testing.T) {
	s := newTestGinServer(t)
	s.readyCheck = func(context.Context) error { return nil }

	if code := getHealthPath(s, "/health/ready"); code != http.StatusOK {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_readyRoute_checkFails_returns503(t *testing.T) {
	s := newTestGinServer(t)
	s.readyCheck = func(context.Context) error { return errors.New("ping failed") }

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinRequestDeadlineMiddleware_anyRequest_carriesDeadlineWithinTimeout(t *testing.T) {
	engine := gin.New()
	engine.Use(GinRequestDeadlineMiddleware(requestTimeout))
	var remaining time.Duration
	var hasDeadline bool
	engine.GET("/probe", func(c *gin.Context) {
		var deadline time.Time
		deadline, hasDeadline = c.Request.Context().Deadline()
		remaining = time.Until(deadline)
		c.Status(http.StatusNoContent)
	})

	engine.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/probe", nil))

	if !hasDeadline {
		t.Fatal("request context has no deadline")
	}
	if remaining <= 0 || remaining > requestTimeout {
		t.Fatalf("remaining deadline = %v, want within (0, %v]", remaining, requestTimeout)
	}
}

func TestGinServer_httpServer_boundsReadWriteAndIdle(t *testing.T) {
	s := newTestGinServer(t)

	if s.server.ReadHeaderTimeout != readHeaderTimeout || s.server.ReadTimeout != readTimeout ||
		s.server.WriteTimeout != writeTimeout || s.server.IdleTimeout != idleTimeout {
		t.Fatalf("server timeouts = header %v, read %v, write %v, idle %v",
			s.server.ReadHeaderTimeout, s.server.ReadTimeout, s.server.WriteTimeout, s.server.IdleTimeout)
	}
	if s.server.WriteTimeout <= requestTimeout {
		t.Fatalf("write timeout %v must exceed request timeout %v", s.server.WriteTimeout, requestTimeout)
	}
}

func TestGinServer_readyRoute_drainingWithSucceedingCheck_responds503(t *testing.T) {
	s := newTestGinServer(t)
	s.readyCheck = func(context.Context) error { return nil }

	s.BeginDrain()

	if code := getHealthPath(s, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("GET /health/ready = %d, want %d", code, http.StatusServiceUnavailable)
	}
}

func TestGinServer_healthRoute_draining_staysOK(t *testing.T) {
	s := newTestGinServer(t)

	s.BeginDrain()

	if code := getHealthPath(s, "/health"); code != http.StatusOK {
		t.Fatalf("GET /health = %d, want %d", code, http.StatusOK)
	}
}

func TestGinServer_shutdown_inFlightRequestPastDeadline_closesServer(t *testing.T) {
	s := newTestGinServer(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	defer close(release)
	s.server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		select {
		case <-release:
		case <-r.Context().Done():
		}
	})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	served := make(chan error, 1)
	go func() { served <- s.server.Serve(listener) }()
	requested := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + listener.Addr().String() + "/slow")
		if err == nil {
			_ = resp.Body.Close()
		}
		requested <- err
	}()
	<-entered

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	err = s.Shutdown(ctx)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown err = %v, want %v", err, context.DeadlineExceeded)
	}
	if serveErr := <-served; !errors.Is(serveErr, http.ErrServerClosed) {
		t.Fatalf("Serve err = %v, want %v", serveErr, http.ErrServerClosed)
	}
	select {
	case reqErr := <-requested:
		if reqErr == nil {
			t.Fatal("in-flight request completed, want connection closed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight request still open after Shutdown")
	}
}
