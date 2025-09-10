package support

import (
	"crypto/tls"
	"fmt"
	"net/http"

	opensearchgo "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/codejsha/bookstore-microservices/catalog/internal/config"
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
		},
	}
	if cfg.Insecure {
		clientCfg.Client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	c, err := opensearchapi.NewClient(clientCfg)
	if err != nil {
		return nil, err
	}
	return &OpensearchClient{Client: c}, nil
}
