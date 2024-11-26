package transfer

import (
	"context"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type Request struct {
	To   string
	File string
}

type Gateway interface {
	Start(context.Context) (<-chan Request, error)
	Shutdown()
	NewRequest(string, string) error
}

func NewGateway(peers *zeroconf.Peers) Gateway {
	return newGateway(peers)
}
