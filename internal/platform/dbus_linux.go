//go:build linux

package platform

import (
	"fmt"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/metalgrid/drift/internal/zeroconf"
)

const (
	busName     = "com.github.metalgrid.Drift"
	objPath     = "/com/github/metalgrid/Drift"
	iface       = "com.github.metalgrid.Drift"
	SigQuestion = iface + ".Question"
	SigNotify   = iface + ".Notify"
)

type dbusService struct {
	mu            sync.Mutex
	bus           *dbus.Conn
	conversations map[string]chan string
	peers         *zeroconf.Peers
	reqch         chan<- Request
}

func newDBusService(peers *zeroconf.Peers, reqch chan<- Request) *dbusService {
	return &dbusService{
		conversations: make(map[string]chan string),
		peers:         peers,
		reqch:         reqch,
	}
}

func (d *dbusService) Start(conn *dbus.Conn) error {
	d.bus = conn

	reply, err := conn.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("failed to request name: %w", err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("service already registered")
	}

	if err := conn.Export(d, dbus.ObjectPath(objPath), iface); err != nil {
		return fmt.Errorf("failed to export object: %w", err)
	}

	methods := introspect.Methods(d)
	node := &introspect.Node{
		Name: objPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    iface,
				Methods: methods,
				Signals: []introspect.Signal{
					{
						Name: "Question",
						Args: []introspect.Arg{
							{Name: "id", Type: "s"},
							{Name: "question", Type: "s"},
						},
					},
					{
						Name: "Notify",
						Args: []introspect.Arg{
							{Name: "message", Type: "s"},
						},
					},
				},
			},
		},
	}

	if err := conn.Export(
		introspect.NewIntrospectable(node),
		dbus.ObjectPath(objPath),
		"org.freedesktop.DBus.Introspectable",
	); err != nil {
		return fmt.Errorf("failed to export introspectable: %w", err)
	}

	return nil
}

func (d *dbusService) Close() {
	if d.bus != nil {
		_, _ = d.bus.ReleaseName(busName)
		_ = d.bus.Export(nil, dbus.ObjectPath(objPath), iface)
	}
}

// RegisterConversation creates a response channel for a conversation ID.
func (d *dbusService) RegisterConversation(id string) chan string {
	ch := make(chan string, 1)
	d.mu.Lock()
	d.conversations[id] = ch
	d.mu.Unlock()
	return ch
}

// RemoveConversation removes a conversation from the map.
func (d *dbusService) RemoveConversation(id string) {
	d.mu.Lock()
	delete(d.conversations, id)
	d.mu.Unlock()
}

// EmitQuestion emits a Question signal on the bus.
func (d *dbusService) EmitQuestion(id, question string) error {
	return d.bus.Emit(dbus.ObjectPath(objPath), SigQuestion, id, question)
}

// EmitNotify emits a Notify signal on the bus.
func (d *dbusService) EmitNotify(message string) error {
	return d.bus.Emit(dbus.ObjectPath(objPath), SigNotify, message)
}

// DBus-exported methods

func (d *dbusService) Request(to, file string) *dbus.Error {
	d.reqch <- Request{To: to, Files: []string{file}}
	return nil
}

func (d *dbusService) Respond(id, answer string) *dbus.Error {
	d.mu.Lock()
	ch, ok := d.conversations[id]
	d.mu.Unlock()

	if !ok {
		return dbus.NewError(iface+".NoSuchQuestion", []any{id})
	}

	ch <- answer
	close(ch)
	return nil
}

func (d *dbusService) ListPeers() ([]string, *dbus.Error) {
	peers := d.peers.All()
	res := make([]string, len(peers))
	for i, peer := range peers {
		res[i] = peer.Instance
	}
	return res, nil
}
