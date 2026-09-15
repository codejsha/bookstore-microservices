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
)

type OpensearchClient struct {
	*opensearchapi.Client
}

func NewOpensearchClient(cfg *config.OpensearchConfig) (*OpensearchClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("opensearch config is nil")
	}
	addr := fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.Port)
	clientCfg := opensearchapi.Config{
		Client: opensearchgo.Config{
			Addresses: []string{addr},
			Username:  cfg.Username,
			Password:  cfg.Password,
			Transport: newOpensearchTransport(cfg.Insecure),
		},
	}
	c, err := opensearchapi.NewClient(clientCfg)
	if err != nil {
		return nil, err
	}
	return &OpensearchClient{Client: c}, nil
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
