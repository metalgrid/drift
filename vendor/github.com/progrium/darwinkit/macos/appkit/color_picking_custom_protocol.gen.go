// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/objc"
)

// A set of methods that provides a way to add color pickers—custom user interfaces for color selection—to an app’s color panel. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nscolorpickingcustom?language=objc
type PColorPickingCustom interface {
	// optional
	SupportsMode(mode ColorPanelMode) bool
	HasSupportsMode() bool

	// optional
	CurrentMode() ColorPanelMode
	HasCurrentMode() bool

	// optional
	SetColor(newColor Color)
	HasSetColor() bool

	// optional
	ProvideNewView(initialRequest bool) View
	HasProvideNewView() bool
}

// ensure impl type implements protocol interface
var _ PColorPickingCustom = (*ColorPickingCustomObject)(nil)

// A concrete type for the [PColorPickingCustom] protocol.
type ColorPickingCustomObject struct {
	objc.Object
}

func (c_ ColorPickingCustomObject) HasSupportsMode() bool {
	return c_.RespondsToSelector(objc.Sel("supportsMode:"))
}

// Returns a Boolean value indicating whether or not the receiver supports the specified picking mode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nscolorpickingcustom/1524683-supportsmode?language=objc
func (c_ ColorPickingCustomObject) SupportsMode(mode ColorPanelMode) bool {
	rv := objc.Call[bool](c_, objc.Sel("supportsMode:"), mode)
	return rv
}

func (c_ ColorPickingCustomObject) HasCurrentMode() bool {
	return c_.RespondsToSelector(objc.Sel("currentMode"))
}

// Returns the receiver’s current mode (or submode, if applicable). [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nscolorpickingcustom/1524671-currentmode?language=objc
func (c_ ColorPickingCustomObject) CurrentMode() ColorPanelMode {
	rv := objc.Call[ColorPanelMode](c_, objc.Sel("currentMode"))
	return rv
}

func (c_ ColorPickingCustomObject) HasSetColor() bool {
	return c_.RespondsToSelector(objc.Sel("setColor:"))
}

// Adjusts the receiver to make the specified color the currently selected color. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nscolorpickingcustom/1526545-setcolor?language=objc
func (c_ ColorPickingCustomObject) SetColor(newColor Color) {
	objc.Call[objc.Void](c_, objc.Sel("setColor:"), newColor)
}

func (c_ ColorPickingCustomObject) HasProvideNewView() bool {
	return c_.RespondsToSelector(objc.Sel("provideNewView:"))
}

// Returns the view containing the receiver’s user interface. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nscolorpickingcustom/1525701-providenewview?language=objc
func (c_ ColorPickingCustomObject) ProvideNewView(initialRequest bool) View {
	rv := objc.Call[View](c_, objc.Sel("provideNewView:"), initialRequest)
	return rv
}
