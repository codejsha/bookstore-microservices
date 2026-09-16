package support

import (
	"context"
	"errors"
	"net"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type recorder struct {
	mu    sync.Mutex
	names []string
}

func (r *recorder) record(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.names = append(r.names, name)
}

func (r *recorder) snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.names)
}

type fakeCloser struct {
	closed atomic.Bool
	err    error
}

func (c *fakeCloser) Close() error {
	c.closed.Store(true)
	return c.err
}

func TestRunShutdownSteps_earlierStepsFailOrOverrun_runsEveryStepInOrder(t *testing.T) {
	rec := &recorder{}
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })

	steps := []shutdownStep{
		{name: "fails", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("fails")
			return errors.New("boom")
		}},
		{name: "ignores-deadline", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("ignores-deadline")
			<-release
			return nil
		}},
		{name: "honors-deadline", budget: 20 * time.Millisecond, run: func(ctx context.Context) error {
			rec.record("honors-deadline")
			<-ctx.Done()
			return ctx.Err()
		}},
		{name: "panics", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("panics")
			panic("step panic")
		}},
		{name: "last", budget: 20 * time.Millisecond, run: func(context.Context) error {
			rec.record("last")
			return nil
		}},
	}

	finished := make(chan struct{})
	go func() {
		runShutdownSteps(steps)
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("runShutdownSteps did not finish within its step budgets")
	}

	want := []string{"fails", "ignores-deadline", "honors-deadline", "panics", "last"}
	if got := rec.snapshot(); !slices.Equal(got, want) {
		t.Fatalf("step order = %v, want %v", got, want)
	}
}

func TestRunShutdownSteps_slowEarlierStep_laterStepKeepsFullBudget(t *testing.T) {
	const laterBudget = 80 * time.Millisecond
	var remaining time.Duration
	var hasDeadline bool
	var startErr error

	runShutdownSteps([]shutdownStep{
		{name: "slow", budget: 40 * time.Millisecond, run: func(ctx context.Context) error {
			<-ctx.Done()
			time.Sleep(20 * time.Millisecond)
			return ctx.Err()
		}},
		{name: "later", budget: laterBudget, run: func(ctx context.Context) error {
			startErr = ctx.Err()
			var deadline time.Time
			deadline, hasDeadline = ctx.Deadline()
			remaining = time.Until(deadline)
			return nil
		}},
	})

	if startErr != nil {
		t.Fatalf("later step context already done: %v", startErr)
	}
	if !hasDeadline {
		t.Fatal("later step context has no deadline")
	}
	if remaining <= 50*time.Millisecond || remaining > laterBudget {
		t.Fatalf("later step remaining budget = %v, want within (50ms, %v]", remaining, laterBudget)
	}
}

func TestShutdownSequence_customerComponents_ordersStepsAndBudgets(t *testing.T) {
	steps := shutdownSequence(newTestGinServer(t), &fakeCloser{}, nil, &TelemetryManager{}, &fakeCloser{})

	type spec struct {
		name   string
		budget time.Duration
	}
	want := []spec{
		{"http-drain-delay", drainStepBudget},
		{"http-shutdown", httpShutdownTimeout},
		{"side-effects-drain", sideEffectsDrainTimeout},
		{"kafka-publisher-close", kafkaCloseTimeout},
		{"grpc-clients-close", grpcClientCloseTimeout},
		{"readiness-db-close", readinessCloseTimeout},
		{"telemetry-shutdown", telemetryShutdownTimeout},
	}
	got := make([]spec, len(steps))
	var total time.Duration
	for i, step := range steps {
		got[i] = spec{step.name, step.budget}
		total += step.budget
	}
	if !slices.Equal(got, want) {
		t.Fatalf("shutdown sequence = %v, want %v", got, want)
	}
	if total >= ShutdownStopTimeout {
		t.Fatalf("summed step budgets %v must stay below stop timeout %v", total, ShutdownStopTimeout)
	}
	if drainStepBudget <= shutdownDrainDelay {
		t.Fatalf("drain step budget %v must exceed drain delay %v", drainStepBudget, shutdownDrainDelay)
	}
}

func TestShutdownSequence_stepsAfterDrain_closeEveryComponent(t *testing.T) {
	publisher := &fakeCloser{}
	grpcClients := []*fakeCloser{{}, {err: errors.New("close failed")}, {}, {}}
	telemetryCalled := false
	telemetry := &TelemetryManager{shutdown: func(context.Context) error {
		telemetryCalled = true
		return nil
	}}

	steps := shutdownSequence(newTestGinServer(t), publisher, nil, telemetry,
		grpcClients[0], grpcClients[1], grpcClients[2], grpcClients[3])
	runShutdownSteps(steps[1:])

	if !publisher.closed.Load() {
		t.Error("kafka publisher was not closed")
	}
	for i, c := range grpcClients {
		if !c.closed.Load() {
			t.Errorf("grpc client %d was not closed", i)
		}
	}
	if !telemetryCalled {
		t.Error("telemetry was not shut down")
	}
}

func TestGinServer_shutdownDeadlineExpires_forceClosesActiveConnections(t *testing.T) {
	s := newTestGinServer(t)
	entered := make(chan struct{})
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	s.engine.GET("/block", func(c *gin.Context) {
		close(entered)
		<-release
		c.Status(http.StatusOK)
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- s.server.Serve(listener) }()

	clientErr := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + listener.Addr().String() + "/block")
		if err == nil {
			_ = resp.Body.Close()
		}
		clientErr <- err
	}()
	<-entered

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := s.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown = %v, want %v", err, context.DeadlineExceeded)
	}

	select {
	case err := <-clientErr:
		if err == nil {
			t.Fatal("in-flight request completed, want connection closed by fallback Close")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight connection was not force-closed")
	}
	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		t.Fatalf("Serve = %v, want %v", err, http.ErrServerClosed)
	}
}
