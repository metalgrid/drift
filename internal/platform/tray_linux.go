//go:build linux

package platform

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"fmt"
	"image/png"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

//go:embed assets/icon.png
var trayIconPNG []byte

// SystemTray implements DBus StatusNotifierItem protocol for system tray integration.
type SystemTray struct {
	conn       *dbus.Conn
	props      *prop.Properties
	onActivate func()
	onQuit     func()
}

// NewSystemTray creates a new system tray icon using DBus StatusNotifierItem protocol.
// conn is a shared DBus session connection (not owned by SystemTray).
// onActivate is called on left-click, onQuit is called on right-click.
func NewSystemTray(conn *dbus.Conn, onActivate, onQuit func()) (*SystemTray, error) {
	tray := &SystemTray{
		conn:       conn,
		onActivate: onActivate,
		onQuit:     onQuit,
	}

	if err := tray.register(); err != nil {
		return nil, err
	}

	return tray, nil
}

// iconPixmap converts a PNG to SNI IconPixmap format: ARGB32 in network byte order.
func iconPixmap() []struct {
	V0 int32
	V1 int32
	V2 []byte
} {
	img, err := png.Decode(bytes.NewReader(trayIconPNG))
	if err != nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Convert to ARGB32 big-endian (network byte order)
	data := make([]byte, w*h*4)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			offset := ((y-bounds.Min.Y)*w + (x - bounds.Min.X)) * 4
			binary.BigEndian.PutUint32(data[offset:], uint32(a>>8)<<24|uint32(r>>8)<<16|uint32(g>>8)<<8|uint32(b>>8))
		}
	}

	return []struct {
		V0 int32
		V1 int32
		V2 []byte
	}{{int32(w), int32(h), data}}
}

func (t *SystemTray) register() error {
	path := dbus.ObjectPath("/StatusNotifierItem")
	sniIface := "org.kde.StatusNotifierItem"

	// Export methods (Activate, SecondaryActivate)
	if err := t.conn.Export(t, path, sniIface); err != nil {
		return fmt.Errorf("failed to export object: %w", err)
	}

	// Build properties. SNI watchers read these via org.freedesktop.DBus.Properties.GetAll.
	propsSpec := map[string]map[string]*prop.Prop{
		sniIface: {
			"Category":   {Value: "ApplicationStatus", Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"Id":         {Value: "drift", Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"Title":      {Value: "Drift", Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"Status":     {Value: "Active", Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"IconName":   {Value: "", Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"IconPixmap": {Value: iconPixmap(), Writable: false, Emit: prop.EmitFalse, Callback: nil},
			"ItemIsMenu": {Value: false, Writable: false, Emit: prop.EmitFalse, Callback: nil},
		},
	}

	var err error
	t.props, err = prop.Export(t.conn, path, propsSpec)
	if err != nil {
		return fmt.Errorf("failed to export properties: %w", err)
	}

	// Export introspection
	intro := &introspect.Node{
		Name: string(path),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:    sniIface,
				Methods: introspect.Methods(t),
				Properties: []introspect.Property{
					{Name: "Category", Type: "s", Access: "read"},
					{Name: "Id", Type: "s", Access: "read"},
					{Name: "Title", Type: "s", Access: "read"},
					{Name: "Status", Type: "s", Access: "read"},
					{Name: "IconName", Type: "s", Access: "read"},
					{Name: "IconPixmap", Type: "a(iiay)", Access: "read"},
					{Name: "ItemIsMenu", Type: "b", Access: "read"},
				},
			},
		},
	}
	if err := t.conn.Export(introspect.NewIntrospectable(intro), path, "org.freedesktop.DBus.Introspectable"); err != nil {
		return fmt.Errorf("failed to export introspection: %w", err)
	}

	// Register with StatusNotifierWatcher
	obj := t.conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	call := obj.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, "/StatusNotifierItem")
	if call.Err != nil {
		return fmt.Errorf("failed to register with watcher: %w", call.Err)
	}

	return nil
}

// DBus methods

func (t *SystemTray) Activate(x, y int32) *dbus.Error {
	if t.onActivate != nil {
		t.onActivate()
	}
	return nil
}

func (t *SystemTray) SecondaryActivate(x, y int32) *dbus.Error {
	if t.onQuit != nil {
		t.onQuit()
	}
	return nil
}

// Close unregisters the tray from DBus. Does NOT close the shared connection.
func (t *SystemTray) Close() {
	if t.conn != nil {
		_ = t.conn.Export(nil, "/StatusNotifierItem", "org.kde.StatusNotifierItem")
	}
}
