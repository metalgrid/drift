// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure an sRGB-to-linear filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cisrgbtonecurvetolinear?language=objc
type PSRGBToneCurveToLinear interface {
	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool
}

// ensure impl type implements protocol interface
var _ PSRGBToneCurveToLinear = (*SRGBToneCurveToLinearObject)(nil)

// A concrete type for the [PSRGBToneCurveToLinear] protocol.
type SRGBToneCurveToLinearObject struct {
	objc.Object
}

func (s_ SRGBToneCurveToLinearObject) HasSetInputImage() bool {
	return s_.RespondsToSelector(objc.Sel("setInputImage:"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cisrgbtonecurvetolinear/3228698-inputimage?language=objc
func (s_ SRGBToneCurveToLinearObject) SetInputImage(value Image) {
	objc.Call[objc.Void](s_, objc.Sel("setInputImage:"), value)
}

func (s_ SRGBToneCurveToLinearObject) HasInputImage() bool {
	return s_.RespondsToSelector(objc.Sel("inputImage"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cisrgbtonecurvetolinear/3228698-inputimage?language=objc
func (s_ SRGBToneCurveToLinearObject) InputImage() Image {
	rv := objc.Call[Image](s_, objc.Sel("inputImage"))
	return rv
}
