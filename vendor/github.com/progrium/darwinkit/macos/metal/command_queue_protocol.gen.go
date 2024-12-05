// Code generated by DarwinKit. DO NOT EDIT.

package metal

import (
	"github.com/progrium/darwinkit/objc"
)

// An instance you use to create, submit, and schedule command buffers to a specific GPU device to run the commands within those buffers. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue?language=objc
type PCommandQueue interface {
	// optional
	CommandBuffer() CommandBufferObject
	HasCommandBuffer() bool

	// optional
	CommandBufferWithUnretainedReferences() CommandBufferObject
	HasCommandBufferWithUnretainedReferences() bool

	// optional
	CommandBufferWithDescriptor(descriptor CommandBufferDescriptor) CommandBufferObject
	HasCommandBufferWithDescriptor() bool

	// optional
	SetLabel(value string)
	HasSetLabel() bool

	// optional
	Label() string
	HasLabel() bool

	// optional
	Device() DeviceObject
	HasDevice() bool
}

// ensure impl type implements protocol interface
var _ PCommandQueue = (*CommandQueueObject)(nil)

// A concrete type for the [PCommandQueue] protocol.
type CommandQueueObject struct {
	objc.Object
}

func (c_ CommandQueueObject) HasCommandBuffer() bool {
	return c_.RespondsToSelector(objc.Sel("commandBuffer"))
}

// Returns a command buffer from the command queue that maintains strong references to resources. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/1508686-commandbuffer?language=objc
func (c_ CommandQueueObject) CommandBuffer() CommandBufferObject {
	rv := objc.Call[CommandBufferObject](c_, objc.Sel("commandBuffer"))
	return rv
}

func (c_ CommandQueueObject) HasCommandBufferWithUnretainedReferences() bool {
	return c_.RespondsToSelector(objc.Sel("commandBufferWithUnretainedReferences"))
}

// Returns a command buffer from the command queue that doesn’t maintain strong references to resources. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/1508684-commandbufferwithunretainedrefer?language=objc
func (c_ CommandQueueObject) CommandBufferWithUnretainedReferences() CommandBufferObject {
	rv := objc.Call[CommandBufferObject](c_, objc.Sel("commandBufferWithUnretainedReferences"))
	return rv
}

func (c_ CommandQueueObject) HasCommandBufferWithDescriptor() bool {
	return c_.RespondsToSelector(objc.Sel("commandBufferWithDescriptor:"))
}

// Returns a command buffer from the command queue that you configure with a descriptor. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/3553957-commandbufferwithdescriptor?language=objc
func (c_ CommandQueueObject) CommandBufferWithDescriptor(descriptor CommandBufferDescriptor) CommandBufferObject {
	rv := objc.Call[CommandBufferObject](c_, objc.Sel("commandBufferWithDescriptor:"), descriptor)
	return rv
}

func (c_ CommandQueueObject) HasSetLabel() bool {
	return c_.RespondsToSelector(objc.Sel("setLabel:"))
}

// An optional name that can help you identify the command queue. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/1508690-label?language=objc
func (c_ CommandQueueObject) SetLabel(value string) {
	objc.Call[objc.Void](c_, objc.Sel("setLabel:"), value)
}

func (c_ CommandQueueObject) HasLabel() bool {
	return c_.RespondsToSelector(objc.Sel("label"))
}

// An optional name that can help you identify the command queue. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/1508690-label?language=objc
func (c_ CommandQueueObject) Label() string {
	rv := objc.Call[string](c_, objc.Sel("label"))
	return rv
}

func (c_ CommandQueueObject) HasDevice() bool {
	return c_.RespondsToSelector(objc.Sel("device"))
}

// The GPU device that creates the command queue. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtlcommandqueue/1508687-device?language=objc
func (c_ CommandQueueObject) Device() DeviceObject {
	rv := objc.Call[DeviceObject](c_, objc.Sel("device"))
	return rv
}
