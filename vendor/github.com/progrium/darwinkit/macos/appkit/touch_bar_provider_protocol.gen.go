// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/objc"
)

// A protocol that an object adopts to create a bar object in your app. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbarprovider?language=objc
type PTouchBarProvider interface {
	// optional
	TouchBar() TouchBar
	HasTouchBar() bool
}

// ensure impl type implements protocol interface
var _ PTouchBarProvider = (*TouchBarProviderObject)(nil)

// A concrete type for the [PTouchBarProvider] protocol.
type TouchBarProviderObject struct {
	objc.Object
}

func (t_ TouchBarProviderObject) HasTouchBar() bool {
	return t_.RespondsToSelector(objc.Sel("touchBar"))
}

// The property you implement to provide a Touch Bar object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbarprovider/2544662-touchbar?language=objc
func (t_ TouchBarProviderObject) TouchBar() TouchBar {
	rv := objc.Call[TouchBar](t_, objc.Sel("touchBar"))
	return rv
}
