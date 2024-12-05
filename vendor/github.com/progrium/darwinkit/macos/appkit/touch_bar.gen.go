// Code generated by DarwinKit. DO NOT EDIT.

package appkit

import (
	"unsafe"

	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [TouchBar] class.
var TouchBarClass = _TouchBarClass{objc.GetClass("NSTouchBar")}

type _TouchBarClass struct {
	objc.Class
}

// An interface definition for the [TouchBar] class.
type ITouchBar interface {
	objc.IObject
	ItemForIdentifier(identifier TouchBarItemIdentifier) TouchBarItem
	TemplateItems() foundation.Set
	SetTemplateItems(value foundation.ISet)
	CustomizationRequiredItemIdentifiers() []TouchBarItemIdentifier
	SetCustomizationRequiredItemIdentifiers(value []TouchBarItemIdentifier)
	PrincipalItemIdentifier() TouchBarItemIdentifier
	SetPrincipalItemIdentifier(value TouchBarItemIdentifier)
	CustomizationAllowedItemIdentifiers() []TouchBarItemIdentifier
	SetCustomizationAllowedItemIdentifiers(value []TouchBarItemIdentifier)
	CustomizationIdentifier() TouchBarCustomizationIdentifier
	SetCustomizationIdentifier(value TouchBarCustomizationIdentifier)
	DefaultItemIdentifiers() []TouchBarItemIdentifier
	SetDefaultItemIdentifiers(value []TouchBarItemIdentifier)
	ItemIdentifiers() []TouchBarItemIdentifier
	IsVisible() bool
	EscapeKeyReplacementItemIdentifier() TouchBarItemIdentifier
	SetEscapeKeyReplacementItemIdentifier(value TouchBarItemIdentifier)
	Delegate() TouchBarDelegateObject
	SetDelegate(value PTouchBarDelegate)
	SetDelegateObject(valueObject objc.IObject)
}

// An object that provides dynamic contextual controls in the Touch Bar of supported models of MacBook Pro. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar?language=objc
type TouchBar struct {
	objc.Object
}

func TouchBarFrom(ptr unsafe.Pointer) TouchBar {
	return TouchBar{
		Object: objc.ObjectFrom(ptr),
	}
}

func (t_ TouchBar) Init() TouchBar {
	rv := objc.Call[TouchBar](t_, objc.Sel("init"))
	return rv
}

func (tc _TouchBarClass) Alloc() TouchBar {
	rv := objc.Call[TouchBar](tc, objc.Sel("alloc"))
	return rv
}

func (tc _TouchBarClass) New() TouchBar {
	rv := objc.Call[TouchBar](tc, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewTouchBar() TouchBar {
	return TouchBarClass.New()
}

// Returns the Touch Bar item that corresponds to a given identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544806-itemforidentifier?language=objc
func (t_ TouchBar) ItemForIdentifier(identifier TouchBarItemIdentifier) TouchBarItem {
	rv := objc.Call[TouchBarItem](t_, objc.Sel("itemForIdentifier:"), identifier)
	return rv
}

// The primary source of items that the Touch Bar uses to fill its private items array, unless you provide items using a delegate. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2646922-templateitems?language=objc
func (t_ TouchBar) TemplateItems() foundation.Set {
	rv := objc.Call[foundation.Set](t_, objc.Sel("templateItems"))
	return rv
}

// The primary source of items that the Touch Bar uses to fill its private items array, unless you provide items using a delegate. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2646922-templateitems?language=objc
func (t_ TouchBar) SetTemplateItems(value foundation.ISet) {
	objc.Call[objc.Void](t_, objc.Sel("setTemplateItems:"), value)
}

// An optional list of identifiers for items you want to always appear in the Touch Bar and which the user can’t remove during customization. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544675-customizationrequireditemidentif?language=objc
func (t_ TouchBar) CustomizationRequiredItemIdentifiers() []TouchBarItemIdentifier {
	rv := objc.Call[[]TouchBarItemIdentifier](t_, objc.Sel("customizationRequiredItemIdentifiers"))
	return rv
}

// An optional list of identifiers for items you want to always appear in the Touch Bar and which the user can’t remove during customization. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544675-customizationrequireditemidentif?language=objc
func (t_ TouchBar) SetCustomizationRequiredItemIdentifiers(value []TouchBarItemIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setCustomizationRequiredItemIdentifiers:"), value)
}

// The identifier of an item you want the system to center in the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544871-principalitemidentifier?language=objc
func (t_ TouchBar) PrincipalItemIdentifier() TouchBarItemIdentifier {
	rv := objc.Call[TouchBarItemIdentifier](t_, objc.Sel("principalItemIdentifier"))
	return rv
}

// The identifier of an item you want the system to center in the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544871-principalitemidentifier?language=objc
func (t_ TouchBar) SetPrincipalItemIdentifier(value TouchBarItemIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setPrincipalItemIdentifier:"), value)
}

// A list of identifiers for items to show in the Touch Bar’s customization UI. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544796-customizationalloweditemidentifi?language=objc
func (t_ TouchBar) CustomizationAllowedItemIdentifiers() []TouchBarItemIdentifier {
	rv := objc.Call[[]TouchBarItemIdentifier](t_, objc.Sel("customizationAllowedItemIdentifiers"))
	return rv
}

// A list of identifiers for items to show in the Touch Bar’s customization UI. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544796-customizationalloweditemidentifi?language=objc
func (t_ TouchBar) SetCustomizationAllowedItemIdentifiers(value []TouchBarItemIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setCustomizationAllowedItemIdentifiers:"), value)
}

// A globally unique string that makes the Touch Bar eligible for user customization. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544730-customizationidentifier?language=objc
func (t_ TouchBar) CustomizationIdentifier() TouchBarCustomizationIdentifier {
	rv := objc.Call[TouchBarCustomizationIdentifier](t_, objc.Sel("customizationIdentifier"))
	return rv
}

// A globally unique string that makes the Touch Bar eligible for user customization. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544730-customizationidentifier?language=objc
func (t_ TouchBar) SetCustomizationIdentifier(value TouchBarCustomizationIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setCustomizationIdentifier:"), value)
}

// A required list of identifiers for items that you want to appear in the Touch Bar after instantiating it. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2646921-defaultitemidentifiers?language=objc
func (t_ TouchBar) DefaultItemIdentifiers() []TouchBarItemIdentifier {
	rv := objc.Call[[]TouchBarItemIdentifier](t_, objc.Sel("defaultItemIdentifiers"))
	return rv
}

// A required list of identifiers for items that you want to appear in the Touch Bar after instantiating it. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2646921-defaultitemidentifiers?language=objc
func (t_ TouchBar) SetDefaultItemIdentifiers(value []TouchBarItemIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setDefaultItemIdentifiers:"), value)
}

// The list of identifiers for the current items in the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544762-itemidentifiers?language=objc
func (t_ TouchBar) ItemIdentifiers() []TouchBarItemIdentifier {
	rv := objc.Call[[]TouchBarItemIdentifier](t_, objc.Sel("itemIdentifiers"))
	return rv
}

// A Boolean value that Indicates whether the Touch Bar is eligible for display. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544756-visible?language=objc
func (t_ TouchBar) IsVisible() bool {
	rv := objc.Call[bool](t_, objc.Sel("isVisible"))
	return rv
}

// The identifier of an item that replaces the system-provided button in the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2787731-escapekeyreplacementitemidentifi?language=objc
func (t_ TouchBar) EscapeKeyReplacementItemIdentifier() TouchBarItemIdentifier {
	rv := objc.Call[TouchBarItemIdentifier](t_, objc.Sel("escapeKeyReplacementItemIdentifier"))
	return rv
}

// The identifier of an item that replaces the system-provided button in the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2787731-escapekeyreplacementitemidentifi?language=objc
func (t_ TouchBar) SetEscapeKeyReplacementItemIdentifier(value TouchBarItemIdentifier) {
	objc.Call[objc.Void](t_, objc.Sel("setEscapeKeyReplacementItemIdentifier:"), value)
}

// The delegate that provides items to the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544666-delegate?language=objc
func (t_ TouchBar) Delegate() TouchBarDelegateObject {
	rv := objc.Call[TouchBarDelegateObject](t_, objc.Sel("delegate"))
	return rv
}

// The delegate that provides items to the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544666-delegate?language=objc
func (t_ TouchBar) SetDelegate(value PTouchBarDelegate) {
	po0 := objc.WrapAsProtocol("NSTouchBarDelegate", value)
	objc.SetAssociatedObject(t_, objc.AssociationKey("setDelegate"), po0, objc.ASSOCIATION_RETAIN)
	objc.Call[objc.Void](t_, objc.Sel("setDelegate:"), po0)
}

// The delegate that provides items to the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/2544666-delegate?language=objc
func (t_ TouchBar) SetDelegateObject(valueObject objc.IObject) {
	objc.Call[objc.Void](t_, objc.Sel("setDelegate:"), valueObject)
}

// A Boolean value indicating whether the main menu contains an item for customizing the contents of the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/3228044-automaticcustomizetouchbarmenuit?language=objc
func (tc _TouchBarClass) AutomaticCustomizeTouchBarMenuItemEnabled() bool {
	rv := objc.Call[bool](tc, objc.Sel("automaticCustomizeTouchBarMenuItemEnabled"))
	return rv
}

// A Boolean value indicating whether the main menu contains an item for customizing the contents of the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/3228044-automaticcustomizetouchbarmenuit?language=objc
func TouchBar_AutomaticCustomizeTouchBarMenuItemEnabled() bool {
	return TouchBarClass.AutomaticCustomizeTouchBarMenuItemEnabled()
}

// A Boolean value indicating whether the main menu contains an item for customizing the contents of the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/3228044-automaticcustomizetouchbarmenuit?language=objc
func (tc _TouchBarClass) SetAutomaticCustomizeTouchBarMenuItemEnabled(value bool) {
	objc.Call[objc.Void](tc, objc.Sel("setAutomaticCustomizeTouchBarMenuItemEnabled:"), value)
}

// A Boolean value indicating whether the main menu contains an item for customizing the contents of the Touch Bar. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/appkit/nstouchbar/3228044-automaticcustomizetouchbarmenuit?language=objc
func TouchBar_SetAutomaticCustomizeTouchBarMenuItemEnabled(value bool) {
	TouchBarClass.SetAutomaticCustomizeTouchBarMenuItemEnabled(value)
}
