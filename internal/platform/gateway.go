package platform

import (
	"context"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type Request struct {
	To   string
	File string
}

type Gateway interface {
	Run(context.Context) (<-chan Request, error)
	Shutdown()
	NewRequest(string, string) error
	Ask(string) string
}

func NewGateway(peers *zeroconf.Peers) Gateway {
	return newGateway(peers)
}
