// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/macos/coregraphics"
	"github.com/progrium/darwinkit/objc"
)

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge?language=objc
type PGlassLozenge interface {
	// optional
	SetRadius(value float32)
	HasSetRadius() bool

	// optional
	Radius() float32
	HasRadius() bool

	// optional
	SetRefraction(value float32)
	HasSetRefraction() bool

	// optional
	Refraction() float32
	HasRefraction() bool

	// optional
	SetPoint0(value coregraphics.Point)
	HasSetPoint0() bool

	// optional
	Point0() coregraphics.Point
	HasPoint0() bool

	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool

	// optional
	SetPoint1(value coregraphics.Point)
	HasSetPoint1() bool

	// optional
	Point1() coregraphics.Point
	HasPoint1() bool
}

// ensure impl type implements protocol interface
var _ PGlassLozenge = (*GlassLozengeObject)(nil)

// A concrete type for the [PGlassLozenge] protocol.
type GlassLozengeObject struct {
	objc.Object
}

func (g_ GlassLozengeObject) HasSetRadius() bool {
	return g_.RespondsToSelector(objc.Sel("setRadius:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600164-radius?language=objc
func (g_ GlassLozengeObject) SetRadius(value float32) {
	objc.Call[objc.Void](g_, objc.Sel("setRadius:"), value)
}

func (g_ GlassLozengeObject) HasRadius() bool {
	return g_.RespondsToSelector(objc.Sel("radius"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600164-radius?language=objc
func (g_ GlassLozengeObject) Radius() float32 {
	rv := objc.Call[float32](g_, objc.Sel("radius"))
	return rv
}

func (g_ GlassLozengeObject) HasSetRefraction() bool {
	return g_.RespondsToSelector(objc.Sel("setRefraction:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600165-refraction?language=objc
func (g_ GlassLozengeObject) SetRefraction(value float32) {
	objc.Call[objc.Void](g_, objc.Sel("setRefraction:"), value)
}

func (g_ GlassLozengeObject) HasRefraction() bool {
	return g_.RespondsToSelector(objc.Sel("refraction"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600165-refraction?language=objc
func (g_ GlassLozengeObject) Refraction() float32 {
	rv := objc.Call[float32](g_, objc.Sel("refraction"))
	return rv
}

func (g_ GlassLozengeObject) HasSetPoint0() bool {
	return g_.RespondsToSelector(objc.Sel("setPoint0:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600162-point0?language=objc
func (g_ GlassLozengeObject) SetPoint0(value coregraphics.Point) {
	objc.Call[objc.Void](g_, objc.Sel("setPoint0:"), value)
}

func (g_ GlassLozengeObject) HasPoint0() bool {
	return g_.RespondsToSelector(objc.Sel("point0"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600162-point0?language=objc
func (g_ GlassLozengeObject) Point0() coregraphics.Point {
	rv := objc.Call[coregraphics.Point](g_, objc.Sel("point0"))
	return rv
}

func (g_ GlassLozengeObject) HasSetInputImage() bool {
	return g_.RespondsToSelector(objc.Sel("setInputImage:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600161-inputimage?language=objc
func (g_ GlassLozengeObject) SetInputImage(value Image) {
	objc.Call[objc.Void](g_, objc.Sel("setInputImage:"), value)
}

func (g_ GlassLozengeObject) HasInputImage() bool {
	return g_.RespondsToSelector(objc.Sel("inputImage"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600161-inputimage?language=objc
func (g_ GlassLozengeObject) InputImage() Image {
	rv := objc.Call[Image](g_, objc.Sel("inputImage"))
	return rv
}

func (g_ GlassLozengeObject) HasSetPoint1() bool {
	return g_.RespondsToSelector(objc.Sel("setPoint1:"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600163-point1?language=objc
func (g_ GlassLozengeObject) SetPoint1(value coregraphics.Point) {
	objc.Call[objc.Void](g_, objc.Sel("setPoint1:"), value)
}

func (g_ GlassLozengeObject) HasPoint1() bool {
	return g_.RespondsToSelector(objc.Sel("point1"))
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/ciglasslozenge/3600163-point1?language=objc
func (g_ GlassLozengeObject) Point1() coregraphics.Point {
	rv := objc.Call[coregraphics.Point](g_, objc.Sel("point1"))
	return rv
}
