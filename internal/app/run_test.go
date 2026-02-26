package app

import (
	"encoding/hex"
	"net/netip"
	"testing"
)

func TestParsePeerPublicKeyHexRequiresExactLength(t *testing.T) {
	valid := make([]byte, 32)
	validHex := hex.EncodeToString(valid)

	if _, err := parsePeerPublicKeyHex(validHex); err != nil {
		t.Fatalf("expected valid key to parse: %v", err)
	}

	shortHex := hex.EncodeToString(make([]byte, 31))
	if _, err := parsePeerPublicKeyHex(shortHex); err == nil {
		t.Fatal("expected error for short key")
	}

	longHex := hex.EncodeToString(make([]byte, 33))
	if _, err := parsePeerPublicKeyHex(longHex); err == nil {
		t.Fatal("expected error for long key")
	}
}

func TestFirstPeerAddressRequiresAtLeastOneAddress(t *testing.T) {
	if _, ok := firstPeerAddress(nil); ok {
		t.Fatal("expected no address for nil slice")
	}

	addr := netip.MustParseAddr("192.168.1.10")
	got, ok := firstPeerAddress([]netip.Addr{addr})
	if !ok {
		t.Fatal("expected first address to exist")
	}
	if got != "192.168.1.10" {
		t.Fatalf("firstPeerAddress() = %q, want %q", got, "192.168.1.10")
	}
}
