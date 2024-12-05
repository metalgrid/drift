// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure a perspective rotate filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate?language=objc
type PPerspectiveRotate interface {
	// optional
	SetRoll(value float32)
	HasSetRoll() bool

	// optional
	Roll() float32
	HasRoll() bool

	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool

	// optional
	SetFocalLength(value float32)
	HasSetFocalLength() bool

	// optional
	FocalLength() float32
	HasFocalLength() bool

	// optional
	SetYaw(value float32)
	HasSetYaw() bool

	// optional
	Yaw() float32
	HasYaw() bool

	// optional
	SetPitch(value float32)
	HasSetPitch() bool

	// optional
	Pitch() float32
	HasPitch() bool
}

// ensure impl type implements protocol interface
var _ PPerspectiveRotate = (*PerspectiveRotateObject)(nil)

// A concrete type for the [PPerspectiveRotate] protocol.
type PerspectiveRotateObject struct {
	objc.Object
}

func (p_ PerspectiveRotateObject) HasSetRoll() bool {
	return p_.RespondsToSelector(objc.Sel("setRoll:"))
}

// The roll angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325540-roll?language=objc
func (p_ PerspectiveRotateObject) SetRoll(value float32) {
	objc.Call[objc.Void](p_, objc.Sel("setRoll:"), value)
}

func (p_ PerspectiveRotateObject) HasRoll() bool {
	return p_.RespondsToSelector(objc.Sel("roll"))
}

// The roll angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325540-roll?language=objc
func (p_ PerspectiveRotateObject) Roll() float32 {
	rv := objc.Call[float32](p_, objc.Sel("roll"))
	return rv
}

func (p_ PerspectiveRotateObject) HasSetInputImage() bool {
	return p_.RespondsToSelector(objc.Sel("setInputImage:"))
}

// The image to process. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325538-inputimage?language=objc
func (p_ PerspectiveRotateObject) SetInputImage(value Image) {
	objc.Call[objc.Void](p_, objc.Sel("setInputImage:"), value)
}

func (p_ PerspectiveRotateObject) HasInputImage() bool {
	return p_.RespondsToSelector(objc.Sel("inputImage"))
}

// The image to process. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325538-inputimage?language=objc
func (p_ PerspectiveRotateObject) InputImage() Image {
	rv := objc.Call[Image](p_, objc.Sel("inputImage"))
	return rv
}

func (p_ PerspectiveRotateObject) HasSetFocalLength() bool {
	return p_.RespondsToSelector(objc.Sel("setFocalLength:"))
}

// The 35mm equivalent focal length of the input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325537-focallength?language=objc
func (p_ PerspectiveRotateObject) SetFocalLength(value float32) {
	objc.Call[objc.Void](p_, objc.Sel("setFocalLength:"), value)
}

func (p_ PerspectiveRotateObject) HasFocalLength() bool {
	return p_.RespondsToSelector(objc.Sel("focalLength"))
}

// The 35mm equivalent focal length of the input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325537-focallength?language=objc
func (p_ PerspectiveRotateObject) FocalLength() float32 {
	rv := objc.Call[float32](p_, objc.Sel("focalLength"))
	return rv
}

func (p_ PerspectiveRotateObject) HasSetYaw() bool {
	return p_.RespondsToSelector(objc.Sel("setYaw:"))
}

// The yaw angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325541-yaw?language=objc
func (p_ PerspectiveRotateObject) SetYaw(value float32) {
	objc.Call[objc.Void](p_, objc.Sel("setYaw:"), value)
}

func (p_ PerspectiveRotateObject) HasYaw() bool {
	return p_.RespondsToSelector(objc.Sel("yaw"))
}

// The yaw angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325541-yaw?language=objc
func (p_ PerspectiveRotateObject) Yaw() float32 {
	rv := objc.Call[float32](p_, objc.Sel("yaw"))
	return rv
}

func (p_ PerspectiveRotateObject) HasSetPitch() bool {
	return p_.RespondsToSelector(objc.Sel("setPitch:"))
}

// The pitch angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325539-pitch?language=objc
func (p_ PerspectiveRotateObject) SetPitch(value float32) {
	objc.Call[objc.Void](p_, objc.Sel("setPitch:"), value)
}

func (p_ PerspectiveRotateObject) HasPitch() bool {
	return p_.RespondsToSelector(objc.Sel("pitch"))
}

// The pitch angle, in radians. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciperspectiverotate/3325539-pitch?language=objc
func (p_ PerspectiveRotateObject) Pitch() float32 {
	rv := objc.Call[float32](p_, objc.Sel("pitch"))
	return rv
}
