// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/objc"
)

// A set of methods that text views need to implement to interact properly with the text input management system. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstextinput?language=objc
type PTextInput interface {
}

// ensure impl type implements protocol interface
var _ PTextInput = (*TextInputObject)(nil)

// A concrete type for the [PTextInput] protocol.
type TextInputObject struct {
	objc.Object
}