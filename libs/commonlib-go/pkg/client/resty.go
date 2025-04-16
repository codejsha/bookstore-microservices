package client

import (
	"github.com/go-resty/resty/v2"
)

func NewRestyClient() *RestyClient {
	client := resty.New()
	return &RestyClient{
		Client: client,
	}
}

type RestyClient struct {
	Client *resty.Client
}
