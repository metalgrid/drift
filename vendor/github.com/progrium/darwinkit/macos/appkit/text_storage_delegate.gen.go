// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// The optional methods that delegates of text storage objects implement to handle text-edit processing. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate?language=objc
type PTextStorageDelegate interface {
	// optional
	TextStorageDidProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)
	HasTextStorageDidProcessEditingRangeChangeInLength() bool

	// optional
	TextStorageWillProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)
	HasTextStorageWillProcessEditingRangeChangeInLength() bool
}

// A delegate implementation builder for the [PTextStorageDelegate] protocol.
type TextStorageDelegate struct {
	_TextStorageDidProcessEditingRangeChangeInLength  func(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)
	_TextStorageWillProcessEditingRangeChangeInLength func(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)
}

func (di *TextStorageDelegate) HasTextStorageDidProcessEditingRangeChangeInLength() bool {
	return di._TextStorageDidProcessEditingRangeChangeInLength != nil
}

// The method the framework calls when a text storage object has finished processing edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534375-textstorage?language=objc
func (di *TextStorageDelegate) SetTextStorageDidProcessEditingRangeChangeInLength(f func(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)) {
	di._TextStorageDidProcessEditingRangeChangeInLength = f
}

// The method the framework calls when a text storage object has finished processing edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534375-textstorage?language=objc
func (di *TextStorageDelegate) TextStorageDidProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int) {
	di._TextStorageDidProcessEditingRangeChangeInLength(textStorage, editedMask, editedRange, delta)
}
func (di *TextStorageDelegate) HasTextStorageWillProcessEditingRangeChangeInLength() bool {
	return di._TextStorageWillProcessEditingRangeChangeInLength != nil
}

// The method the framework calls when a text storage object is about to process edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534795-textstorage?language=objc
func (di *TextStorageDelegate) SetTextStorageWillProcessEditingRangeChangeInLength(f func(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int)) {
	di._TextStorageWillProcessEditingRangeChangeInLength = f
}

// The method the framework calls when a text storage object is about to process edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534795-textstorage?language=objc
func (di *TextStorageDelegate) TextStorageWillProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int) {
	di._TextStorageWillProcessEditingRangeChangeInLength(textStorage, editedMask, editedRange, delta)
}

// ensure impl type implements protocol interface
var _ PTextStorageDelegate = (*TextStorageDelegateObject)(nil)

// A concrete type for the [PTextStorageDelegate] protocol.
type TextStorageDelegateObject struct {
	objc.Object
}

func (t_ TextStorageDelegateObject) HasTextStorageDidProcessEditingRangeChangeInLength() bool {
	return t_.RespondsToSelector(objc.Sel("textStorage:didProcessEditing:range:changeInLength:"))
}

// The method the framework calls when a text storage object has finished processing edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534375-textstorage?language=objc
func (t_ TextStorageDelegateObject) TextStorageDidProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int) {
	objc.Call[objc.Void](t_, objc.Sel("textStorage:didProcessEditing:range:changeInLength:"), textStorage, editedMask, editedRange, delta)
}

func (t_ TextStorageDelegateObject) HasTextStorageWillProcessEditingRangeChangeInLength() bool {
	return t_.RespondsToSelector(objc.Sel("textStorage:willProcessEditing:range:changeInLength:"))
}

// The method the framework calls when a text storage object is about to process edits. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uikit/nstextstoragedelegate/1534795-textstorage?language=objc
func (t_ TextStorageDelegateObject) TextStorageWillProcessEditingRangeChangeInLength(textStorage TextStorage, editedMask TextStorageEditActions, editedRange foundation.Range, delta int) {
	objc.Call[objc.Void](t_, objc.Sel("textStorage:willProcessEditing:range:changeInLength:"), textStorage, editedMask, editedRange, delta)
}