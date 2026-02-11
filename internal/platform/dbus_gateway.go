//go:build linux
// +build linux

package platform

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/metalgrid/drift/internal/zeroconf"
)

const (
	busName          = "com.github.metalgrid.Drift"
	objPath          = "/com/github/metalgrid/Drift"
	iface            = "com.github.metalgrid.Drift"
	SigQuestion      = iface + ".Question"
	SigNotify        = iface + ".Notify"
	ErrNoSuchQestion = iface + ".NoSuchQuestion"
	ResponseTimeout  = 30 * time.Second
)

type DBusGateway struct {
	mu            *sync.Mutex
	conversations map[string]chan string
	bus           *dbus.Conn
	peers         *zeroconf.Peers
	reqch         chan<- Request
}

func (g *DBusGateway) Run(ctx context.Context) error {
	busConn, err := dbus.ConnectSessionBus()
	if err != nil {
		return fmt.Errorf("failed to connect to session bus: %w", err)
	}

	g.bus = busConn

	reply, err := busConn.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf("failed to request name: %w", err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf("service already registered")
	}

	err = busConn.Export(g, objPath, iface)
	if err != nil {
		return fmt.Errorf("failed to export object: %w", err)
	}

	methods := introspect.Methods(g)
	i8t := &introspect.Node{
		Name: objPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    iface,
				Methods: methods,
				Signals: []introspect.Signal{
					{
						Name: SigQuestion,
						Args: []introspect.Arg{
							{Name: "id", Type: "s"},
							{Name: "question", Type: "s"},
						},
					},
					{
						Name: SigNotify,
						Args: []introspect.Arg{
							{Name: "message", Type: "s"},
						},
					},
				},
			},
		},
	}

	err = busConn.Export(introspect.NewIntrospectable(i8t), objPath, "org.freedesktop.DBus.Introspectable")

	if err != nil {
		return fmt.Errorf("failed to export introspectable: %w", err)
	}

	<-ctx.Done()
	_ = busConn.Close()

	return nil
}

func (g *DBusGateway) Shutdown() {
	close(g.reqch)
}

func (g *DBusGateway) NewRequest(to, file string) error {
	g.reqch <- Request{To: to, Files: []string{file}}
	return nil
}

func (g *DBusGateway) Request(to, file string) *dbus.Error {
	err := g.NewRequest(to, file)
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func (g *DBusGateway) Ask(question string) string {
	idBytes := make([]byte, 32)
	rand.Read(idBytes)
	id := hex.EncodeToString(idBytes)
	g.mu.Lock()
	g.conversations[id] = make(chan string, 1)
	g.mu.Unlock()

	defer func() {
		g.mu.Lock()
		delete(g.conversations, id)
		g.mu.Unlock()
	}()

	err := g.bus.Emit(objPath, SigQuestion, id, question)
	if err != nil {
		fmt.Println("failed to emit question:", err)
		return "DECLINE"
	}

	select {
	case <-time.After(ResponseTimeout):
		return "DECLINE"
	case answer := <-g.conversations[id]:
		return answer
	}
}

func (g *DBusGateway) Notify(message string) {
	_ = g.bus.Emit(objPath, SigNotify, message)
}

func (g *DBusGateway) Respond(id, answer string) *dbus.Error {
	g.mu.Lock()
	conversation, ok := g.conversations[id]
	g.mu.Unlock()

	if !ok {
		return dbus.NewError(ErrNoSuchQestion, []any{id})
	}

	conversation <- answer
	close(conversation)
	return nil
}

func (g *DBusGateway) ListPeers() ([]string, *dbus.Error) {
	peers := g.peers.All()
	res := make([]string, len(peers))
	for idx, peer := range peers {
		res[idx] = peer.Instance
	}
	return res, nil
}

func newGateway(peers *zeroconf.Peers, requests chan<- Request) Gateway {
	return &DBusGateway{
		&sync.Mutex{},
		make(map[string]chan string),
		nil,
		peers,
		requests,
	}
}
