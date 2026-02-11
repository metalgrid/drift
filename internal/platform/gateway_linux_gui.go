//go:build linux && gui
// +build linux,gui

package platform

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type guiGateway struct {
	mu     *sync.Mutex
	peers  *zeroconf.Peers
	reqch  chan<- Request
	app    *gtk.Application
	window *gtk.ApplicationWindow
}

func (g *guiGateway) Run(ctx context.Context) error {
	runtime.LockOSThread()

	g.app = gtk.NewApplication("com.github.metalgrid.drift", gio.ApplicationFlagsNone)

	g.app.ConnectActivate(func() {
		// Create window 400x500
		g.window = gtk.NewApplicationWindow(g.app)
		g.window.SetDefaultSize(400, 500)

		// Header bar with title
		header := gtk.NewHeaderBar()
		header.SetShowTitleButtons(true)
		g.window.SetTitlebar(header)
		g.window.SetTitle("Drift")

		// Empty box as placeholder
		box := gtk.NewBox(gtk.OrientationVertical, 0)
		g.window.SetChild(box)

		g.window.Show()
	})

	// Watch for context cancellation
	go func() {
		<-ctx.Done()
		glib.IdleAdd(func() {
			g.app.Quit()
		})
	}()

	g.app.Run(nil) // Blocks here
	return nil
}

func (g *guiGateway) Shutdown() {
	glib.IdleAdd(func() {
		if g.app != nil {
			g.app.Quit()
		}
	})
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
