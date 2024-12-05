//go:build darwin
// +build darwin

package platform

import (
	"context"
	"strconv"

	"github.com/progrium/darwinkit/macos"
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type macGateway struct {
}

func (m *macGateway) Run(ctx context.Context) error {
	macos.RunApp(launched)
	return nil
}

func (m *macGateway) Shutdown() {
}

func (m *macGateway) NewRequest(peer string, file string) error {
	return nil
}

func (m *macGateway) Ask(question string) string {
	return ""
}

func (m *macGateway) Notify(msg string) {
}

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return &macGateway{}
}

func launched(app appkit.Application, delegate *appkit.ApplicationDelegate) {
	w := appkit.NewWindowWithSize(720, 440)
	objc.Retain(&w)
	w.SetTitle("Decoder")

	tabView := appkit.NewTabView()
	tabView.SetTranslatesAutoresizingMaskIntoConstraints(false)

	// add tabs
	tabView.AddTabViewItem(createNewView(1))
	tabView.AddTabViewItem(createNewView(2))

	w.SetContentView(tabView)
	w.Center()
	w.MakeKeyAndOrderFront(nil)

	delegate.SetApplicationWillFinishLaunching(func(foundation.Notification) {
		w.SetFrameAutosaveName("tab-test")
	})
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return true
	})
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
	app.ActivateIgnoringOtherApps(true)

}

func createNewView(idx int) appkit.ITabViewItem {
	sv := appkit.NewStackView()
	sv.SetTranslatesAutoresizingMaskIntoConstraints(true)
	sv.AddViewInGravity(appkit.NewButtonWithTitle("button"), appkit.StackViewGravityTop)
	sv.AddViewInGravity(appkit.NewTextField(), appkit.StackViewGravityTop)
	ti := appkit.NewTabViewItem()
	ti.SetLabel("Tab" + strconv.Itoa(idx))
	ti.SetView(sv)
	return ti
}
