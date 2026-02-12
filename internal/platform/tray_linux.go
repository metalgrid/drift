//go:build linux && gui
// +build linux,gui

package platform

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

// SystemTray implements DBus StatusNotifierItem protocol for system tray integration.
type SystemTray struct {
	conn       *dbus.Conn
	onActivate func()
	onQuit     func()
}

// NewSystemTray creates a new system tray icon using DBus StatusNotifierItem protocol.
// onActivate is called on left-click, onQuit is called on right-click.
func NewSystemTray(onActivate, onQuit func()) (*SystemTray, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	tray := &SystemTray{
		conn:       conn,
		onActivate: onActivate,
		onQuit:     onQuit,
	}

	// Register StatusNotifierItem
	if err := tray.register(); err != nil {
		conn.Close()
		return nil, err
	}

	return tray, nil
}

func (t *SystemTray) register() error {
	// Export object on DBus
	path := dbus.ObjectPath("/StatusNotifierItem")

	// Export methods and properties
	if err := t.conn.Export(t, path, "org.kde.StatusNotifierItem"); err != nil {
		return fmt.Errorf("failed to export object: %w", err)
	}

	// Export introspection data
	intro := &introspect.Node{
		Name: string(path),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    "org.kde.StatusNotifierItem",
				Methods: introspect.Methods(t),
				Properties: []introspect.Property{
					{Name: "Category", Type: "s", Access: "read"},
					{Name: "Id", Type: "s", Access: "read"},
					{Name: "Title", Type: "s", Access: "read"},
					{Name: "IconName", Type: "s", Access: "read"},
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

// DBus property getters (required by SNI spec)

func (t *SystemTray) Category() (string, *dbus.Error) {
	return "ApplicationStatus", nil
}

func (t *SystemTray) Id() (string, *dbus.Error) {
	return "drift", nil
}

func (t *SystemTray) Title() (string, *dbus.Error) {
	return "Drift", nil
}

func (t *SystemTray) IconName() (string, *dbus.Error) {
	return "network-transmit-receive", nil
}

func (t *SystemTray) ItemIsMenu() (bool, *dbus.Error) {
	return false, nil
}

// DBus methods (signals)

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

// Close closes the DBus connection.
func (t *SystemTray) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}
