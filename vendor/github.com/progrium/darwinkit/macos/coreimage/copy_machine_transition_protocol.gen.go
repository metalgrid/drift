// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/macos/coregraphics"
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure a copy machine transition filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition?language=objc
type PCopyMachineTransition interface {
	// optional
	SetColor(value Color)
	HasSetColor() bool

	// optional
	Color() Color
	HasColor() bool

	// optional
	SetExtent(value coregraphics.Rect)
	HasSetExtent() bool

	// optional
	Extent() coregraphics.Rect
	HasExtent() bool

	// optional
	SetOpacity(value float32)
	HasSetOpacity() bool

	// optional
	Opacity() float32
	HasOpacity() bool

	// optional
	SetWidth(value float32)
	HasSetWidth() bool

	// optional
	Width() float32
	HasWidth() bool

	// optional
	SetAngle(value float32)
	HasSetAngle() bool

	// optional
	Angle() float32
	HasAngle() bool
}

// ensure impl type implements protocol interface
var _ PCopyMachineTransition = (*CopyMachineTransitionObject)(nil)

// A concrete type for the [PCopyMachineTransition] protocol.
type CopyMachineTransitionObject struct {
	objc.Object
}

func (c_ CopyMachineTransitionObject) HasSetColor() bool {
	return c_.RespondsToSelector(objc.Sel("setColor:"))
}

// The color of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228190-color?language=objc
func (c_ CopyMachineTransitionObject) SetColor(value Color) {
	objc.Call[objc.Void](c_, objc.Sel("setColor:"), value)
}

func (c_ CopyMachineTransitionObject) HasColor() bool {
	return c_.RespondsToSelector(objc.Sel("color"))
}

// The color of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228190-color?language=objc
func (c_ CopyMachineTransitionObject) Color() Color {
	rv := objc.Call[Color](c_, objc.Sel("color"))
	return rv
}

func (c_ CopyMachineTransitionObject) HasSetExtent() bool {
	return c_.RespondsToSelector(objc.Sel("setExtent:"))
}

// A rectangle that defines the extent of the effect. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228191-extent?language=objc
func (c_ CopyMachineTransitionObject) SetExtent(value coregraphics.Rect) {
	objc.Call[objc.Void](c_, objc.Sel("setExtent:"), value)
}

func (c_ CopyMachineTransitionObject) HasExtent() bool {
	return c_.RespondsToSelector(objc.Sel("extent"))
}

// A rectangle that defines the extent of the effect. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228191-extent?language=objc
func (c_ CopyMachineTransitionObject) Extent() coregraphics.Rect {
	rv := objc.Call[coregraphics.Rect](c_, objc.Sel("extent"))
	return rv
}

func (c_ CopyMachineTransitionObject) HasSetOpacity() bool {
	return c_.RespondsToSelector(objc.Sel("setOpacity:"))
}

// The opacity of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228192-opacity?language=objc
func (c_ CopyMachineTransitionObject) SetOpacity(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setOpacity:"), value)
}

func (c_ CopyMachineTransitionObject) HasOpacity() bool {
	return c_.RespondsToSelector(objc.Sel("opacity"))
}

// The opacity of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228192-opacity?language=objc
func (c_ CopyMachineTransitionObject) Opacity() float32 {
	rv := objc.Call[float32](c_, objc.Sel("opacity"))
	return rv
}

func (c_ CopyMachineTransitionObject) HasSetWidth() bool {
	return c_.RespondsToSelector(objc.Sel("setWidth:"))
}

// The width of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228193-width?language=objc
func (c_ CopyMachineTransitionObject) SetWidth(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setWidth:"), value)
}

func (c_ CopyMachineTransitionObject) HasWidth() bool {
	return c_.RespondsToSelector(objc.Sel("width"))
}

// The width of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228193-width?language=objc
func (c_ CopyMachineTransitionObject) Width() float32 {
	rv := objc.Call[float32](c_, objc.Sel("width"))
	return rv
}

func (c_ CopyMachineTransitionObject) HasSetAngle() bool {
	return c_.RespondsToSelector(objc.Sel("setAngle:"))
}

// The angle of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228189-angle?language=objc
func (c_ CopyMachineTransitionObject) SetAngle(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setAngle:"), value)
}

func (c_ CopyMachineTransitionObject) HasAngle() bool {
	return c_.RespondsToSelector(objc.Sel("angle"))
}

// The angle of the copier light. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicopymachinetransition/3228189-angle?language=objc
func (c_ CopyMachineTransitionObject) Angle() float32 {
	rv := objc.Call[float32](c_, objc.Sel("angle"))
	return rv
}
