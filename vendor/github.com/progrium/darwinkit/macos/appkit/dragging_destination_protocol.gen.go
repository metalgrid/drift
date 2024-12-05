// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"github.com/progrium/darwinkit/objc"
)

// A set of methods that the destination object (or recipient) of a dragged image must implement. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination?language=objc
type PDraggingDestination interface {
	// optional
	DraggingExited(sender DraggingInfoObject)
	HasDraggingExited() bool

	// optional
	WantsPeriodicDraggingUpdates() bool
	HasWantsPeriodicDraggingUpdates() bool

	// optional
	ConcludeDragOperation(sender DraggingInfoObject)
	HasConcludeDragOperation() bool

	// optional
	DraggingEntered(sender DraggingInfoObject) DragOperation
	HasDraggingEntered() bool

	// optional
	PerformDragOperation(sender DraggingInfoObject) bool
	HasPerformDragOperation() bool

	// optional
	UpdateDraggingItemsForDrag(sender DraggingInfoObject)
	HasUpdateDraggingItemsForDrag() bool

	// optional
	PrepareForDragOperation(sender DraggingInfoObject) bool
	HasPrepareForDragOperation() bool

	// optional
	DraggingUpdated(sender DraggingInfoObject) DragOperation
	HasDraggingUpdated() bool

	// optional
	DraggingEnded(sender DraggingInfoObject)
	HasDraggingEnded() bool
}

// ensure impl type implements protocol interface
var _ PDraggingDestination = (*DraggingDestinationObject)(nil)

// A concrete type for the [PDraggingDestination] protocol.
type DraggingDestinationObject struct {
	objc.Object
}

func (d_ DraggingDestinationObject) HasDraggingExited() bool {
	return d_.RespondsToSelector(objc.Sel("draggingExited:"))
}

// Invoked when the dragged image exits the destination’s bounds rectangle (in the case of a view object) or its frame rectangle (in the case of a window object). [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416056-draggingexited?language=objc
func (d_ DraggingDestinationObject) DraggingExited(sender DraggingInfoObject) {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	objc.Call[objc.Void](d_, objc.Sel("draggingExited:"), po0)
}

func (d_ DraggingDestinationObject) HasWantsPeriodicDraggingUpdates() bool {
	return d_.RespondsToSelector(objc.Sel("wantsPeriodicDraggingUpdates"))
}

// Asks the destination object whether it wants to receive periodic [appkit/nsdraggingdestination/draggingupdated] messages. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416049-wantsperiodicdraggingupdates?language=objc
func (d_ DraggingDestinationObject) WantsPeriodicDraggingUpdates() bool {
	rv := objc.Call[bool](d_, objc.Sel("wantsPeriodicDraggingUpdates"))
	return rv
}

func (d_ DraggingDestinationObject) HasConcludeDragOperation() bool {
	return d_.RespondsToSelector(objc.Sel("concludeDragOperation:"))
}

// Invoked when the dragging operation is complete, signaling the receiver to perform any necessary clean-up. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416010-concludedragoperation?language=objc
func (d_ DraggingDestinationObject) ConcludeDragOperation(sender DraggingInfoObject) {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	objc.Call[objc.Void](d_, objc.Sel("concludeDragOperation:"), po0)
}

func (d_ DraggingDestinationObject) HasDraggingEntered() bool {
	return d_.RespondsToSelector(objc.Sel("draggingEntered:"))
}

// Invoked when the dragged image enters destination bounds or frame; delegate returns dragging operation to perform. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416019-draggingentered?language=objc
func (d_ DraggingDestinationObject) DraggingEntered(sender DraggingInfoObject) DragOperation {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	rv := objc.Call[DragOperation](d_, objc.Sel("draggingEntered:"), po0)
	return rv
}

func (d_ DraggingDestinationObject) HasPerformDragOperation() bool {
	return d_.RespondsToSelector(objc.Sel("performDragOperation:"))
}

// Invoked after the released image has been removed from the screen, signaling the receiver to import the pasteboard data. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1415970-performdragoperation?language=objc
func (d_ DraggingDestinationObject) PerformDragOperation(sender DraggingInfoObject) bool {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	rv := objc.Call[bool](d_, objc.Sel("performDragOperation:"), po0)
	return rv
}

func (d_ DraggingDestinationObject) HasUpdateDraggingItemsForDrag() bool {
	return d_.RespondsToSelector(objc.Sel("updateDraggingItemsForDrag:"))
}

// Invoked when the dragging images should be changed. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416050-updatedraggingitemsfordrag?language=objc
func (d_ DraggingDestinationObject) UpdateDraggingItemsForDrag(sender DraggingInfoObject) {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	objc.Call[objc.Void](d_, objc.Sel("updateDraggingItemsForDrag:"), po0)
}

func (d_ DraggingDestinationObject) HasPrepareForDragOperation() bool {
	return d_.RespondsToSelector(objc.Sel("prepareForDragOperation:"))
}

// Invoked when the image is released, allowing the receiver to agree to or refuse drag operation. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416066-preparefordragoperation?language=objc
func (d_ DraggingDestinationObject) PrepareForDragOperation(sender DraggingInfoObject) bool {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	rv := objc.Call[bool](d_, objc.Sel("prepareForDragOperation:"), po0)
	return rv
}

func (d_ DraggingDestinationObject) HasDraggingUpdated() bool {
	return d_.RespondsToSelector(objc.Sel("draggingUpdated:"))
}

// Invoked periodically as the image is held within the destination area, allowing modification of the dragging operation or mouse-pointer position. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1415998-draggingupdated?language=objc
func (d_ DraggingDestinationObject) DraggingUpdated(sender DraggingInfoObject) DragOperation {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	rv := objc.Call[DragOperation](d_, objc.Sel("draggingUpdated:"), po0)
	return rv
}

func (d_ DraggingDestinationObject) HasDraggingEnded() bool {
	return d_.RespondsToSelector(objc.Sel("draggingEnded:"))
}

// Called when a drag operation ends. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nsdraggingdestination/1416096-draggingended?language=objc
func (d_ DraggingDestinationObject) DraggingEnded(sender DraggingInfoObject) {
	po0 := objc.WrapAsProtocol("NSDraggingInfo", sender)
	objc.Call[objc.Void](d_, objc.Sel("draggingEnded:"), po0)
}
