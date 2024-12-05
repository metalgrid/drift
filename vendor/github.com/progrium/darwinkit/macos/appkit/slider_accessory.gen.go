// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"unsafe"

	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [SliderAccessory] class.
var SliderAccessoryClass = _SliderAccessoryClass{objc.GetClass("NSSliderAccessory")}

type _SliderAccessoryClass struct {
	objc.Class
}

// An interface definition for the [SliderAccessory] class.
type ISliderAccessory interface {
	objc.IObject
	IsEnabled() bool
	SetEnabled(value bool)
	Behavior() SliderAccessoryBehavior
	SetBehavior(value ISliderAccessoryBehavior)
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory?language=objc
type SliderAccessory struct {
	objc.Object
}

func SliderAccessoryFrom(ptr unsafe.Pointer) SliderAccessory {
	return SliderAccessory{
		Object: objc.ObjectFrom(ptr),
	}
}

func (sc _SliderAccessoryClass) Alloc() SliderAccessory {
	rv := objc.Call[SliderAccessory](sc, objc.Sel("alloc"))
	return rv
}

func (sc _SliderAccessoryClass) New() SliderAccessory {
	rv := objc.Call[SliderAccessory](sc, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewSliderAccessory() SliderAccessory {
	return SliderAccessoryClass.New()
}

func (s_ SliderAccessory) Init() SliderAccessory {
	rv := objc.Call[SliderAccessory](s_, objc.Sel("init"))
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544660-accessorywithimage?language=objc
func (sc _SliderAccessoryClass) AccessoryWithImage(image IImage) SliderAccessory {
	rv := objc.Call[SliderAccessory](sc, objc.Sel("accessoryWithImage:"), image)
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544660-accessorywithimage?language=objc
func SliderAccessory_AccessoryWithImage(image IImage) SliderAccessory {
	return SliderAccessoryClass.AccessoryWithImage(image)
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544779-enabled?language=objc
func (s_ SliderAccessory) IsEnabled() bool {
	rv := objc.Call[bool](s_, objc.Sel("isEnabled"))
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544779-enabled?language=objc
func (s_ SliderAccessory) SetEnabled(value bool) {
	objc.Call[objc.Void](s_, objc.Sel("setEnabled:"), value)
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544656-behavior?language=objc
func (s_ SliderAccessory) Behavior() SliderAccessoryBehavior {
	rv := objc.Call[SliderAccessoryBehavior](s_, objc.Sel("behavior"))
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsslideraccessory/2544656-behavior?language=objc
func (s_ SliderAccessory) SetBehavior(value ISliderAccessoryBehavior) {
	objc.Call[objc.Void](s_, objc.Sel("setBehavior:"), value)
}