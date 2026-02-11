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

func (g *guiGateway) buildPeerList() *gtk.ListBox {
	listBox := gtk.NewListBox()
	listBox.SetSelectionMode(gtk.SelectionNone)

	peers := g.peers.All()
	for _, peer := range peers {
		row := gtk.NewBox(gtk.OrientationHorizontal, 10)
		row.SetMarginTop(5)
		row.SetMarginBottom(5)
		row.SetMarginStart(10)
		row.SetMarginEnd(10)

		// Peer name (bold)
		nameLabel := gtk.NewLabel(peer.GetInstance())
		nameLabel.SetMarkup("<b>" + peer.GetInstance() + "</b>")
		nameLabel.SetHExpand(true)
		nameLabel.SetXAlign(0)
		row.Append(nameLabel)

		// OS badge
		osLabel := gtk.NewLabel(peer.GetRecord("os"))
		row.Append(osLabel)

		// IP address (first address)
		if len(peer.Addresses) > 0 {
			ipLabel := gtk.NewLabel(peer.Addresses[0].String())
			row.Append(ipLabel)
		}

		listBox.Append(row)
	}

	return listBox
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

		// Create scrolled window for peer list
		scrolled := gtk.NewScrolledWindow()
		scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
		scrolled.SetVExpand(true)

		peerList := g.buildPeerList()
		scrolled.SetChild(peerList)

		// Main container
		box := gtk.NewBox(gtk.OrientationVertical, 0)
		box.Append(scrolled)

		g.window.SetChild(box)
		g.window.Show()

		// Register peer change observer
		g.peers.OnChange(func() {
			glib.IdleAdd(func() {
				// Rebuild peer list on change
				newList := g.buildPeerList()
				scrolled.SetChild(newList)
			})
		})
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
