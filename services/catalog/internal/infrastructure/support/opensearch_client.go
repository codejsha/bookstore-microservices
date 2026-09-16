package support

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	opensearchgo "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
)

const (
	opensearchDialTimeout           = 5 * time.Second
	opensearchKeepAlive             = 30 * time.Second
	opensearchTLSHandshakeTimeout   = 5 * time.Second
	opensearchResponseHeaderTimeout = 10 * time.Second
	opensearchIdleConnTimeout       = 90 * time.Second
	opensearchMaxIdleConnsPerHost   = 16
	opensearchMaxRetries            = 2
	opensearchRetryBaseBackoff      = 200 * time.Millisecond
	opensearchRetryMaxBackoff       = 2 * time.Second
)

type OpensearchClient struct {
	*opensearchapi.Client
}

func NewOpensearchClient(cfg *config.OpensearchConfig) (*OpensearchClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("opensearch config is nil")
	}
	c, err := opensearchapi.NewClient(newOpensearchConfig(cfg))
	if err != nil {
		return nil, err
	}
	return &OpensearchClient{Client: c}, nil
}

func newOpensearchConfig(cfg *config.OpensearchConfig) opensearchapi.Config {
	addr := fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.Port)
	return opensearchapi.Config{
		Client: opensearchgo.Config{
			Addresses:     []string{addr},
			Username:      cfg.Username,
			Password:      cfg.Password,
			Transport:     newOpensearchTransport(cfg.Insecure),
			MaxRetries:    opensearchMaxRetries,
			RetryOnStatus: opensearchRetryStatuses(),
			RetryBackoff:  opensearchRetryBackoff,
		},
	}
}

func opensearchRetryStatuses() []int {
	return []int{
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	}
}

func opensearchRetryBackoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	backoff := opensearchRetryBaseBackoff << (attempt - 1)
	if backoff <= 0 || backoff > opensearchRetryMaxBackoff {
		return opensearchRetryMaxBackoff
	}
	return backoff
}

func newOpensearchTransport(insecure bool) *http.Transport {
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: opensearchDialTimeout, KeepAlive: opensearchKeepAlive}).DialContext,
		ForceAttemptHTTP2:     true,
		TLSHandshakeTimeout:   opensearchTLSHandshakeTimeout,
		ResponseHeaderTimeout: opensearchResponseHeaderTimeout,
		IdleConnTimeout:       opensearchIdleConnTimeout,
		MaxIdleConnsPerHost:   opensearchMaxIdleConnsPerHost,
	}
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return transport
}
