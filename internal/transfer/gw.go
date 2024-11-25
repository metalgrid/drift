package transfer

import "context"

type Request struct {
	To   string
	File string
}

type Gateway interface {
	Start(context.Context) (<-chan Request, error)
	Shutdown()
	NewRequest(string, string) error
}

func NewGateway() Gateway {
	return newGateway()
}
