package support

import (
	"net/http"
	"testing"

	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
)

func TestNewOpensearchTransport_insecure_boundsDialHandshakeAndResponseHeaders(t *testing.T) {
	transport := newOpensearchTransport(true)

	if transport.TLSHandshakeTimeout != opensearchTLSHandshakeTimeout ||
		transport.ResponseHeaderTimeout != opensearchResponseHeaderTimeout ||
		transport.MaxIdleConnsPerHost != opensearchMaxIdleConnsPerHost {
		t.Fatalf("transport limits = handshake %v, response header %v, idle per host %d",
			transport.TLSHandshakeTimeout, transport.ResponseHeaderTimeout, transport.MaxIdleConnsPerHost)
	}
	if transport.TLSClientConfig == nil || !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("insecure transport must skip certificate verification")
	}
}

func TestNewOpensearchTransport_secure_keepsCertificateVerification(t *testing.T) {
	transport := newOpensearchTransport(false)

	if transport.TLSClientConfig != nil {
		t.Fatal("secure transport must use the default TLS verification")
	}
	if transport.ResponseHeaderTimeout != opensearchResponseHeaderTimeout {
		t.Fatalf("response header timeout = %v, want %v", transport.ResponseHeaderTimeout, opensearchResponseHeaderTimeout)
	}
}

func TestNewOpensearchConfig_clusterOutage_retriesBoundedWithBackoff(t *testing.T) {
	cfg := newOpensearchConfig(&config.OpensearchConfig{
		Scheme: "http", Host: "opensearch", Port: 9200, Index: "books",
	})

	if cfg.Client.DisableRetry {
		t.Fatal("retries must stay enabled for transient upstream failures")
	}
	if cfg.Client.MaxRetries != opensearchMaxRetries {
		t.Fatalf("max retries = %d, want %d", cfg.Client.MaxRetries, opensearchMaxRetries)
	}
	if cfg.Client.RetryBackoff == nil {
		t.Fatal("retry backoff must be set, otherwise retries hammer the cluster with no delay")
	}

	want := []int{http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout}
	if len(cfg.Client.RetryOnStatus) != len(want) {
		t.Fatalf("retry statuses = %v, want %v", cfg.Client.RetryOnStatus, want)
	}
	for i, status := range want {
		if cfg.Client.RetryOnStatus[i] != status {
			t.Fatalf("retry statuses = %v, want %v", cfg.Client.RetryOnStatus, want)
		}
	}
}

func TestOpensearchRetryBackoff_successiveAttempts_growsAndStaysCapped(t *testing.T) {
	if got := opensearchRetryBackoff(1); got != opensearchRetryBaseBackoff {
		t.Fatalf("attempt 1 backoff = %v, want %v", got, opensearchRetryBaseBackoff)
	}
	if got := opensearchRetryBackoff(2); got != 2*opensearchRetryBaseBackoff {
		t.Fatalf("attempt 2 backoff = %v, want %v", got, 2*opensearchRetryBaseBackoff)
	}

	previous := opensearchRetryBackoff(1)
	for attempt := 2; attempt <= 64; attempt++ {
		got := opensearchRetryBackoff(attempt)
		if got < previous {
			t.Fatalf("attempt %d backoff = %v, must not shrink below %v", attempt, got, previous)
		}
		if got > opensearchRetryMaxBackoff {
			t.Fatalf("attempt %d backoff = %v, exceeds cap %v", attempt, got, opensearchRetryMaxBackoff)
		}
		previous = got
	}
	if got := opensearchRetryBackoff(64); got != opensearchRetryMaxBackoff {
		t.Fatalf("late attempt backoff = %v, want the cap %v", got, opensearchRetryMaxBackoff)
	}
}
