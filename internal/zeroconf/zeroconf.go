package zeroconf

import (
	"context"
	"fmt"
	"maps"
	"net"
	"os"
	"os/user"
	"slices"
	"strconv"
	"strings"
	"sync"

	zc "github.com/grandcat/zeroconf"
)

const (
	serviceType   = "_dropzone._tcp"
	serviceDomain = "local."
)

type PeerInfo struct {
	*zc.ServiceEntry
}

func (pi *PeerInfo) GetInstance() string {
	instance, err := strconv.Unquote(pi.Instance)
	if err != nil {
		return pi.Instance
	}
	return instance
}

func (pi *PeerInfo) GetRecord(key string) string {
	// mDNS/Zeroconf text entries are a string in the form key=value,
	// so we add the equal sign to ensure exact lookup.
	key += "="
	for _, record := range pi.Text {
		if strings.HasPrefix(record, key) {
			return record[len(key):]
		}
	}
	return ""
}

type Peers struct {
	mu    *sync.RWMutex
	peers map[string]*PeerInfo
}

func (p *Peers) All() []*PeerInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return slices.Collect(maps.Values(p.peers))
}

func (p *Peers) GetByService(service string) *PeerInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if peer, exists := p.peers[service]; exists {
		return peer
	}

	return nil
}

func (p *Peers) GetByInstance(instance string) *PeerInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, pi := range p.peers {
		if pi.Instance == instance {
			return pi
		}
	}

	return nil
}

func (p *Peers) GetByAddr(addr net.Addr) *PeerInfo {
	remoteIP := addr.(*net.TCPAddr).IP

	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, pi := range p.peers {
		if slices.ContainsFunc(append(pi.AddrIPv4, pi.AddrIPv6...), func(ip net.IP) bool {
			return ip.Equal(remoteIP)
		}) {
			return pi
		}
	}
	return nil
}

func (p *Peers) add(pi *PeerInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers[pi.ServiceInstanceName()] = pi
}

type ZeroconfService struct {
	servicePort int
	pubkey      string
	instance    string
	peers       *Peers
	server      *zc.Server
}

func (svc *ZeroconfService) Shutdown() {
	svc.server.Shutdown()
}

func (svc *ZeroconfService) Advertise() error {
	var err error
	svc.server, err = zc.Register(
		svc.instance,
		serviceType,
		serviceDomain,
		svc.servicePort,
		[]string{
			"v=0.1",
			"pk=" + svc.pubkey,
		},
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed registering to zeroconf: %w", err)
	}

	svc.server.TTL(300)

	return err
}

func (svc *ZeroconfService) Discover(ctx context.Context) error {
	resolver, err := zc.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("failed initializing service discovery: %w", err)
	}

	entries := make(chan *zc.ServiceEntry)

	err = resolver.Browse(ctx, serviceType, serviceDomain, entries)
	if err != nil {
		return fmt.Errorf("failed discovering services: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case service := <-entries:
				// Skip ourselves
				if service.Instance == svc.instance {
					continue
				}
				svc.peers.add(&PeerInfo{
					ServiceEntry: service,
				})
			}
		}
	}()

	return nil
}

func (svc *ZeroconfService) Peers() *Peers {
	return svc.peers
}

func NewZeroconfService(port int, pubkey string) (*ZeroconfService, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed determining local machine's hostname")
	}

	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed determining local user")
	}

	username := user.Username
	if user.Name != "" {
		username = user.Name
	}

	svc := &ZeroconfService{
		servicePort: port,
		pubkey:      pubkey,
		instance:    fmt.Sprintf("%s’s Dropzone on %s", username, hostname),
		peers: &Peers{
			mu:    &sync.RWMutex{},
			peers: make(map[string]*PeerInfo),
		},
		server: nil,
	}

	return svc, nil
}
