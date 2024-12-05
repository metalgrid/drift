// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// A role-based protocol that declares the minimum interface necessary for an accessibility element to act as a checkbox. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsaccessibilitycheckbox?language=objc
type PAccessibilityCheckBox interface {
	// optional
	AccessibilityValue() foundation.Number
	HasAccessibilityValue() bool
}

// ensure impl type implements protocol interface
var _ PAccessibilityCheckBox = (*AccessibilityCheckBoxObject)(nil)

// A concrete type for the [PAccessibilityCheckBox] protocol.
type AccessibilityCheckBoxObject struct {
	objc.Object
}

func (a_ AccessibilityCheckBoxObject) HasAccessibilityValue() bool {
	return a_.RespondsToSelector(objc.Sel("accessibilityValue"))
}

// Returns the checkbox’s value. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsaccessibilitycheckbox/1524299-accessibilityvalue?language=objc
func (a_ AccessibilityCheckBoxObject) AccessibilityValue() foundation.Number {
	rv := objc.Call[foundation.Number](a_, objc.Sel("accessibilityValue"))
	return rv
}
