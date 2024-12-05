// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/macos/coregraphics"
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure an affine tile filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciaffinetile?language=objc
type PAffineTile interface {
	// optional
	SetTransform(value coregraphics.AffineTransform)
	HasSetTransform() bool

	// optional
	Transform() coregraphics.AffineTransform
	HasTransform() bool

	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool
}

// ensure impl type implements protocol interface
var _ PAffineTile = (*AffineTileObject)(nil)

// A concrete type for the [PAffineTile] protocol.
type AffineTileObject struct {
	objc.Object
}

func (a_ AffineTileObject) HasSetTransform() bool {
	return a_.RespondsToSelector(objc.Sel("setTransform:"))
}

// The transform to apply to the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciaffinetile/3228058-transform?language=objc
func (a_ AffineTileObject) SetTransform(value coregraphics.AffineTransform) {
	objc.Call[objc.Void](a_, objc.Sel("setTransform:"), value)
}

func (a_ AffineTileObject) HasTransform() bool {
	return a_.RespondsToSelector(objc.Sel("transform"))
}

// The transform to apply to the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciaffinetile/3228058-transform?language=objc
func (a_ AffineTileObject) Transform() coregraphics.AffineTransform {
	rv := objc.Call[coregraphics.AffineTransform](a_, objc.Sel("transform"))
	return rv
}

func (a_ AffineTileObject) HasSetInputImage() bool {
	return a_.RespondsToSelector(objc.Sel("setInputImage:"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciaffinetile/3228057-inputimage?language=objc
func (a_ AffineTileObject) SetInputImage(value Image) {
	objc.Call[objc.Void](a_, objc.Sel("setInputImage:"), value)
}

func (a_ AffineTileObject) HasInputImage() bool {
	return a_.RespondsToSelector(objc.Sel("inputImage"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciaffinetile/3228057-inputimage?language=objc
func (a_ AffineTileObject) InputImage() Image {
	rv := objc.Call[Image](a_, objc.Sel("inputImage"))
	return rv
}
