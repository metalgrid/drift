// Code generated by DarwinKit. DO NOT EDIT.

package metal

import (
	"unsafe"

	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [TextureDescriptor] class.
var TextureDescriptorClass = _TextureDescriptorClass{objc.GetClass("MTLTextureDescriptor")}

type _TextureDescriptorClass struct {
	objc.Class
}

// An interface definition for the [TextureDescriptor] class.
type ITextureDescriptor interface {
	objc.IObject
	TextureType() TextureType
	SetTextureType(value TextureType)
	MipmapLevelCount() uint
	SetMipmapLevelCount(value uint)
	ResourceOptions() ResourceOptions
	SetResourceOptions(value ResourceOptions)
	PixelFormat() PixelFormat
	SetPixelFormat(value PixelFormat)
	Swizzle() TextureSwizzleChannels
	SetSwizzle(value TextureSwizzleChannels)
	HazardTrackingMode() HazardTrackingMode
	SetHazardTrackingMode(value HazardTrackingMode)
	ArrayLength() uint
	SetArrayLength(value uint)
	AllowGPUOptimizedContents() bool
	SetAllowGPUOptimizedContents(value bool)
	Height() uint
	SetHeight(value uint)
	StorageMode() StorageMode
	SetStorageMode(value StorageMode)
	Width() uint
	SetWidth(value uint)
	CpuCacheMode() CPUCacheMode
	SetCpuCacheMode(value CPUCacheMode)
	Usage() TextureUsage
	SetUsage(value TextureUsage)
	Depth() uint
	SetDepth(value uint)
	CompressionType() TextureCompressionType
	SetCompressionType(value TextureCompressionType)
	SampleCount() uint
	SetSampleCount(value uint)
}

// An object that you use to configure new Metal texture objects. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor?language=objc
type TextureDescriptor struct {
	objc.Object
}

func TextureDescriptorFrom(ptr unsafe.Pointer) TextureDescriptor {
	return TextureDescriptor{
		Object: objc.ObjectFrom(ptr),
	}
}

func (tc _TextureDescriptorClass) Alloc() TextureDescriptor {
	rv := objc.Call[TextureDescriptor](tc, objc.Sel("alloc"))
	return rv
}

func (tc _TextureDescriptorClass) New() TextureDescriptor {
	rv := objc.Call[TextureDescriptor](tc, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewTextureDescriptor() TextureDescriptor {
	return TextureDescriptorClass.New()
}

func (t_ TextureDescriptor) Init() TextureDescriptor {
	rv := objc.Call[TextureDescriptor](t_, objc.Sel("init"))
	return rv
}

// Creates a texture descriptor object for a cube texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516090-texturecubedescriptorwithpixelfo?language=objc
func (tc _TextureDescriptorClass) TextureCubeDescriptorWithPixelFormatSizeMipmapped(pixelFormat PixelFormat, size uint, mipmapped bool) TextureDescriptor {
	rv := objc.Call[TextureDescriptor](tc, objc.Sel("textureCubeDescriptorWithPixelFormat:size:mipmapped:"), pixelFormat, size, mipmapped)
	return rv
}

// Creates a texture descriptor object for a cube texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516090-texturecubedescriptorwithpixelfo?language=objc
func TextureDescriptor_TextureCubeDescriptorWithPixelFormatSizeMipmapped(pixelFormat PixelFormat, size uint, mipmapped bool) TextureDescriptor {
	return TextureDescriptorClass.TextureCubeDescriptorWithPixelFormatSizeMipmapped(pixelFormat, size, mipmapped)
}

// Creates a texture descriptor object for a 2D texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515511-texture2ddescriptorwithpixelform?language=objc
func (tc _TextureDescriptorClass) Texture2DDescriptorWithPixelFormatWidthHeightMipmapped(pixelFormat PixelFormat, width uint, height uint, mipmapped bool) TextureDescriptor {
	rv := objc.Call[TextureDescriptor](tc, objc.Sel("texture2DDescriptorWithPixelFormat:width:height:mipmapped:"), pixelFormat, width, height, mipmapped)
	return rv
}

// Creates a texture descriptor object for a 2D texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515511-texture2ddescriptorwithpixelform?language=objc
func TextureDescriptor_Texture2DDescriptorWithPixelFormatWidthHeightMipmapped(pixelFormat PixelFormat, width uint, height uint, mipmapped bool) TextureDescriptor {
	return TextureDescriptorClass.Texture2DDescriptorWithPixelFormatWidthHeightMipmapped(pixelFormat, width, height, mipmapped)
}

// Creates a texture descriptor object for a texture buffer. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/2966642-texturebufferdescriptorwithpixel?language=objc
func (tc _TextureDescriptorClass) TextureBufferDescriptorWithPixelFormatWidthResourceOptionsUsage(pixelFormat PixelFormat, width uint, resourceOptions ResourceOptions, usage TextureUsage) TextureDescriptor {
	rv := objc.Call[TextureDescriptor](tc, objc.Sel("textureBufferDescriptorWithPixelFormat:width:resourceOptions:usage:"), pixelFormat, width, resourceOptions, usage)
	return rv
}

// Creates a texture descriptor object for a texture buffer. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/2966642-texturebufferdescriptorwithpixel?language=objc
func TextureDescriptor_TextureBufferDescriptorWithPixelFormatWidthResourceOptionsUsage(pixelFormat PixelFormat, width uint, resourceOptions ResourceOptions, usage TextureUsage) TextureDescriptor {
	return TextureDescriptorClass.TextureBufferDescriptorWithPixelFormatWidthResourceOptionsUsage(pixelFormat, width, resourceOptions, usage)
}

// The dimension and arrangement of texture image data. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516228-texturetype?language=objc
func (t_ TextureDescriptor) TextureType() TextureType {
	rv := objc.Call[TextureType](t_, objc.Sel("textureType"))
	return rv
}

// The dimension and arrangement of texture image data. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516228-texturetype?language=objc
func (t_ TextureDescriptor) SetTextureType(value TextureType) {
	objc.Call[objc.Void](t_, objc.Sel("setTextureType:"), value)
}

// The number of mipmap levels for this texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516300-mipmaplevelcount?language=objc
func (t_ TextureDescriptor) MipmapLevelCount() uint {
	rv := objc.Call[uint](t_, objc.Sel("mipmapLevelCount"))
	return rv
}

// The number of mipmap levels for this texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516300-mipmaplevelcount?language=objc
func (t_ TextureDescriptor) SetMipmapLevelCount(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setMipmapLevelCount:"), value)
}

// The behavior of a new memory allocation. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515776-resourceoptions?language=objc
func (t_ TextureDescriptor) ResourceOptions() ResourceOptions {
	rv := objc.Call[ResourceOptions](t_, objc.Sel("resourceOptions"))
	return rv
}

// The behavior of a new memory allocation. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515776-resourceoptions?language=objc
func (t_ TextureDescriptor) SetResourceOptions(value ResourceOptions) {
	objc.Call[objc.Void](t_, objc.Sel("setResourceOptions:"), value)
}

// The size and bit layout of all pixels in the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515450-pixelformat?language=objc
func (t_ TextureDescriptor) PixelFormat() PixelFormat {
	rv := objc.Call[PixelFormat](t_, objc.Sel("pixelFormat"))
	return rv
}

// The size and bit layout of all pixels in the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515450-pixelformat?language=objc
func (t_ TextureDescriptor) SetPixelFormat(value PixelFormat) {
	objc.Call[objc.Void](t_, objc.Sel("setPixelFormat:"), value)
}

// The pattern you want the GPU to apply to pixels when you read or sample pixels from the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3114305-swizzle?language=objc
func (t_ TextureDescriptor) Swizzle() TextureSwizzleChannels {
	rv := objc.Call[TextureSwizzleChannels](t_, objc.Sel("swizzle"))
	return rv
}

// The pattern you want the GPU to apply to pixels when you read or sample pixels from the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3114305-swizzle?language=objc
func (t_ TextureDescriptor) SetSwizzle(value TextureSwizzleChannels) {
	objc.Call[objc.Void](t_, objc.Sel("setSwizzle:"), value)
}

// The texture's hazard tracking mode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3131697-hazardtrackingmode?language=objc
func (t_ TextureDescriptor) HazardTrackingMode() HazardTrackingMode {
	rv := objc.Call[HazardTrackingMode](t_, objc.Sel("hazardTrackingMode"))
	return rv
}

// The texture's hazard tracking mode. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3131697-hazardtrackingmode?language=objc
func (t_ TextureDescriptor) SetHazardTrackingMode(value HazardTrackingMode) {
	objc.Call[objc.Void](t_, objc.Sel("setHazardTrackingMode:"), value)
}

// The number of array elements for this texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515331-arraylength?language=objc
func (t_ TextureDescriptor) ArrayLength() uint {
	rv := objc.Call[uint](t_, objc.Sel("arrayLength"))
	return rv
}

// The number of array elements for this texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515331-arraylength?language=objc
func (t_ TextureDescriptor) SetArrayLength(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setArrayLength:"), value)
}

// A Boolean value indicating whether the GPU is allowed to adjust the texture's contents to improve GPU performance. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/2966641-allowgpuoptimizedcontents?language=objc
func (t_ TextureDescriptor) AllowGPUOptimizedContents() bool {
	rv := objc.Call[bool](t_, objc.Sel("allowGPUOptimizedContents"))
	return rv
}

// A Boolean value indicating whether the GPU is allowed to adjust the texture's contents to improve GPU performance. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/2966641-allowgpuoptimizedcontents?language=objc
func (t_ TextureDescriptor) SetAllowGPUOptimizedContents(value bool) {
	objc.Call[objc.Void](t_, objc.Sel("setAllowGPUOptimizedContents:"), value)
}

// The height of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516000-height?language=objc
func (t_ TextureDescriptor) Height() uint {
	rv := objc.Call[uint](t_, objc.Sel("height"))
	return rv
}

// The height of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516000-height?language=objc
func (t_ TextureDescriptor) SetHeight(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setHeight:"), value)
}

// The location and access permissions of the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516262-storagemode?language=objc
func (t_ TextureDescriptor) StorageMode() StorageMode {
	rv := objc.Call[StorageMode](t_, objc.Sel("storageMode"))
	return rv
}

// The location and access permissions of the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516262-storagemode?language=objc
func (t_ TextureDescriptor) SetStorageMode(value StorageMode) {
	objc.Call[objc.Void](t_, objc.Sel("setStorageMode:"), value)
}

// The width of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515649-width?language=objc
func (t_ TextureDescriptor) Width() uint {
	rv := objc.Call[uint](t_, objc.Sel("width"))
	return rv
}

// The width of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515649-width?language=objc
func (t_ TextureDescriptor) SetWidth(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setWidth:"), value)
}

// The CPU cache mode used for the CPU mapping of the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515375-cpucachemode?language=objc
func (t_ TextureDescriptor) CpuCacheMode() CPUCacheMode {
	rv := objc.Call[CPUCacheMode](t_, objc.Sel("cpuCacheMode"))
	return rv
}

// The CPU cache mode used for the CPU mapping of the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515375-cpucachemode?language=objc
func (t_ TextureDescriptor) SetCpuCacheMode(value CPUCacheMode) {
	objc.Call[objc.Void](t_, objc.Sel("setCpuCacheMode:"), value)
}

// Options that determine how you can use the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515783-usage?language=objc
func (t_ TextureDescriptor) Usage() TextureUsage {
	rv := objc.Call[TextureUsage](t_, objc.Sel("usage"))
	return rv
}

// Options that determine how you can use the texture. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1515783-usage?language=objc
func (t_ TextureDescriptor) SetUsage(value TextureUsage) {
	objc.Call[objc.Void](t_, objc.Sel("setUsage:"), value)
}

// The depth of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516298-depth?language=objc
func (t_ TextureDescriptor) Depth() uint {
	rv := objc.Call[uint](t_, objc.Sel("depth"))
	return rv
}

// The depth of the texture image for the base level mipmap, in pixels. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516298-depth?language=objc
func (t_ TextureDescriptor) SetDepth(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setDepth:"), value)
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3763055-compressiontype?language=objc
func (t_ TextureDescriptor) CompressionType() TextureCompressionType {
	rv := objc.Call[TextureCompressionType](t_, objc.Sel("compressionType"))
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/3763055-compressiontype?language=objc
func (t_ TextureDescriptor) SetCompressionType(value TextureCompressionType) {
	objc.Call[objc.Void](t_, objc.Sel("setCompressionType:"), value)
}

// The number of samples in each fragment. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516260-samplecount?language=objc
func (t_ TextureDescriptor) SampleCount() uint {
	rv := objc.Call[uint](t_, objc.Sel("sampleCount"))
	return rv
}

// The number of samples in each fragment. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/metal/mtltexturedescriptor/1516260-samplecount?language=objc
func (t_ TextureDescriptor) SetSampleCount(value uint) {
	objc.Call[objc.Void](t_, objc.Sel("setSampleCount:"), value)
}