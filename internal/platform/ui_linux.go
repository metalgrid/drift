//go:build linux

package platform

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// buildPeerPopover creates a small undecorated window for the peer list.
// It is toggled by the tray icon's left-click.
func (g *linuxGateway) buildPeerPopover() *gtk.Window {
	win := gtk.NewWindow()
	win.SetTitle("Drift - Peers")
	win.SetDefaultSize(350, 400)
	win.SetDecorated(false)
	win.SetResizable(false)

	// Escape key dismisses the popover
	esc := gtk.NewEventControllerKey()
	esc.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) bool {
		if keyval == 0xff1b { // GDK_KEY_Escape
			win.SetVisible(false)
			return true
		}
		return false
	})
	win.AddController(esc)

	g.rebuildPeerList(win)

	return win
}

// rebuildPeerList refreshes the peer list contents inside the popover window.
func (g *linuxGateway) rebuildPeerList(win *gtk.Window) {
	scrolled := gtk.NewScrolledWindow()
	scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
	scrolled.SetVExpand(true)

	listBox := gtk.NewListBox()
	listBox.SetSelectionMode(gtk.SelectionNone)

	peers := g.peers.All()
	for _, peer := range peers {
		row := gtk.NewBox(gtk.OrientationHorizontal, 10)
		row.SetMarginTop(8)
		row.SetMarginBottom(8)
		row.SetMarginStart(12)
		row.SetMarginEnd(12)

		nameLabel := gtk.NewLabel("")
		nameLabel.SetMarkup("<b>" + html.EscapeString(peer.GetInstance()) + "</b>")
		nameLabel.SetHExpand(true)
		nameLabel.SetXAlign(0)
		row.Append(nameLabel)

		osLabel := gtk.NewLabel(peer.GetRecord("os"))
		row.Append(osLabel)

		if len(peer.Addresses) > 0 {
			ipLabel := gtk.NewLabel(peer.Addresses[0].String())
			row.Append(ipLabel)
		}

		listBox.Append(row)

		// Click row to open drop window
		peerInstance := peer.Instance
		gesture := gtk.NewGestureClick()
		gesture.ConnectReleased(func(nPress int, x, y float64) {
			g.openDropWindow(peerInstance)
			win.SetVisible(false)
		})
		row.AddController(gesture)
	}

	if len(peers) == 0 {
		emptyLabel := gtk.NewLabel("No peers discovered")
		emptyLabel.SetMarginTop(20)
		emptyLabel.SetMarginBottom(20)
		listBox.Append(emptyLabel)
	}

	scrolled.SetChild(listBox)
	win.SetChild(scrolled)
}

// openDropWindow opens (or focuses) a drop target window for a specific peer.
func (g *linuxGateway) openDropWindow(peerInstance string) {
	g.mu.Lock()
	if existing, ok := g.dropWindows[peerInstance]; ok {
		g.mu.Unlock()
		existing.Present()
		return
	}
	g.mu.Unlock()

	peer := g.peers.GetByInstance(peerInstance)
	displayName := peerInstance
	if peer != nil {
		displayName = peer.GetInstance()
	}

	win := gtk.NewWindow()
	win.SetTitle("Send to " + displayName)
	win.SetDefaultSize(300, 250)

	box := gtk.NewBox(gtk.OrientationVertical, 12)
	box.SetMarginTop(16)
	box.SetMarginBottom(16)
	box.SetMarginStart(16)
	box.SetMarginEnd(16)

	// Drop zone
	dropZone := gtk.NewBox(gtk.OrientationVertical, 8)
	dropZone.SetVExpand(true)
	dropZone.SetHAlign(gtk.AlignFill)
	dropZone.SetVAlign(gtk.AlignFill)
	dropZone.AddCSSClass("drop-zone")

	// Apply dashed border style
	css := gtk.NewCSSProvider()
	css.LoadFromData(`
		.drop-zone {
			border: 2px dashed @borders;
			border-radius: 8px;
			padding: 16px;
		}
		.drop-zone-active {
			border-color: @accent_color;
			background: alpha(@accent_color, 0.1);
		}
	`)
	gtk.StyleContextAddProviderForDisplay(
		gdk.DisplayGetDefault(),
		css,
		gtk.STYLE_PROVIDER_PRIORITY_APPLICATION,
	)

	hintLabel := gtk.NewLabel("Drop files here")
	hintLabel.SetVAlign(gtk.AlignCenter)
	hintLabel.SetVExpand(true)
	dropZone.Append(hintLabel)

	// Drop target
	dropTarget := gtk.NewDropTarget(glib.TypeString, gdk.ActionCopy)
	dropTarget.ConnectDrop(func(val *glib.Value, x, y float64) bool {
		str, ok := val.GoValue().(string)
		if !ok {
			return false
		}

		paths := parseFileURIs(str)
		if len(paths) == 0 {
			return false
		}

		g.reqch <- Request{To: peerInstance, Files: paths}
		return true
	})
	dropTarget.ConnectEnter(func(x, y float64) gdk.DragAction {
		dropZone.AddCSSClass("drop-zone-active")
		return gdk.ActionCopy
	})
	dropTarget.ConnectLeave(func() {
		dropZone.RemoveCSSClass("drop-zone-active")
	})
	dropZone.AddController(dropTarget)

	box.Append(dropZone)

	// Choose Files button
	chooseBtn := gtk.NewButtonWithLabel("Choose Files...")
	chooseBtn.ConnectClicked(func() {
		dialog := gtk.NewFileChooserNative(
			"Select Files to Send",
			win,
			gtk.FileChooserActionOpen,
			"Send",
			"Cancel",
		)
		dialog.SetSelectMultiple(true)
		dialog.ConnectResponse(func(responseID int) {
			if responseID != int(gtk.ResponseAccept) {
				return
			}
			model := dialog.Files()
			var paths []string
			for i := uint(0); i < model.NItems(); i++ {
				obj := model.Item(i)
				casted := obj.Cast()
				if f, ok := casted.(*gio.File); ok {
					if p := f.Path(); p != "" {
						paths = append(paths, p)
					}
				}
			}
			if len(paths) > 0 {
				g.reqch <- Request{To: peerInstance, Files: paths}
			}
		})
		dialog.Show()
	})
	box.Append(chooseBtn)

	win.SetChild(box)

	// Track and clean up on close
	win.ConnectCloseRequest(func() bool {
		g.mu.Lock()
		delete(g.dropWindows, peerInstance)
		g.mu.Unlock()
		return false
	})

	g.mu.Lock()
	g.dropWindows[peerInstance] = win
	g.mu.Unlock()

	win.Present()
}

// showTransferDetail shows a window for an incoming transfer prompt.
func (g *linuxGateway) showTransferDetail(req promptRequest) {
	win := gtk.NewWindow()
	win.SetTitle("Incoming Transfer")
	win.SetDefaultSize(400, 350)

	box := gtk.NewBox(gtk.OrientationVertical, 10)
	box.SetMarginTop(12)
	box.SetMarginBottom(12)
	box.SetMarginStart(12)
	box.SetMarginEnd(12)

	// Header
	if req.peerName != "" {
		headerLabel := gtk.NewLabel("")
		headerLabel.SetMarkup("<b>Incoming files from " + html.EscapeString(req.peerName) + "</b>")
		box.Append(headerLabel)
	} else if req.question != "" {
		questionLabel := gtk.NewLabel(req.question)
		questionLabel.SetWrap(true)
		box.Append(questionLabel)
	}

	// File list
	if len(req.files) > 0 {
		listBox := gtk.NewListBox()
		listBox.SetSelectionMode(gtk.SelectionNone)
		totalSize := int64(0)

		for _, file := range req.files {
			row := gtk.NewBox(gtk.OrientationHorizontal, 10)
			row.SetMarginTop(4)
			row.SetMarginBottom(4)
			row.SetMarginStart(8)
			row.SetMarginEnd(8)

			nameLabel := gtk.NewLabel(file.Filename)
			nameLabel.SetHExpand(true)
			nameLabel.SetXAlign(0)
			row.Append(nameLabel)

			sizeLabel := gtk.NewLabel(formatSize(file.Size))
			row.Append(sizeLabel)

			listBox.Append(row)
			totalSize += file.Size
		}

		scrolled := gtk.NewScrolledWindow()
		scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
		scrolled.SetVExpand(true)
		scrolled.SetChild(listBox)
		box.Append(scrolled)

		totalLabel := gtk.NewLabel(fmt.Sprintf("%d files, %s total", len(req.files), formatSize(totalSize)))
		box.Append(totalLabel)
	}

	// Progress bar + countdown
	progressBar := gtk.NewProgressBar()
	progressBar.SetFraction(1.0)
	box.Append(progressBar)

	deadline := time.Now().Add(30 * time.Second)
	countdownLabel := gtk.NewLabel("Auto-declining in 30s")
	box.Append(countdownLabel)

	// Buttons
	btnBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	btnBox.SetHAlign(gtk.AlignEnd)
	btnBox.SetMarginTop(8)

	responded := false

	declineBtn := gtk.NewButtonWithLabel("Decline")
	acceptBtn := gtk.NewButtonWithLabel("Accept")
	acceptBtn.AddCSSClass("suggested-action")

	declineBtn.ConnectClicked(func() {
		if !responded {
			responded = true
			req.response <- "DECLINE"
			win.Destroy()
		}
	})

	acceptBtn.ConnectClicked(func() {
		if !responded {
			responded = true
			req.response <- "ACCEPT"
			win.Destroy()
		}
	})

	btnBox.Append(declineBtn)
	btnBox.Append(acceptBtn)
	box.Append(btnBox)

	// Countdown timer
	glib.TimeoutAdd(1000, func() bool {
		if responded {
			return false
		}
		remaining := time.Until(deadline)
		if remaining <= 0 {
			responded = true
			req.response <- "DECLINE"
			win.Destroy()
			return false
		}
		secs := int(remaining.Seconds())
		progressBar.SetFraction(float64(secs) / 30.0)
		countdownLabel.SetLabel(fmt.Sprintf("Auto-declining in %ds", secs))
		return true
	})

	win.ConnectCloseRequest(func() bool {
		if !responded {
			responded = true
			req.response <- "DECLINE"
		}
		return false
	})

	win.SetChild(box)
	win.Present()
}

// parseFileURIs extracts file paths from a newline-separated list of file:// URIs.
func parseFileURIs(s string) []string {
	var paths []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		u, err := url.Parse(line)
		if err != nil || u.Scheme != "file" {
			continue
		}
		path := u.Path
		if path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}

func formatSize(bytes int64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
		TiB = GiB * 1024
	)

	switch {
	case bytes >= TiB:
		return fmt.Sprintf("%.2f TiB", float64(bytes)/float64(TiB))
	case bytes >= GiB:
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GiB))
	case bytes >= MiB:
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MiB))
	case bytes >= KiB:
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KiB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
