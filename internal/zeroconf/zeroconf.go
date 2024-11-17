package zeroconf

import (
	"context"
	"fmt"
	"maps"
	"net"
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

var me string = "iso's\\ Dropzone\\ on\\ archbtw"

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

func Advertise(port int, username, hostname, pubkey string) (*zc.Server, error) {
	server, err := zc.Register(
		fmt.Sprintf("%s's Dropzone on %s", username, hostname),
		serviceType,
		serviceDomain,
		port,
		[]string{
			"v=0.1",
			"u=" + username,
			"pk=" + pubkey,
		},
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed registering to zeroconf: %w", err)
	}

	server.TTL(300)

	return server, nil
}

func Discover(ctx context.Context) (*Peers, error) {
	resolver, err := zc.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("failed initializing service discovery: %w", err)
	}

	entries := make(chan *zc.ServiceEntry)

	err = resolver.Browse(ctx, serviceType, serviceDomain, entries)
	if err != nil {
		return nil, fmt.Errorf("failed discovering services: %w", err)
	}

	peers := &Peers{
		mu:    &sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case service := <-entries:
				if service.Instance != me {
					peers.add(&PeerInfo{
						ServiceEntry: service,
					})
				}
				fmt.Printf("Discovered service:\n%v", service)
			}
		}
	}()

	return peers, nil
}
