// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure a thermal filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cithermal?language=objc
type PThermal interface {
	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool
}

// ensure impl type implements protocol interface
var _ PThermal = (*ThermalObject)(nil)

// A concrete type for the [PThermal] protocol.
type ThermalObject struct {
	objc.Object
}

func (t_ ThermalObject) HasSetInputImage() bool {
	return t_.RespondsToSelector(objc.Sel("setInputImage:"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cithermal/3228790-inputimage?language=objc
func (t_ ThermalObject) SetInputImage(value Image) {
	objc.Call[objc.Void](t_, objc.Sel("setInputImage:"), value)
}

func (t_ ThermalObject) HasInputImage() bool {
	return t_.RespondsToSelector(objc.Sel("inputImage"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cithermal/3228790-inputimage?language=objc
func (t_ ThermalObject) InputImage() Image {
	rv := objc.Call[Image](t_, objc.Sel("inputImage"))
	return rv
}
