// Code generated by DarwinKit. DO NOT EDIT.

package foundation

import (
	"github.com/progrium/darwinkit/objc"
)

// The protocol that delegates to the XPC listener use to accept or reject new connections. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsxpclistenerdelegate?language=objc
type PXPCListenerDelegate interface {
	// optional
	ListenerShouldAcceptNewConnection(listener XPCListener, newConnection XPCConnection) bool
	HasListenerShouldAcceptNewConnection() bool
}

// A delegate implementation builder for the [PXPCListenerDelegate] protocol.
type XPCListenerDelegate struct {
	_ListenerShouldAcceptNewConnection func(listener XPCListener, newConnection XPCConnection) bool
}

func (di *XPCListenerDelegate) HasListenerShouldAcceptNewConnection() bool {
	return di._ListenerShouldAcceptNewConnection != nil
}

// Accepts or rejects a new connection to the listener. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsxpclistenerdelegate/1410381-listener?language=objc
func (di *XPCListenerDelegate) SetListenerShouldAcceptNewConnection(f func(listener XPCListener, newConnection XPCConnection) bool) {
	di._ListenerShouldAcceptNewConnection = f
}

// Accepts or rejects a new connection to the listener. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsxpclistenerdelegate/1410381-listener?language=objc
func (di *XPCListenerDelegate) ListenerShouldAcceptNewConnection(listener XPCListener, newConnection XPCConnection) bool {
	return di._ListenerShouldAcceptNewConnection(listener, newConnection)
}

// ensure impl type implements protocol interface
var _ PXPCListenerDelegate = (*XPCListenerDelegateObject)(nil)

// A concrete type for the [PXPCListenerDelegate] protocol.
type XPCListenerDelegateObject struct {
	objc.Object
}

func (x_ XPCListenerDelegateObject) HasListenerShouldAcceptNewConnection() bool {
	return x_.RespondsToSelector(objc.Sel("listener:shouldAcceptNewConnection:"))
}

// Accepts or rejects a new connection to the listener. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsxpclistenerdelegate/1410381-listener?language=objc
func (x_ XPCListenerDelegateObject) ListenerShouldAcceptNewConnection(listener XPCListener, newConnection XPCConnection) bool {
	rv := objc.Call[bool](x_, objc.Sel("listener:shouldAcceptNewConnection:"), listener, newConnection)
	return rv
}
