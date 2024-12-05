// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure a Code 128 barcode generator filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator?language=objc
type PCode128BarcodeGenerator interface {
	// optional
	SetBarcodeHeight(value float32)
	HasSetBarcodeHeight() bool

	// optional
	BarcodeHeight() float32
	HasBarcodeHeight() bool

	// optional
	SetQuietSpace(value float32)
	HasSetQuietSpace() bool

	// optional
	QuietSpace() float32
	HasQuietSpace() bool

	// optional
	SetMessage(value []byte)
	HasSetMessage() bool

	// optional
	Message() []byte
	HasMessage() bool
}

// ensure impl type implements protocol interface
var _ PCode128BarcodeGenerator = (*Code128BarcodeGeneratorObject)(nil)

// A concrete type for the [PCode128BarcodeGenerator] protocol.
type Code128BarcodeGeneratorObject struct {
	objc.Object
}

func (c_ Code128BarcodeGeneratorObject) HasSetBarcodeHeight() bool {
	return c_.RespondsToSelector(objc.Sel("setBarcodeHeight:"))
}

// The height, in pixels, of the generated barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228116-barcodeheight?language=objc
func (c_ Code128BarcodeGeneratorObject) SetBarcodeHeight(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setBarcodeHeight:"), value)
}

func (c_ Code128BarcodeGeneratorObject) HasBarcodeHeight() bool {
	return c_.RespondsToSelector(objc.Sel("barcodeHeight"))
}

// The height, in pixels, of the generated barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228116-barcodeheight?language=objc
func (c_ Code128BarcodeGeneratorObject) BarcodeHeight() float32 {
	rv := objc.Call[float32](c_, objc.Sel("barcodeHeight"))
	return rv
}

func (c_ Code128BarcodeGeneratorObject) HasSetQuietSpace() bool {
	return c_.RespondsToSelector(objc.Sel("setQuietSpace:"))
}

// The number of empty white pixels that should surround the barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228118-quietspace?language=objc
func (c_ Code128BarcodeGeneratorObject) SetQuietSpace(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setQuietSpace:"), value)
}

func (c_ Code128BarcodeGeneratorObject) HasQuietSpace() bool {
	return c_.RespondsToSelector(objc.Sel("quietSpace"))
}

// The number of empty white pixels that should surround the barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228118-quietspace?language=objc
func (c_ Code128BarcodeGeneratorObject) QuietSpace() float32 {
	rv := objc.Call[float32](c_, objc.Sel("quietSpace"))
	return rv
}

func (c_ Code128BarcodeGeneratorObject) HasSetMessage() bool {
	return c_.RespondsToSelector(objc.Sel("setMessage:"))
}

// The message to encode in the Code 128 barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228117-message?language=objc
func (c_ Code128BarcodeGeneratorObject) SetMessage(value []byte) {
	objc.Call[objc.Void](c_, objc.Sel("setMessage:"), value)
}

func (c_ Code128BarcodeGeneratorObject) HasMessage() bool {
	return c_.RespondsToSelector(objc.Sel("message"))
}

// The message to encode in the Code 128 barcode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicode128barcodegenerator/3228117-message?language=objc
func (c_ Code128BarcodeGeneratorObject) Message() []byte {
	rv := objc.Call[[]byte](c_, objc.Sel("message"))
	return rv
}
