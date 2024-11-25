//go:build linux
// +build linux

package transfer

import (
	"context"
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const (
	busName = "com.github.metalgrid.Dropzone"
	objPath = "/com/github/metalgrid/Dropzone"
	iface   = "com.github.metalgrid.Dropzone"
)

type DBusGateway struct {
	reqch chan Request
}

func NewDBusGateway() *DBusGateway {
	return &DBusGateway{}
}

func (g *DBusGateway) Start(ctx context.Context) (<-chan Request, error) {
	busConn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	reply, err := busConn.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return nil, fmt.Errorf("failed to request name: %w", err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return nil, fmt.Errorf("service already registered")
	}

	err = busConn.Export(g, objPath, iface)
	if err != nil {
		return nil, fmt.Errorf("failed to export object: %w", err)
	}

	methods := introspect.Methods(g)
	i8t := &introspect.Node{
		Name: objPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    iface,
				Methods: methods,
			},
		},
	}

	err = busConn.Export(introspect.NewIntrospectable(i8t), objPath, "org.freedesktop.DBus.Introspectable")

	if err != nil {
		return nil, fmt.Errorf("failed to export introspectable: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = busConn.Close()
	}()

	g.reqch = make(chan Request)
	return g.reqch, nil
}

func (g *DBusGateway) Shutdown() {
	close(g.reqch)
}

func (g *DBusGateway) NewRequest(to, file string) error {
	g.reqch <- Request{To: to, File: file}
	return nil
}

func (g *DBusGateway) Request(to, file string) *dbus.Error {
	err := g.NewRequest(to, file)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (g *DBusGateway) ListPeers() ([]string, *dbus.Error) {
	return nil, nil
}

func newGateway() Gateway {
	return &DBusGateway{}
}
