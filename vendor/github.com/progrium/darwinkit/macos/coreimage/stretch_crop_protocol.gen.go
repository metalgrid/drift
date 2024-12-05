// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/macos/coregraphics"
	"github.com/progrium/darwinkit/objc"
)

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop?language=objc
type PStretchCrop interface {
	// optional
	SetSize(value coregraphics.Point)
	HasSetSize() bool

	// optional
	Size() coregraphics.Point
	HasSize() bool

	// optional
	SetCropAmount(value float32)
	HasSetCropAmount() bool

	// optional
	CropAmount() float32
	HasCropAmount() bool

	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool

	// optional
	SetCenterStretchAmount(value float32)
	HasSetCenterStretchAmount() bool

	// optional
	CenterStretchAmount() float32
	HasCenterStretchAmount() bool
}

// ensure impl type implements protocol interface
var _ PStretchCrop = (*StretchCropObject)(nil)

// A concrete type for the [PStretchCrop] protocol.
type StretchCropObject struct {
	objc.Object
}

func (s_ StretchCropObject) HasSetSize() bool {
	return s_.RespondsToSelector(objc.Sel("setSize:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600195-size?language=objc
func (s_ StretchCropObject) SetSize(value coregraphics.Point) {
	objc.Call[objc.Void](s_, objc.Sel("setSize:"), value)
}

func (s_ StretchCropObject) HasSize() bool {
	return s_.RespondsToSelector(objc.Sel("size"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600195-size?language=objc
func (s_ StretchCropObject) Size() coregraphics.Point {
	rv := objc.Call[coregraphics.Point](s_, objc.Sel("size"))
	return rv
}

func (s_ StretchCropObject) HasSetCropAmount() bool {
	return s_.RespondsToSelector(objc.Sel("setCropAmount:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600193-cropamount?language=objc
func (s_ StretchCropObject) SetCropAmount(value float32) {
	objc.Call[objc.Void](s_, objc.Sel("setCropAmount:"), value)
}

func (s_ StretchCropObject) HasCropAmount() bool {
	return s_.RespondsToSelector(objc.Sel("cropAmount"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600193-cropamount?language=objc
func (s_ StretchCropObject) CropAmount() float32 {
	rv := objc.Call[float32](s_, objc.Sel("cropAmount"))
	return rv
}

func (s_ StretchCropObject) HasSetInputImage() bool {
	return s_.RespondsToSelector(objc.Sel("setInputImage:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600194-inputimage?language=objc
func (s_ StretchCropObject) SetInputImage(value Image) {
	objc.Call[objc.Void](s_, objc.Sel("setInputImage:"), value)
}

func (s_ StretchCropObject) HasInputImage() bool {
	return s_.RespondsToSelector(objc.Sel("inputImage"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600194-inputimage?language=objc
func (s_ StretchCropObject) InputImage() Image {
	rv := objc.Call[Image](s_, objc.Sel("inputImage"))
	return rv
}

func (s_ StretchCropObject) HasSetCenterStretchAmount() bool {
	return s_.RespondsToSelector(objc.Sel("setCenterStretchAmount:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600192-centerstretchamount?language=objc
func (s_ StretchCropObject) SetCenterStretchAmount(value float32) {
	objc.Call[objc.Void](s_, objc.Sel("setCenterStretchAmount:"), value)
}

func (s_ StretchCropObject) HasCenterStretchAmount() bool {
	return s_.RespondsToSelector(objc.Sel("centerStretchAmount"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cistretchcrop/3600192-centerstretchamount?language=objc
func (s_ StretchCropObject) CenterStretchAmount() float32 {
	rv := objc.Call[float32](s_, objc.Sel("centerStretchAmount"))
	return rv
}
