//go:build linux && gui
// +build linux,gui

package platform

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	gio "github.com/diamondburned/gotk4/pkg/gio/v2"
	gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/metalgrid/drift/internal/zeroconf"
)

type promptRequest struct {
	question string
	files    []FileInfo
	peerName string
	response chan string
}

type TransferState struct {
	ID         string
	PeerName   string
	Filename   string
	Direction  string // "↑" for upload, "↓" for download
	Total      int64
	Current    int64
	Speed      float64 // bytes per second
	LastUpdate time.Time
	Status     string // "active", "complete", "failed"
}

type guiGateway struct {
	mu           *sync.Mutex
	peers        *zeroconf.Peers
	reqch        chan<- Request
	app          *gtk.Application
	window       *gtk.ApplicationWindow
	tray         *SystemTray
	transfers    map[string]*TransferState
	transferList *gtk.ListBox
	transferBox  *gtk.Box
	prompts      chan promptRequest
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

		// Send File button
		sendBtn := gtk.NewButton()
		sendBtn.SetLabel("Send File")
		peerInstance := peer.GetInstance() // Capture peer instance in closure
		sendBtn.ConnectClicked(func() {
			// Create file chooser dialog
			dialog := gtk.NewFileChooserNative(
				"Select Files to Send",
				&g.window.Window,
				gtk.FileChooserActionOpen,
				"Send",
				"Cancel",
			)
			dialog.SetSelectMultiple(true)

			dialog.Show()

			dialog.ConnectResponse(func(responseID int) {
				if responseID == int(gtk.ResponseAccept) {
					// Get selected files
					files := dialog.GetFiles()
					var filePaths []string

					for i := uint(0); i < files.NItems(); i++ {
						item := files.Item(i)
						if fileObj, ok := item.(*gio.File); ok {
							path := fileObj.Path()
							if path != "" {
								filePaths = append(filePaths, path)
							}
						}
					}

					// Send request
					if len(filePaths) > 0 {
						g.reqch <- Request{To: peerInstance, Files: filePaths}
					}
				}
			})
		})
		row.Append(sendBtn)

		// Drag-and-drop target
		peerInstanceForDrop := peer.GetInstance() // Capture before drop handler
		drop := gtk.NewDropTarget(glib.TypeString, gdk.ActionCopy)
		drop.ConnectDrop(func(drop *gtk.DropTarget, val *glib.Value, x, y float64) bool {
			// Extract file URI from dropped value
			str, ok := val.GoValue().(string)
			if !ok {
				return false
			}

			// Parse file:// URI
			if !strings.HasPrefix(str, "file://") {
				return false
			}

			// Extract path from URI
			path := strings.TrimPrefix(str, "file://")
			path = strings.TrimSpace(path)

			if path == "" {
				return false
			}

			// Send request
			g.reqch <- Request{To: peerInstanceForDrop, Files: []string{path}}
			return true
		})

		row.AddController(drop)

		listBox.Append(row)
	}

	return listBox
}

func (g *guiGateway) buildTransferList() *gtk.ListBox {
	listBox := gtk.NewListBox()
	listBox.SetSelectionMode(gtk.SelectionNone)

	g.mu.Lock()
	defer g.mu.Unlock()

	for _, transfer := range g.transfers {
		if transfer.Status != "active" {
			continue
		}

		row := gtk.NewBox(gtk.OrientationHorizontal, 10)
		row.SetMarginTop(5)
		row.SetMarginBottom(5)
		row.SetMarginStart(10)
		row.SetMarginEnd(10)

		label := gtk.NewLabel(transfer.Direction + " " + transfer.Filename)
		label.SetHExpand(true)
		label.SetXAlign(0)
		row.Append(label)

		progress := gtk.NewProgressBar()
		progress.SetSizeRequest(200, -1)
		if transfer.Total > 0 {
			fraction := float64(transfer.Current) / float64(transfer.Total)
			progress.SetFraction(fraction)
		}
		row.Append(progress)

		percentage := 0.0
		if transfer.Total > 0 {
			percentage = (float64(transfer.Current) / float64(transfer.Total)) * 100
		}
		speedText := fmt.Sprintf("%.1f MB/s - %.0f%%", transfer.Speed/1024/1024, percentage)
		speedLabel := gtk.NewLabel(speedText)
		row.Append(speedLabel)

		listBox.Append(row)
	}

	return listBox
}

func (g *guiGateway) UpdateTransfer(id string, current int64) {
	g.mu.Lock()
	transfer, exists := g.transfers[id]
	if !exists {
		g.mu.Unlock()
		return
	}

	now := time.Now()
	elapsed := now.Sub(transfer.LastUpdate).Seconds()
	if elapsed > 0 {
		delta := current - transfer.Current
		transfer.Speed = float64(delta) / elapsed
	}

	transfer.Current = current
	transfer.LastUpdate = now
	g.mu.Unlock()

	glib.IdleAdd(func() {
		newList := g.buildTransferList()
		if g.transferBox != nil {
			child := g.transferBox.FirstChild()
			if child != nil {
				g.transferBox.Remove(child)
			}
			g.transferBox.Append(newList)
		}
	})
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

		// Transfer list section
		transferLabel := gtk.NewLabel("")
		transferLabel.SetMarkup("<b>Active Transfers</b>")
		transferLabel.SetXAlign(0)
		transferLabel.SetMarginStart(10)
		transferLabel.SetMarginTop(10)
		box.Append(transferLabel)

		g.transferBox = gtk.NewBox(gtk.OrientationVertical, 0)
		g.transferList = g.buildTransferList()
		g.transferBox.Append(g.transferList)
		box.Append(g.transferBox)

		g.window.SetChild(box)

		// Initialize system tray
		tray, err := NewSystemTray(
			func() {
				// Toggle window visibility
				glib.IdleAdd(func() {
					if g.window.IsVisible() {
						g.window.Hide()
					} else {
						g.window.Show()
					}
				})
			},
			func() {
				// Quit app
				glib.IdleAdd(func() {
					g.app.Quit()
				})
			},
		)
		if err != nil {
			fmt.Printf("Failed to create system tray: %v\n", err)
			// Continue without tray
		} else {
			g.tray = tray
		}

		// Start window hidden (tray click will show it)
		g.window.SetVisible(false)

		// Handle window close button → hide to tray instead of quit
		g.window.ConnectCloseRequest(func() bool {
			g.window.Hide()
			return true // Prevent default close behavior
		})

		// Register peer change observer
		g.peers.OnChange(func() {
			glib.IdleAdd(func() {
				// Rebuild peer list on change
				newList := g.buildPeerList()
				scrolled.SetChild(newList)
			})
		})

		// Start prompt handler
		go func() {
			for req := range g.prompts {
				reqCopy := req // Capture for closure
				glib.IdleAdd(func() {
					response := g.showDialog(reqCopy)
					reqCopy.response <- response
				})
			}
		}()
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
	if g.tray != nil {
		g.tray.Close()
	}
	glib.IdleAdd(func() {
		if g.app != nil {
			g.app.Quit()
		}
	})
	close(g.reqch)
}

func (g *guiGateway) NewRequest(to, file string) error {
	g.reqch <- Request{To: to, Files: []string{file}}
	return nil
}

func (g *guiGateway) Ask(question string) string {
	responseCh := make(chan string, 1)
	g.prompts <- promptRequest{
		question: question,
		response: responseCh,
	}
	return <-responseCh
}

func (g *guiGateway) Notify(message string) {
	iconPath := "internal/platform/assets/drift-icon.svg"
	_ = SendNotification("Drift", message, iconPath)
}

func (g *guiGateway) AskBatch(peerName string, files []FileInfo) string {
	responseCh := make(chan string, 1)
	g.prompts <- promptRequest{
		peerName: peerName,
		files:    files,
		response: responseCh,
	}
	return <-responseCh
}

func (g *guiGateway) showDialog(req promptRequest) string {
	if req.peerName != "" {
		totalSize := int64(0)
		for _, f := range req.files {
			totalSize += f.Size
		}
		msg := fmt.Sprintf("Incoming transfer from %s: %d files (%s)",
			req.peerName, len(req.files), formatSize(totalSize))
		_ = SendNotification("Drift", msg, "internal/platform/assets/drift-icon.svg")
	} else if req.question != "" {
		_ = SendNotification("Drift", req.question, "internal/platform/assets/drift-icon.svg")
	}

	dialog := gtk.NewDialog()
	dialog.SetTitle("Incoming Transfer")
	dialog.SetModal(true)
	dialog.SetTransientFor(&g.window.Window)
	dialog.SetDefaultSize(400, 300)

	content := dialog.ContentArea()
	box := gtk.NewBox(gtk.OrientationVertical, 10)
	box.SetMarginTop(10)
	box.SetMarginBottom(10)
	box.SetMarginStart(10)
	box.SetMarginEnd(10)

	if req.peerName != "" {
		headerLabel := gtk.NewLabel("")
		headerLabel.SetMarkup(fmt.Sprintf("<b>Incoming files from %s</b>", req.peerName))
		box.Append(headerLabel)

		listBox := gtk.NewListBox()
		for _, file := range req.files {
			row := gtk.NewBox(gtk.OrientationHorizontal, 10)
			nameLabel := gtk.NewLabel(file.Filename)
			nameLabel.SetHExpand(true)
			nameLabel.SetXAlign(0)
			row.Append(nameLabel)
			sizeLabel := gtk.NewLabel(formatSize(file.Size))
			row.Append(sizeLabel)
			listBox.Append(row)
		}
		scrolled := gtk.NewScrolledWindow()
		scrolled.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)
		scrolled.SetVExpand(true)
		scrolled.SetChild(listBox)
		box.Append(scrolled)

		totalSize := int64(0)
		for _, f := range req.files {
			totalSize += f.Size
		}
		totalLabel := gtk.NewLabel(fmt.Sprintf("%d files, %s total", len(req.files), formatSize(totalSize)))
		box.Append(totalLabel)
	} else {
		questionLabel := gtk.NewLabel(req.question)
		questionLabel.SetWrap(true)
		box.Append(questionLabel)
	}

	progressBar := gtk.NewProgressBar()
	progressBar.SetFraction(1.0)
	box.Append(progressBar)

	countdownLabel := gtk.NewLabel("Auto-declining in 30s")
	box.Append(countdownLabel)

	content.Append(box)

	dialog.AddButton("Decline", int(gtk.ResponseReject))
	acceptBtn := dialog.AddButton("Accept", int(gtk.ResponseAccept))
	acceptBtn.AddCssClass("suggested-action")
	dialog.SetDefaultResponse(int(gtk.ResponseAccept))

	timeLeft := 30
	timerActive := true
	timeoutID := glib.TimeoutAdd(1000, func() bool {
		if !timerActive {
			return false
		}
		timeLeft--
		if timeLeft <= 0 {
			dialog.Response(int(gtk.ResponseReject))
			return false
		}
		progressBar.SetFraction(float64(timeLeft) / 30.0)
		countdownLabel.SetLabel(fmt.Sprintf("Auto-declining in %ds", timeLeft))
		return true
	})

	responseID := dialog.Run()
	timerActive = false
	glib.SourceRemove(timeoutID)
	dialog.Destroy()

	if responseID == int(gtk.ResponseAccept) {
		return "ACCEPT"
	}
	return "DECLINE"
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

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return &guiGateway{
		mu:        &sync.Mutex{},
		peers:     peers,
		reqch:     requests,
		transfers: make(map[string]*TransferState),
		prompts:   make(chan promptRequest),
	}
}
