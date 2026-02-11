//go:build linux && gui
// +build linux,gui

package platform

import (
	"github.com/godbus/dbus/v5"
)

// SendNotification sends a desktop notification via freedesktop DBus.
// Uses org.freedesktop.Notifications.Notify method.
func SendNotification(summary, body, icon string) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	call := obj.Call("org.freedesktop.Notifications.Notify", 0,
		"Drift",                   // app_name
		uint32(0),                 // replaces_id
		icon,                      // app_icon
		summary,                   // summary
		body,                      // body
		[]string{},                // actions
		map[string]dbus.Variant{}, // hints
		int32(5000),               // expire_timeout (milliseconds)
	)

	return call.Err
}
