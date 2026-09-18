package protostub

import (
	"encoding/json"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestKeepaliveParameters_idleConnection_pingedWithBoundedTimeout(t *testing.T) {
	params := keepaliveParameters()

	if params.Time < 10*time.Second {
		t.Fatalf("keepalive time = %v, grpc clamps anything below 10s", params.Time)
	}
	if params.Timeout <= 0 || params.Timeout >= params.Time {
		t.Fatalf("keepalive timeout = %v, want within (0, %v)", params.Timeout, params.Time)
	}
	if !params.PermitWithoutStream {
		t.Fatal("a dead peer must be detected while the connection is idle, not only during an rpc")
	}
}

func TestRoundRobinServiceConfig_backendPods_balancedPerRpc(t *testing.T) {
	var parsed struct {
		LoadBalancingConfig []map[string]json.RawMessage `json:"loadBalancingConfig"`
	}
	if err := json.Unmarshal([]byte(roundRobinServiceConfig), &parsed); err != nil {
		t.Fatalf("service config must be valid json: %v", err)
	}
	if len(parsed.LoadBalancingConfig) != 1 {
		t.Fatalf("load balancing config = %v, want exactly one policy", parsed.LoadBalancingConfig)
	}
	if _, ok := parsed.LoadBalancingConfig[0]["round_robin"]; !ok {
		t.Fatalf("load balancing policy = %v, want round_robin", parsed.LoadBalancingConfig[0])
	}
}

func TestBaseDialOptions_grpcClient_acceptedAtDial(t *testing.T) {
	conn, err := grpc.NewClient("passthrough:///order:9090", baseDialOptions()...)
	if err != nil {
		t.Fatalf("grpc rejected the dial options: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
}
