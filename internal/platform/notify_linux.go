//go:build linux

package platform

import (
	_ "embed"
	"os"
	"path/filepath"
	"sync"

	"github.com/godbus/dbus/v5"
)

//go:embed assets/icon.png
var notifyIconPNG []byte

// cachedIconPath stores the temp file path for the notification icon.
var cachedIconPath string

// notifyIconPath returns a file path to the embedded icon for use with notifications.
func notifyIconPath() string {
	if cachedIconPath != "" {
		return cachedIconPath
	}
	dir := filepath.Join(os.TempDir(), "drift")
	_ = os.MkdirAll(dir, 0o700)
	p := filepath.Join(dir, "icon.png")
	if err := os.WriteFile(p, notifyIconPNG, 0o644); err != nil {
		return ""
	}
	cachedIconPath = p
	return p
}

type notifier struct {
	bus     *dbus.Conn
	mu      sync.Mutex
	pending map[uint32]func(string)
	signal  chan *dbus.Signal
}

func newNotifier() *notifier {
	return &notifier{
		pending: make(map[uint32]func(string)),
		signal:  make(chan *dbus.Signal, 16),
	}
}

func (n *notifier) Start(conn *dbus.Conn) error {
	n.bus = conn

	if err := conn.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.Notifications"),
		dbus.WithMatchMember("ActionInvoked"),
	); err != nil {
		return err
	}

	conn.Signal(n.signal)

	go n.listenActions()

	return nil
}

func (n *notifier) listenActions() {
	for sig := range n.signal {
		if sig.Name != "org.freedesktop.Notifications.ActionInvoked" {
			continue
		}
		if len(sig.Body) < 2 {
			continue
		}
		id, ok1 := sig.Body[0].(uint32)
		actionKey, ok2 := sig.Body[1].(string)
		if !ok1 || !ok2 {
			continue
		}

		n.mu.Lock()
		cb, exists := n.pending[id]
		if exists {
			delete(n.pending, id)
		}
		n.mu.Unlock()

		if exists && cb != nil {
			cb(actionKey)
		}
	}
}

// Send sends a desktop notification with optional actions.
// onAction is called with the action key when the user clicks an action.
func (n *notifier) Send(summary, body, icon string, actions map[string]string, onAction func(string)) (uint32, error) {
	// Build actions list: [key, label, key, label, ...]
	var actionList []string
	for key, label := range actions {
		actionList = append(actionList, key, label)
	}

	obj := n.bus.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Drift",                   // app_name
		uint32(0),                 // replaces_id
		icon,                      // app_icon
		summary,                   // summary
		body,                      // body
		actionList,                // actions
		map[string]dbus.Variant{}, // hints
		int32(10000),              // expire_timeout
	)
	if call.Err != nil {
		return 0, call.Err
	}

	var id uint32
	if err := call.Store(&id); err != nil {
		return 0, err
	}

	if onAction != nil && len(actions) > 0 {
		n.mu.Lock()
		n.pending[id] = onAction
		n.mu.Unlock()
	}

	return id, nil
}

func (n *notifier) Close() {
	if n.bus != nil {
		n.bus.RemoveSignal(n.signal)
		_ = n.bus.RemoveMatchSignal(
			dbus.WithMatchInterface("org.freedesktop.Notifications"),
			dbus.WithMatchMember("ActionInvoked"),
		)
	}
	close(n.signal)
}
