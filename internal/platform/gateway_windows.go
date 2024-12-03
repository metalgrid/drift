//go:build windows
// +build windows

package platform

//go:generate go run github.com/akavel/rsrc@latest -manifest ../../windows/drift.manifest -ico ../../windows/drift.ico

import (
	"context"
	"fmt"

	"github.com/tailscale/walk"
	// . "github.com/tailscale/walk/declarative"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type Win32Gateway struct {
	app      *walk.Application
	trayIcon *walk.NotifyIcon
	peers    *zeroconf.Peers
	reqch    chan<- Request
}

func (g *Win32Gateway) Run(ctx context.Context) error {
	app, err := walk.InitApp()
	if err != nil {
		return fmt.Errorf("platform initialization failed: %w", err)
	}

	// icon, err := walk.NewIconFromResource("icon")
	// if err != nil {
	// 	return fmt.Errorf("failed setting app icon: %w", err)
	// }

	tray, err := walk.NewNotifyIcon()
	if err != nil {
		return fmt.Errorf("failed creating notification icon: %w", err)
	}

	// err = tray.SetIcon(icon)
	if err != nil {
		return fmt.Errorf("failed setting app icon: %w", err)
	}

	tray.SetToolTip("Drift - pain free, secure file transfer")
	tray.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		for _, peer := range g.peers.All() {
			action := walk.NewAction()
			action.SetText(peer.Instance)
			action.Triggered().Attach(func() {
				picker := walk.FileDialog{
					Title: "Send file",
				}

				if ok, err := picker.ShowOpen(nil); ok && err == nil {
					g.NewRequest(peer.Instance, picker.FilePath)
				}
			})
			tray.ContextMenu().Actions().Add(action)
		}
	})

	tray.SetVisible(true) //when would this not work, windows pls

	app.Run()
	return nil
}

func (g *Win32Gateway) Shutdown() {}
func (g *Win32Gateway) NewRequest(peer, file string) error {
	g.reqch <- Request{
		To:   peer,
		File: file,
	}
	return nil
}
func (g *Win32Gateway) Ask(string) string { return "" }
func (g *Win32Gateway) Notify(msg string) {
	fmt.Println(msg)
}

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	_ = peers
	return &Win32Gateway{
		peers: peers,
		reqch: requests,
	}
}
