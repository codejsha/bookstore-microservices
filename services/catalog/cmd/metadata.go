package cmd

import (
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/config"
)

var (
	Name    string
	Version string
	Commit  string
)

var meta = &config.Metadata{
	Name:    Name,
	Version: Version,
	Commit:  Commit,
}

func metadata() *config.Metadata {
	return meta
}
