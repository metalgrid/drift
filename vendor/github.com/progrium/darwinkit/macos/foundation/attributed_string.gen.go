// Code generated by DarwinKit. DO NOT EDIT.

package foundation

import (
	"unsafe"

	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [AttributedString] class.
var AttributedStringClass = _AttributedStringClass{objc.GetClass("NSAttributedString")}

type _AttributedStringClass struct {
	objc.Class
}

// An interface definition for the [AttributedString] class.
type IAttributedString interface {
	objc.IObject
	Size() Size
	AttributedStringByInflectingString() AttributedString
	ContainsAttachmentsInRange(range_ Range) bool
	RulerAttributesInRange(range_ Range) map[AttributedStringKey]objc.Object
	RTFDFileWrapperFromRangeDocumentAttributes(range_ Range, dict Dictionary) FileWrapper
	DataFromRangeDocumentAttributesError(range_ Range, dict Dictionary, error unsafe.Pointer) []byte
	AttributesAtIndexEffectiveRange(location uint, range_ RangePointer) map[AttributedStringKey]objc.Object
	RangeOfTextTableAtIndex(table objc.IObject, location uint) Range
	FontAttributesInRange(range_ Range) map[AttributedStringKey]objc.Object
	AttributesAtIndexLongestEffectiveRangeInRange(location uint, range_ RangePointer, rangeLimit Range) map[AttributedStringKey]objc.Object
	LineBreakBeforeIndexWithinRange(location uint, aRange Range) uint
	DocFormatFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte
	LineBreakByHyphenatingBeforeIndexWithinRange(location uint, aRange Range) uint
	EnumerateAttributesInRangeOptionsUsingBlock(enumerationRange Range, opts AttributedStringEnumerationOptions, block func(attrs map[AttributedStringKey]objc.Object, range_ Range, stop *bool))
	AttributeAtIndexEffectiveRange(attrName AttributedStringKey, location uint, range_ RangePointer) objc.Object
	DoubleClickAtIndex(location uint) Range
	RangeOfTextListAtIndex(list objc.IObject, location uint) Range
	FileWrapperFromRangeDocumentAttributesError(range_ Range, dict Dictionary, error unsafe.Pointer) FileWrapper
	BoundingRectWithSizeOptionsContext(size Size, options objc.IObject, context objc.IObject) Rect
	DrawInRect(rect Rect)
	EnumerateAttributeInRangeOptionsUsingBlock(attrName AttributedStringKey, enumerationRange Range, opts AttributedStringEnumerationOptions, block func(value objc.Object, range_ Range, stop *bool))
	DrawAtPoint(point Point)
	IsEqualToAttributedString(other IAttributedString) bool
	NextWordFromIndexForward(location uint, isForward bool) uint
	RangeOfTextBlockAtIndex(block objc.IObject, location uint) Range
	RTFFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte
	DrawWithRectOptionsContext(rect Rect, options objc.IObject, context objc.IObject)
	RTFDFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte
	ItemNumberInTextListAtIndex(list objc.IObject, location uint) int
	AttributedSubstringFromRange(range_ Range) AttributedString
	AttributeAtIndexLongestEffectiveRangeInRange(attrName AttributedStringKey, location uint, range_ RangePointer, rangeLimit Range) objc.Object
	String() string
	Length() uint
}

// A string with associated attributes (such as visual style, hyperlinks, or accessibility data) for portions of its text. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring?language=objc
type AttributedString struct {
	objc.Object
}

func AttributedStringFrom(ptr unsafe.Pointer) AttributedString {
	return AttributedString{
		Object: objc.ObjectFrom(ptr),
	}
}

func (a_ AttributedString) InitWithHTMLBaseURLDocumentAttributes(data []byte, base IURL, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithHTML:baseURL:documentAttributes:"), data, base, dict)
	return rv
}

// Creates an attributed string from the HTML in the specified data object and base URL. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1524624-initwithhtml?language=objc
func NewAttributedStringWithHTMLBaseURLDocumentAttributes(data []byte, base IURL, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithHTMLBaseURLDocumentAttributes(data, base, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithContentsOfMarkdownFileAtURLOptionsBaseURLError(markdownFile IURL, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithContentsOfMarkdownFileAtURL:options:baseURL:error:"), markdownFile, options, baseURL, error)
	return rv
}

// Creates an attributed string from the contents of a specified URL that contains Markdown-formatted data using the provided options. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746872-initwithcontentsofmarkdownfileat?language=objc
func NewAttributedStringWithContentsOfMarkdownFileAtURLOptionsBaseURLError(markdownFile IURL, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithContentsOfMarkdownFileAtURLOptionsBaseURLError(markdownFile, options, baseURL, error)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithStringAttributes(str string, attrs map[AttributedStringKey]objc.IObject) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithString:attributes:"), str, attrs)
	return rv
}

// Creates an attributed string with the specified string and attributes. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1408136-initwithstring?language=objc
func NewAttributedStringWithStringAttributes(str string, attrs map[AttributedStringKey]objc.IObject) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithStringAttributes(str, attrs)
	instance.Autorelease()
	return instance
}

func (ac _AttributedStringClass) LocalizedAttributedStringWithFormat(format IAttributedString, args ...any) AttributedString {
	rv := objc.Call[AttributedString](ac, objc.Sel("localizedAttributedStringWithFormat:"), append([]any{format}, args...)...)
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746877-localizedattributedstringwithfor?language=objc
func AttributedString_LocalizedAttributedStringWithFormat(format IAttributedString, args ...any) AttributedString {
	return AttributedStringClass.LocalizedAttributedStringWithFormat(format, args...)
}

func (a_ AttributedString) InitWithAttributedString(attrStr IAttributedString) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithAttributedString:"), attrStr)
	return rv
}

// Creates an attributed string with the characters and attributes of the specified attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1415342-initwithattributedstring?language=objc
func NewAttributedStringWithAttributedString(attrStr IAttributedString) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithAttributedString(attrStr)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithDataOptionsDocumentAttributesError(data []byte, options Dictionary, dict unsafe.Pointer, error unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithData:options:documentAttributes:error:"), data, options, dict, error)
	return rv
}

// Creates an attributed string from the data in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1524613-initwithdata?language=objc
func NewAttributedStringWithDataOptionsDocumentAttributesError(data []byte, options Dictionary, dict unsafe.Pointer, error unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithDataOptionsDocumentAttributesError(data, options, dict, error)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithRTFDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithRTF:documentAttributes:"), data, dict)
	return rv
}

// Creates an attributed string by decoding the stream of RTF commands and data in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1532912-initwithrtf?language=objc
func NewAttributedStringWithRTFDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithRTFDocumentAttributes(data, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithFormatOptionsLocale(format IAttributedString, options AttributedStringFormattingOptions, locale ILocale, args ...any) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithFormat:options:locale:"), append([]any{format, options, locale}, args...)...)
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746873-initwithformat?language=objc
func NewAttributedStringWithFormatOptionsLocale(format IAttributedString, options AttributedStringFormattingOptions, locale ILocale, args ...any) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithFormatOptionsLocale(format, options, locale, args...)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithRTFDFileWrapperDocumentAttributes(wrapper IFileWrapper, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithRTFDFileWrapper:documentAttributes:"), wrapper, dict)
	return rv
}

// Creates an attributed string from the specified file wrapper that contains an RTFD document. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1533594-initwithrtfdfilewrapper?language=objc
func NewAttributedStringWithRTFDFileWrapperDocumentAttributes(wrapper IFileWrapper, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithRTFDFileWrapperDocumentAttributes(wrapper, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithMarkdownOptionsBaseURLError(markdown []byte, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithMarkdown:options:baseURL:error:"), markdown, options, baseURL, error)
	return rv
}

// Creates an attributed string from Markdown-formatted data using the provided options. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746875-initwithmarkdown?language=objc
func NewAttributedStringWithMarkdownOptionsBaseURLError(markdown []byte, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithMarkdownOptionsBaseURLError(markdown, options, baseURL, error)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithHTMLOptionsDocumentAttributes(data []byte, options Dictionary, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithHTML:options:documentAttributes:"), data, options, dict)
	return rv
}

// Creates an attributed string from the HTML in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1535412-initwithhtml?language=objc
func NewAttributedStringWithHTMLOptionsDocumentAttributes(data []byte, options Dictionary, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithHTMLOptionsDocumentAttributes(data, options, dict)
	instance.Autorelease()
	return instance
}

func (ac _AttributedStringClass) LocalizedAttributedStringWithFormatOptions(format IAttributedString, options AttributedStringFormattingOptions, args ...any) AttributedString {
	rv := objc.Call[AttributedString](ac, objc.Sel("localizedAttributedStringWithFormat:options:"), append([]any{format, options}, args...)...)
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746878-localizedattributedstringwithfor?language=objc
func AttributedString_LocalizedAttributedStringWithFormatOptions(format IAttributedString, options AttributedStringFormattingOptions, args ...any) AttributedString {
	return AttributedStringClass.LocalizedAttributedStringWithFormatOptions(format, options, args...)
}

func (a_ AttributedString) InitWithDocFormatDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithDocFormat:documentAttributes:"), data, dict)
	return rv
}

// Creates an attributed string from Microsoft Word format data in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534329-initwithdocformat?language=objc
func NewAttributedStringWithDocFormatDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithDocFormatDocumentAttributes(data, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithMarkdownStringOptionsBaseURLError(markdownString string, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithMarkdownString:options:baseURL:error:"), markdownString, options, baseURL, error)
	return rv
}

// Creates an attributed string from a Markdown-formatted string using the provided options. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746876-initwithmarkdownstring?language=objc
func NewAttributedStringWithMarkdownStringOptionsBaseURLError(markdownString string, options IAttributedStringMarkdownParsingOptions, baseURL IURL, error unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithMarkdownStringOptionsBaseURLError(markdownString, options, baseURL, error)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithURLOptionsDocumentAttributesError(url IURL, options Dictionary, dict unsafe.Pointer, error unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithURL:options:documentAttributes:error:"), url, options, dict, error)
	return rv
}

// Creates an attributed string from the contents of the specified URL. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1530490-initwithurl?language=objc
func NewAttributedStringWithURLOptionsDocumentAttributesError(url IURL, options Dictionary, dict unsafe.Pointer, error unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithURLOptionsDocumentAttributesError(url, options, dict, error)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithHTMLDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithHTML:documentAttributes:"), data, dict)
	return rv
}

// Creates an attributed string from the HTML in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1525953-initwithhtml?language=objc
func NewAttributedStringWithHTMLDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithHTMLDocumentAttributes(data, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithRTFDDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithRTFD:documentAttributes:"), data, dict)
	return rv
}

// Creates an attributed string by decoding the stream of RTFD commands and data in the specified data object. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1530987-initwithrtfd?language=objc
func NewAttributedStringWithRTFDDocumentAttributes(data []byte, dict unsafe.Pointer) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithRTFDDocumentAttributes(data, dict)
	instance.Autorelease()
	return instance
}

func (a_ AttributedString) InitWithString(str string) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("initWithString:"), str)
	return rv
}

// Creates an attributed string with the characters of the specified string and no attribute information. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1407481-initwithstring?language=objc
func NewAttributedStringWithString(str string) AttributedString {
	instance := AttributedStringClass.Alloc().InitWithString(str)
	instance.Autorelease()
	return instance
}

func (ac _AttributedStringClass) Alloc() AttributedString {
	rv := objc.Call[AttributedString](ac, objc.Sel("alloc"))
	return rv
}

func (ac _AttributedStringClass) New() AttributedString {
	rv := objc.Call[AttributedString](ac, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewAttributedString() AttributedString {
	return AttributedStringClass.New()
}

func (a_ AttributedString) Init() AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("init"))
	return rv
}

// Returns the size necessary to draw the string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1528362-size?language=objc
func (a_ AttributedString) Size() Size {
	rv := objc.Call[Size](a_, objc.Sel("size"))
	return rv
}

//	[Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/3746871-attributedstringbyinflectingstri?language=objc
func (a_ AttributedString) AttributedStringByInflectingString() AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("attributedStringByInflectingString"))
	return rv
}

// Returns a Boolean value that indicates if the attributed string contains an attachment in the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1525086-containsattachmentsinrange?language=objc
func (a_ AttributedString) ContainsAttachmentsInRange(range_ Range) bool {
	rv := objc.Call[bool](a_, objc.Sel("containsAttachmentsInRange:"), range_)
	return rv
}

// Returns the ruler (paragraph) attributes in effect for the characters within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1524469-rulerattributesinrange?language=objc
func (a_ AttributedString) RulerAttributesInRange(range_ Range) map[AttributedStringKey]objc.Object {
	rv := objc.Call[map[AttributedStringKey]objc.Object](a_, objc.Sel("rulerAttributesInRange:"), range_)
	return rv
}

// Returns a file wrapper object that contains an RTFD document corresponding to the characters and attributes within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1527502-rtfdfilewrapperfromrange?language=objc
func (a_ AttributedString) RTFDFileWrapperFromRangeDocumentAttributes(range_ Range, dict Dictionary) FileWrapper {
	rv := objc.Call[FileWrapper](a_, objc.Sel("RTFDFileWrapperFromRange:documentAttributes:"), range_, dict)
	return rv
}

// Returns a data object that contains a text stream corresponding to the characters and attributes within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534090-datafromrange?language=objc
func (a_ AttributedString) DataFromRangeDocumentAttributesError(range_ Range, dict Dictionary, error unsafe.Pointer) []byte {
	rv := objc.Call[[]byte](a_, objc.Sel("dataFromRange:documentAttributes:error:"), range_, dict, error)
	return rv
}

// Returns the attributes for the character at the specified index. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1415682-attributesatindex?language=objc
func (a_ AttributedString) AttributesAtIndexEffectiveRange(location uint, range_ RangePointer) map[AttributedStringKey]objc.Object {
	rv := objc.Call[map[AttributedStringKey]objc.Object](a_, objc.Sel("attributesAtIndex:effectiveRange:"), location, range_)
	return rv
}

// Returns the range of the specified text table that contains the specified location. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534365-rangeoftexttable?language=objc
func (a_ AttributedString) RangeOfTextTableAtIndex(table objc.IObject, location uint) Range {
	rv := objc.Call[Range](a_, objc.Sel("rangeOfTextTable:atIndex:"), table, location)
	return rv
}

// Returns the font attributes in effect for the character at the specified location. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1528371-fontattributesinrange?language=objc
func (a_ AttributedString) FontAttributesInRange(range_ Range) map[AttributedStringKey]objc.Object {
	rv := objc.Call[map[AttributedStringKey]objc.Object](a_, objc.Sel("fontAttributesInRange:"), range_)
	return rv
}

// Returns the attributes for the character at the specified index and, by reference, the range where the attributes apply. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1410494-attributesatindex?language=objc
func (a_ AttributedString) AttributesAtIndexLongestEffectiveRangeInRange(location uint, range_ RangePointer, rangeLimit Range) map[AttributedStringKey]objc.Object {
	rv := objc.Call[map[AttributedStringKey]objc.Object](a_, objc.Sel("attributesAtIndex:longestEffectiveRange:inRange:"), location, range_, rangeLimit)
	return rv
}

// Returns the appropriate line break when the character at the index doesn’t fit on the same line as the character at the beginning of the range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1526887-linebreakbeforeindex?language=objc
func (a_ AttributedString) LineBreakBeforeIndexWithinRange(location uint, aRange Range) uint {
	rv := objc.Call[uint](a_, objc.Sel("lineBreakBeforeIndex:withinRange:"), location, aRange)
	return rv
}

// Returns a data object that contains a Microsoft Word–format stream corresponding to the characters and attributes within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1528592-docformatfromrange?language=objc
func (a_ AttributedString) DocFormatFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte {
	rv := objc.Call[[]byte](a_, objc.Sel("docFormatFromRange:documentAttributes:"), range_, dict)
	return rv
}

// Returns the index of the closest character before the specified index, and within the specified range, that can fit on a new line by hyphenating. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534114-linebreakbyhyphenatingbeforeinde?language=objc
func (a_ AttributedString) LineBreakByHyphenatingBeforeIndexWithinRange(location uint, aRange Range) uint {
	rv := objc.Call[uint](a_, objc.Sel("lineBreakByHyphenatingBeforeIndex:withinRange:"), location, aRange)
	return rv
}

// Executes the specified block for each range of attributes in the attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1412070-enumerateattributesinrange?language=objc
func (a_ AttributedString) EnumerateAttributesInRangeOptionsUsingBlock(enumerationRange Range, opts AttributedStringEnumerationOptions, block func(attrs map[AttributedStringKey]objc.Object, range_ Range, stop *bool)) {
	objc.Call[objc.Void](a_, objc.Sel("enumerateAttributesInRange:options:usingBlock:"), enumerationRange, opts, block)
}

// Returns the value for an attribute with the specified name of the character at the specified index and, by reference, the range where the attribute applies. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1408174-attribute?language=objc
func (a_ AttributedString) AttributeAtIndexEffectiveRange(attrName AttributedStringKey, location uint, range_ RangePointer) objc.Object {
	rv := objc.Call[objc.Object](a_, objc.Sel("attribute:atIndex:effectiveRange:"), attrName, location, range_)
	return rv
}

// Returns the range of characters that form a word (or other linguistic unit) surrounding the specified index, taking language characteristics into account. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534748-doubleclickatindex?language=objc
func (a_ AttributedString) DoubleClickAtIndex(location uint) Range {
	rv := objc.Call[Range](a_, objc.Sel("doubleClickAtIndex:"), location)
	return rv
}

// Returns the range of the specified text list that contains the specified location. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1525314-rangeoftextlist?language=objc
func (a_ AttributedString) RangeOfTextListAtIndex(list objc.IObject, location uint) Range {
	rv := objc.Call[Range](a_, objc.Sel("rangeOfTextList:atIndex:"), list, location)
	return rv
}

// Returns a file wrapper object that contains a text stream corresponding to the characters and attributes within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1530461-filewrapperfromrange?language=objc
func (a_ AttributedString) FileWrapperFromRangeDocumentAttributesError(range_ Range, dict Dictionary, error unsafe.Pointer) FileWrapper {
	rv := objc.Call[FileWrapper](a_, objc.Sel("fileWrapperFromRange:documentAttributes:error:"), range_, dict, error)
	return rv
}

// Returns the bounding rectangle necessary to draw the string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1529154-boundingrectwithsize?language=objc
func (a_ AttributedString) BoundingRectWithSizeOptionsContext(size Size, options objc.IObject, context objc.IObject) Rect {
	rv := objc.Call[Rect](a_, objc.Sel("boundingRectWithSize:options:context:"), size, options, context)
	return rv
}

// Draws the attributed string inside the specified bounding rectangle in the current graphics context. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1531631-drawinrect?language=objc
func (a_ AttributedString) DrawInRect(rect Rect) {
	objc.Call[objc.Void](a_, objc.Sel("drawInRect:"), rect)
}

// Executes the specified block for each range of a particular attribute in the attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1412461-enumerateattribute?language=objc
func (a_ AttributedString) EnumerateAttributeInRangeOptionsUsingBlock(attrName AttributedStringKey, enumerationRange Range, opts AttributedStringEnumerationOptions, block func(value objc.Object, range_ Range, stop *bool)) {
	objc.Call[objc.Void](a_, objc.Sel("enumerateAttribute:inRange:options:usingBlock:"), attrName, enumerationRange, opts, block)
}

// Draws the attributed string starting at the specified point in the current graphics context. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1529478-drawatpoint?language=objc
func (a_ AttributedString) DrawAtPoint(point Point) {
	objc.Call[objc.Void](a_, objc.Sel("drawAtPoint:"), point)
}

// Returns a Boolean value that indicates whether the attributed string is equal to another attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1414808-isequaltoattributedstring?language=objc
func (a_ AttributedString) IsEqualToAttributedString(other IAttributedString) bool {
	rv := objc.Call[bool](a_, objc.Sel("isEqualToAttributedString:"), other)
	return rv
}

// Creates an attributed string with an attachment. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1508376-attributedstringwithattachment?language=objc
func (ac _AttributedStringClass) AttributedStringWithAttachment(attachment objc.IObject) AttributedString {
	rv := objc.Call[AttributedString](ac, objc.Sel("attributedStringWithAttachment:"), attachment)
	return rv
}

// Creates an attributed string with an attachment. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1508376-attributedstringwithattachment?language=objc
func AttributedString_AttributedStringWithAttachment(attachment objc.IObject) AttributedString {
	return AttributedStringClass.AttributedStringWithAttachment(attachment)
}

// Returns the index of the first character of the word after or before the specified index. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1535305-nextwordfromindex?language=objc
func (a_ AttributedString) NextWordFromIndexForward(location uint, isForward bool) uint {
	rv := objc.Call[uint](a_, objc.Sel("nextWordFromIndex:forward:"), location, isForward)
	return rv
}

// Returns the range of the individual text block that contains the specified location. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1534045-rangeoftextblock?language=objc
func (a_ AttributedString) RangeOfTextBlockAtIndex(block objc.IObject, location uint) Range {
	rv := objc.Call[Range](a_, objc.Sel("rangeOfTextBlock:atIndex:"), block, location)
	return rv
}

// Returns a data object that contains an RTF stream corresponding to the characters and attributes within the specified range, omitting all attachment attributes. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1535158-rtffromrange?language=objc
func (a_ AttributedString) RTFFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte {
	rv := objc.Call[[]byte](a_, objc.Sel("RTFFromRange:documentAttributes:"), range_, dict)
	return rv
}

// Draws the attributed string in the specified bounding rectangle using the provided options. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1524971-drawwithrect?language=objc
func (a_ AttributedString) DrawWithRectOptionsContext(rect Rect, options objc.IObject, context objc.IObject) {
	objc.Call[objc.Void](a_, objc.Sel("drawWithRect:options:context:"), rect, options, context)
}

// Returns a data object that contains an RTFD stream corresponding to the characters and attributes within the specified range. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1530578-rtfdfromrange?language=objc
func (a_ AttributedString) RTFDFromRangeDocumentAttributes(range_ Range, dict Dictionary) []byte {
	rv := objc.Call[[]byte](a_, objc.Sel("RTFDFromRange:documentAttributes:"), range_, dict)
	return rv
}

// Returns the index of the item at the specified location within the list. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1527303-itemnumberintextlist?language=objc
func (a_ AttributedString) ItemNumberInTextListAtIndex(list objc.IObject, location uint) int {
	rv := objc.Call[int](a_, objc.Sel("itemNumberInTextList:atIndex:"), list, location)
	return rv
}

// Returns an attributed string consisting of the characters and attributes within the specified range in the attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1414283-attributedsubstringfromrange?language=objc
func (a_ AttributedString) AttributedSubstringFromRange(range_ Range) AttributedString {
	rv := objc.Call[AttributedString](a_, objc.Sel("attributedSubstringFromRange:"), range_)
	return rv
}

// Returns the value for the attribute with the specified name of the character at the specified index and, by reference, the range where the attribute applies. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1412701-attribute?language=objc
func (a_ AttributedString) AttributeAtIndexLongestEffectiveRangeInRange(attrName AttributedStringKey, location uint, range_ RangePointer, rangeLimit Range) objc.Object {
	rv := objc.Call[objc.Object](a_, objc.Sel("attribute:atIndex:longestEffectiveRange:inRange:"), attrName, location, range_, rangeLimit)
	return rv
}

// An array of UTI strings that identify the file types that attributed strings support directly. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1528269-textunfilteredtypes?language=objc
func (ac _AttributedStringClass) TextUnfilteredTypes() []string {
	rv := objc.Call[[]string](ac, objc.Sel("textUnfilteredTypes"))
	return rv
}

// An array of UTI strings that identify the file types that attributed strings support directly. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1528269-textunfilteredtypes?language=objc
func AttributedString_TextUnfilteredTypes() []string {
	return AttributedStringClass.TextUnfilteredTypes()
}

// The character contents of the attributed string as a string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1412616-string?language=objc
func (a_ AttributedString) String() string {
	rv := objc.Call[string](a_, objc.Sel("string"))
	return rv
}

// An array of UTI strings that identify the file types that attributed strings support, either directly or through a user-installed filter service. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1535409-texttypes?language=objc
func (ac _AttributedStringClass) TextTypes() []string {
	rv := objc.Call[[]string](ac, objc.Sel("textTypes"))
	return rv
}

// An array of UTI strings that identify the file types that attributed strings support, either directly or through a user-installed filter service. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1535409-texttypes?language=objc
func AttributedString_TextTypes() []string {
	return AttributedStringClass.TextTypes()
}

// The length of the attributed string. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsattributedstring/1418432-length?language=objc
func (a_ AttributedString) Length() uint {
	rv := objc.Call[uint](a_, objc.Sel("length"))
	return rv
}