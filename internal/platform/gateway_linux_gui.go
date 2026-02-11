//go:build linux && gui
// +build linux,gui

package platform

import (
	"context"
	"fmt"
	"sync"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type guiGateway struct {
	mu    *sync.Mutex
	peers *zeroconf.Peers
	reqch chan<- Request
}

func (g *guiGateway) Run(ctx context.Context) error {
	fmt.Println("GUI gateway: Run() not implemented")
	<-ctx.Done()
	return nil
}

func (g *guiGateway) Shutdown() {
	fmt.Println("GUI gateway: Shutdown() not implemented")
	close(g.reqch)
}

func (g *guiGateway) NewRequest(to, file string) error {
	fmt.Println("GUI gateway: NewRequest() not implemented")
	g.reqch <- Request{To: to, Files: []string{file}}
	return nil
}

func (g *guiGateway) Ask(question string) string {
	fmt.Println("GUI gateway: Ask() not implemented")
	return "DECLINE"
}

func (g *guiGateway) Notify(message string) {
	iconPath := "internal/platform/assets/drift-icon.svg"
	_ = SendNotification("Drift", message, iconPath)
}

func (g *guiGateway) AskBatch(peerName string, files []FileInfo) string {
	fmt.Println("GUI gateway: AskBatch() not implemented")
	return "DECLINE"
}

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return &guiGateway{
		&sync.Mutex{},
		peers,
		requests,
	}
}
