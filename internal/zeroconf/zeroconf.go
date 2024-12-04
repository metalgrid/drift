package zeroconf

import (
	"context"
	"fmt"
	"maps"
	"net"
	"net/netip"
	"os"
	"os/user"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"

	zc "github.com/betamos/zeroconf"
)

const (
	serviceType   = "_drift._tcp"
	serviceDomain = "local."
)

type PeerInfo struct {
	Service   string
	Instance  string
	Domain    string
	Port      int
	Records   []string
	Addresses []netip.Addr
}

func (pi *PeerInfo) String() string {
	return fmt.Sprintf("%s.%s.%s", pi.Instance, pi.Service, pi.Domain)
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
	for _, record := range pi.Records {
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
	remoteIP := netip.MustParseAddrPort(addr.String())

	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, pi := range p.peers {
		if slices.ContainsFunc(pi.Addresses, func(ip netip.Addr) bool {
			return ip == remoteIP.Addr()
		}) {
			return pi
		}
	}
	return nil
}

func (p *Peers) add(pi *PeerInfo) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers[pi.String()] = pi
}

func (p *Peers) remove(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.peers, key)
}

type ZeroconfService struct {
	servicePort int
	pubkey      string
	instance    string
	peers       *Peers
	client      *zc.Client
}

func (svc *ZeroconfService) Shutdown() {
	_ = svc.client.Close()
}

func (svc *ZeroconfService) Start(ctx context.Context) error {
	kind := zc.NewType(serviceType)
	service := zc.NewService(kind, svc.instance, uint16(svc.servicePort))
	service.Text = []string{
		"v=0.1",
		"pk=" + svc.pubkey,
		"os=" + runtime.GOOS,
	}

	svc.client.Publish(service)

	svc.client.Browse(func(e zc.Event) {
		fmt.Println(e.String())
		switch e.Op {
		case zc.OpAdded, zc.OpUpdated:
			svc.peers.add(&PeerInfo{
				Service:   e.Type.Name,
				Domain:    e.Type.Domain,
				Instance:  e.Name,
				Records:   e.Text,
				Addresses: e.Addrs,
			})
		case zc.OpRemoved:
			svc.peers.remove(e.Service.String())
		}
	}, kind)

	_, err := svc.client.Open()
	return err
}

func (svc *ZeroconfService) Peers() *Peers {
	return svc.peers
}

type ZeroconfOptions struct {
	Identity string
}

func NewZeroconfService(port int, pubkey string, options *ZeroconfOptions) (*ZeroconfService, error) {

	var identity string

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

	identity = fmt.Sprintf("%s’s %s", username, hostname)
	if options != nil {
		if options.Identity != "" {
			identity = options.Identity
		}
	}

	client := zc.New()

	svc := &ZeroconfService{
		servicePort: port,
		pubkey:      pubkey,
		instance:    identity,
		peers: &Peers{
			mu:    &sync.RWMutex{},
			peers: make(map[string]*PeerInfo),
		},
		client: client,
	}

	return svc, nil
}
