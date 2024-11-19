// Code generated by DarwinKit. DO NOT EDIT.

package foundation

import (
	"github.com/progrium/darwinkit/objc"
)

// An interface that delegates of a stream instance use to handle events on the stream. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsstreamdelegate?language=objc
type PStreamDelegate interface {
	// optional
	StreamHandleEvent(aStream Stream, eventCode StreamEvent)
	HasStreamHandleEvent() bool
}

// A delegate implementation builder for the [PStreamDelegate] protocol.
type StreamDelegate struct {
	_StreamHandleEvent func(aStream Stream, eventCode StreamEvent)
}

func (di *StreamDelegate) HasStreamHandleEvent() bool {
	return di._StreamHandleEvent != nil
}

// The delegate receives this message when a given event has occurred on a given stream. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsstreamdelegate/1410079-stream?language=objc
func (di *StreamDelegate) SetStreamHandleEvent(f func(aStream Stream, eventCode StreamEvent)) {
	di._StreamHandleEvent = f
}

// The delegate receives this message when a given event has occurred on a given stream. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsstreamdelegate/1410079-stream?language=objc
func (di *StreamDelegate) StreamHandleEvent(aStream Stream, eventCode StreamEvent) {
	di._StreamHandleEvent(aStream, eventCode)
}

// ensure impl type implements protocol interface
var _ PStreamDelegate = (*StreamDelegateObject)(nil)

// A concrete type for the [PStreamDelegate] protocol.
type StreamDelegateObject struct {
	objc.Object
}

func (s_ StreamDelegateObject) HasStreamHandleEvent() bool {
	return s_.RespondsToSelector(objc.Sel("stream:handleEvent:"))
}

// The delegate receives this message when a given event has occurred on a given stream. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsstreamdelegate/1410079-stream?language=objc
func (s_ StreamDelegateObject) StreamHandleEvent(aStream Stream, eventCode StreamEvent) {
	objc.Call[objc.Void](s_, objc.Sel("stream:handleEvent:"), aStream, eventCode)
}