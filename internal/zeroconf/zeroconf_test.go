package zeroconf

import (
	"net/netip"
	"sync"
	"testing"
)

func TestPeersOnChangeCalledOnAdd(t *testing.T) {
	peers := &Peers{
		mu:    &sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
	}

	called := false
	peers.OnChange(func() {
		called = true
	})

	pi := &PeerInfo{
		Service:   "_drift._tcp",
		Instance:  "test-peer",
		Domain:    "local.",
		Port:      38473,
		Records:   []string{"v=0.1"},
		Addresses: []netip.Addr{netip.MustParseAddr("192.168.1.100")},
	}

	peers.add(pi)

	if !called {
		t.Error("Expected observer to be called on add, but it was not")
	}
}

func TestPeersOnChangeCalledOnRemove(t *testing.T) {
	peers := &Peers{
		mu:    &sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
	}

	addCalled := false
	removeCalled := false
	callCount := 0

	peers.OnChange(func() {
		callCount++
		if callCount == 1 {
			addCalled = true
		} else if callCount == 2 {
			removeCalled = true
		}
	})

	pi := &PeerInfo{
		Service:  "_drift._tcp",
		Instance: "test-peer",
		Domain:   "local.",
		Port:     38473,
	}

	peers.add(pi)
	peers.remove(pi.String())

	if !addCalled {
		t.Error("Expected observer to be called on add")
	}
	if !removeCalled {
		t.Error("Expected observer to be called on remove")
	}
	if callCount != 2 {
		t.Errorf("Expected observer to be called exactly 2 times, got %d", callCount)
	}
}

func TestPeersMultipleObservers(t *testing.T) {
	peers := &Peers{
		mu:    &sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
	}

	observer1Called := false
	observer2Called := false

	peers.OnChange(func() {
		observer1Called = true
	})

	peers.OnChange(func() {
		observer2Called = true
	})

	pi := &PeerInfo{
		Service:  "_drift._tcp",
		Instance: "test-peer",
		Domain:   "local.",
		Port:     38473,
	}

	peers.add(pi)

	if !observer1Called {
		t.Error("Expected first observer to be called")
	}
	if !observer2Called {
		t.Error("Expected second observer to be called")
	}
}

func TestPeersAllReturnsCopy(t *testing.T) {
	peers := &Peers{
		mu:    &sync.RWMutex{},
		peers: make(map[string]*PeerInfo),
	}

	pi1 := &PeerInfo{
		Service:  "_drift._tcp",
		Instance: "peer-1",
		Domain:   "local.",
		Port:     38473,
	}

	pi2 := &PeerInfo{
		Service:  "_drift._tcp",
		Instance: "peer-2",
		Domain:   "local.",
		Port:     38474,
	}

	peers.add(pi1)
	peers.add(pi2)

	snapshot := peers.All()

	if len(snapshot) != 2 {
		t.Errorf("Expected 2 peers in snapshot, got %d", len(snapshot))
	}

	// Verify it's a snapshot by adding another peer and checking original slice unchanged
	pi3 := &PeerInfo{
		Service:  "_drift._tcp",
		Instance: "peer-3",
		Domain:   "local.",
		Port:     38475,
	}
	peers.add(pi3)

	if len(snapshot) != 2 {
		t.Error("Expected snapshot to remain unchanged after adding new peer")
	}

	newSnapshot := peers.All()
	if len(newSnapshot) != 3 {
		t.Errorf("Expected 3 peers in new snapshot, got %d", len(newSnapshot))
	}
}
