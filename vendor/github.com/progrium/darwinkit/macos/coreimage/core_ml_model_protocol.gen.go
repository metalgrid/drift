// Code generated by DarwinKit. DO NOT EDIT.

package coreimage

import (
	"github.com/progrium/darwinkit/macos/coreml"
	"github.com/progrium/darwinkit/objc"
)

// The properties you use to configure a Core ML model filter. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel?language=objc
type PCoreMLModel interface {
	// optional
	SetSoftmaxNormalization(value bool)
	HasSetSoftmaxNormalization() bool

	// optional
	SoftmaxNormalization() bool
	HasSoftmaxNormalization() bool

	// optional
	SetHeadIndex(value float32)
	HasSetHeadIndex() bool

	// optional
	HeadIndex() float32
	HasHeadIndex() bool

	// optional
	SetModel(value coreml.Model)
	HasSetModel() bool

	// optional
	Model() coreml.Model
	HasModel() bool

	// optional
	SetInputImage(value Image)
	HasSetInputImage() bool

	// optional
	InputImage() Image
	HasInputImage() bool
}

// ensure impl type implements protocol interface
var _ PCoreMLModel = (*CoreMLModelObject)(nil)

// A concrete type for the [PCoreMLModel] protocol.
type CoreMLModelObject struct {
	objc.Object
}

func (c_ CoreMLModelObject) HasSetSoftmaxNormalization() bool {
	return c_.RespondsToSelector(objc.Sel("setSoftmaxNormalization:"))
}

// A Boolean value that specifies whether to apply Softmax normalization to the output of the model. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228198-softmaxnormalization?language=objc
func (c_ CoreMLModelObject) SetSoftmaxNormalization(value bool) {
	objc.Call[objc.Void](c_, objc.Sel("setSoftmaxNormalization:"), value)
}

func (c_ CoreMLModelObject) HasSoftmaxNormalization() bool {
	return c_.RespondsToSelector(objc.Sel("softmaxNormalization"))
}

// A Boolean value that specifies whether to apply Softmax normalization to the output of the model. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228198-softmaxnormalization?language=objc
func (c_ CoreMLModelObject) SoftmaxNormalization() bool {
	rv := objc.Call[bool](c_, objc.Sel("softmaxNormalization"))
	return rv
}

func (c_ CoreMLModelObject) HasSetHeadIndex() bool {
	return c_.RespondsToSelector(objc.Sel("setHeadIndex:"))
}

// A number that specifies which output of a multihead Core ML model applies the effect on the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228195-headindex?language=objc
func (c_ CoreMLModelObject) SetHeadIndex(value float32) {
	objc.Call[objc.Void](c_, objc.Sel("setHeadIndex:"), value)
}

func (c_ CoreMLModelObject) HasHeadIndex() bool {
	return c_.RespondsToSelector(objc.Sel("headIndex"))
}

// A number that specifies which output of a multihead Core ML model applies the effect on the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228195-headindex?language=objc
func (c_ CoreMLModelObject) HeadIndex() float32 {
	rv := objc.Call[float32](c_, objc.Sel("headIndex"))
	return rv
}

func (c_ CoreMLModelObject) HasSetModel() bool {
	return c_.RespondsToSelector(objc.Sel("setModel:"))
}

// The Core ML model used to apply the effect on the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228197-model?language=objc
func (c_ CoreMLModelObject) SetModel(value coreml.Model) {
	objc.Call[objc.Void](c_, objc.Sel("setModel:"), value)
}

func (c_ CoreMLModelObject) HasModel() bool {
	return c_.RespondsToSelector(objc.Sel("model"))
}

// The Core ML model used to apply the effect on the image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228197-model?language=objc
func (c_ CoreMLModelObject) Model() coreml.Model {
	rv := objc.Call[coreml.Model](c_, objc.Sel("model"))
	return rv
}

func (c_ CoreMLModelObject) HasSetInputImage() bool {
	return c_.RespondsToSelector(objc.Sel("setInputImage:"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228196-inputimage?language=objc
func (c_ CoreMLModelObject) SetInputImage(value Image) {
	objc.Call[objc.Void](c_, objc.Sel("setInputImage:"), value)
}

func (c_ CoreMLModelObject) HasInputImage() bool {
	return c_.RespondsToSelector(objc.Sel("inputImage"))
}

// The image to use as an input image. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreimage/cicoremlmodel/3228196-inputimage?language=objc
func (c_ CoreMLModelObject) InputImage() Image {
	rv := objc.Call[Image](c_, objc.Sel("inputImage"))
	return rv
}
