package config

import (
	"fmt"
)

type Metadata struct {
	Name    string
	Version string
	Commit  string
}

func (m *Metadata) String() string {
	return fmt.Sprintf("Name: %s, Version: %s, Commit: %s", m.Name, m.Version, m.Commit)
}
