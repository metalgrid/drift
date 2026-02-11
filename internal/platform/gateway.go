package platform

import (
	"context"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type Request struct {
	To    string
	Files []string
}

type FileInfo struct {
	Filename string
	Size     int64
}

type Gateway interface {
	// Run starts the platform-dependent logic. Blocks execution until the context is cancelled.
	// This method may call FFI or UI event loops, so it's guaranteed to run on the main thread.
	// Using `runtime.LockOSThread()` in implementations is encouraged.
	Run(context.Context) error
	Shutdown()
	NewRequest(string, string) error
	Ask(string) string
	Notify(string)
}

type BatchGateway interface {
	Gateway
	AskBatch(peerName string, files []FileInfo) string
}

func NewGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return newGateway(peers, requests)
}
