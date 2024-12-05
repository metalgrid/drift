// Code generated by DarwinKit. DO NOT EDIT.

package metal

import (
	"github.com/progrium/darwinkit/objc"
)

// An individual counter a GPU device lists within one of its counter sets. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcounter?language=objc
type PCounter interface {
	// optional
	Name() string
	HasName() bool
}

// ensure impl type implements protocol interface
var _ PCounter = (*CounterObject)(nil)

// A concrete type for the [PCounter] protocol.
type CounterObject struct {
	objc.Object
}

func (c_ CounterObject) HasName() bool {
	return c_.RespondsToSelector(objc.Sel("name"))
}

// The name of a GPU’s counter instance. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcounter/3081701-name?language=objc
func (c_ CounterObject) Name() string {
	rv := objc.Call[string](c_, objc.Sel("name"))
	return rv
}