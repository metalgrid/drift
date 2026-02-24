//go:build linux

package platform

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/godbus/dbus/v5"

	"github.com/metalgrid/drift/internal/zeroconf"
)

const responseTimeout = 30 * time.Second

type promptRequest struct {
	question string
	files    []FileInfo
	peerName string
	response chan string
}

type linuxGateway struct {
	mu      sync.Mutex
	peers   *zeroconf.Peers
	reqch   chan<- Request
	app     *gtk.Application
	busConn *dbus.Conn
	dbus    *dbusService
	tray    *SystemTray
	notif   *notifier
	prompts chan promptRequest

	peerWindow  *gtk.Window
	dropWindows map[string]*gtk.Window
}

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return &linuxGateway{
		peers:       peers,
		reqch:       requests,
		prompts:     make(chan promptRequest),
		dropWindows: make(map[string]*gtk.Window),
	}
}

func (g *linuxGateway) Run(ctx context.Context) error {
	runtime.LockOSThread()

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}
	g.busConn = conn

	g.dbus = newDBusService(g.peers, g.reqch)
	g.notif = newNotifier()

	g.app = gtk.NewApplication("com.github.metalgrid.drift", gio.ApplicationFlagsNone)

	g.app.ConnectActivate(func() {
		// Keep the GTK application alive even without windows (tray-only app)
		g.app.Hold()

		// Start DBus service
		if err := g.dbus.Start(conn); err != nil {
			fmt.Printf("Failed to start DBus service: %v\n", err)
		}

		// Start notification listener
		if err := g.notif.Start(conn); err != nil {
			fmt.Printf("Failed to start notifier: %v\n", err)
		}

		// Start system tray
		tray, err := NewSystemTray(
			conn,
			func() {
				// Left-click: toggle peer popover
				glib.IdleAdd(func() {
					if g.peerWindow == nil {
						g.peerWindow = g.buildPeerPopover()
					}
					g.peerWindow.SetVisible(!g.peerWindow.IsVisible())
				})
			},
			func() {
				// Right-click: quit
				glib.IdleAdd(func() {
					g.app.Quit()
				})
			},
		)
		if err != nil {
			fmt.Printf("Failed to create system tray: %v\n", err)
		} else {
			g.tray = tray
		}

		// Observe peer changes -> rebuild popover
		g.peers.OnChange(func() {
			glib.IdleAdd(func() {
				if g.peerWindow != nil {
					g.rebuildPeerList(g.peerWindow)
				}
			})
		})

		// Prompt handler goroutine: reads from channel, shows UI on GTK thread
		go func() {
			for req := range g.prompts {
				reqCopy := req
				glib.IdleAdd(func() {
					g.showTransferDetail(reqCopy)
				})
			}
		}()
	})

	// Watch context cancellation
	go func() {
		<-ctx.Done()
		glib.IdleAdd(func() {
			g.app.Quit()
		})
	}()

	g.app.Run(nil)
	return nil
}

func (g *linuxGateway) Shutdown() {
	if g.tray != nil {
		g.tray.Close()
	}
	if g.notif != nil {
		g.notif.Close()
	}
	if g.dbus != nil {
		g.dbus.Close()
	}
	if g.busConn != nil {
		_ = g.busConn.Close()
	}
	if g.app != nil {
		g.app.Quit()
	}
	close(g.reqch)
}

func (g *linuxGateway) NewRequest(to, file string) error {
	g.reqch <- Request{To: to, Files: []string{file}}
	return nil
}

func (g *linuxGateway) Ask(question string) string {
	id := g.generateID()

	ch := g.dbus.RegisterConversation(id)
	defer g.dbus.RemoveConversation(id)

	// Emit DBus signal so CLI tools can respond
	if err := g.dbus.EmitQuestion(id, question); err != nil {
		fmt.Printf("failed to emit question signal: %v\n", err)
	}

	// Send desktop notification with action
	g.notif.Send("Drift", question, notifyIconPath(),
		map[string]string{"show": "Show Details"},
		func(actionKey string) {
			glib.IdleAdd(func() {
				// Raise/focus the detail window if it exists
				// The prompt is already being shown via the prompts channel
			})
		},
	)

	// Show detail window on GTK thread
	responseCh := make(chan string, 1)
	g.prompts <- promptRequest{
		question: question,
		response: responseCh,
	}

	// Wait for response from either GTK window or DBus Respond call
	select {
	case answer := <-responseCh:
		return answer
	case answer := <-ch:
		return answer
	case <-time.After(responseTimeout):
		return "DECLINE"
	}
}

func (g *linuxGateway) AskBatch(peerName string, files []FileInfo) string {
	id := g.generateID()

	ch := g.dbus.RegisterConversation(id)
	defer g.dbus.RemoveConversation(id)

	// Build question text for DBus signal
	totalSize := int64(0)
	for _, f := range files {
		totalSize += f.Size
	}
	question := fmt.Sprintf("Incoming transfer from %s: %d files (%s)",
		peerName, len(files), formatSize(totalSize))

	if err := g.dbus.EmitQuestion(id, question); err != nil {
		fmt.Printf("failed to emit question signal: %v\n", err)
	}

	// Send desktop notification with action
	g.notif.Send("Drift", question, notifyIconPath(),
		map[string]string{"show": "Show Details"},
		func(actionKey string) {
			glib.IdleAdd(func() {
				// Window is already shown via prompts channel
			})
		},
	)

	// Show detail window on GTK thread
	responseCh := make(chan string, 1)
	g.prompts <- promptRequest{
		peerName: peerName,
		files:    files,
		response: responseCh,
	}

	select {
	case answer := <-responseCh:
		return answer
	case answer := <-ch:
		return answer
	case <-time.After(responseTimeout):
		return "DECLINE"
	}
}

func (g *linuxGateway) Notify(message string) {
	// Emit DBus signal
	if g.dbus != nil {
		_ = g.dbus.EmitNotify(message)
	}

	// Send desktop notification
	if g.notif != nil {
		g.notif.Send("Drift", message, notifyIconPath(), nil, nil)
	}
}

func (g *linuxGateway) generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
