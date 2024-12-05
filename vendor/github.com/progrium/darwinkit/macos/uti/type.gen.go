// Code generated by DarwinKit. DO NOT EDIT.

package uti

import (
	"unsafe"

	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [Type] class.
var TypeClass = _TypeClass{objc.GetClass("UTType")}

type _TypeClass struct {
	objc.Class
}

// An interface definition for the [Type] class.
type IType interface {
	objc.IObject
	ConformsToType(type_ IType) bool
	IsSubtypeOfType(type_ IType) bool
	IsSupertypeOfType(type_ IType) bool
	Identifier() string
	IsPublicType() bool
	Supertypes() foundation.Set
	Version() foundation.Number
	ReferenceURL() foundation.URL
	Tags() map[string][]string
	PreferredFilenameExtension() string
	PreferredMIMEType() string
	IsDynamic() bool
	LocalizedDescription() string
	IsDeclared() bool
}

// An object that represents a type of data to load, send, or receive. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype?language=objc
type Type struct {
	objc.Object
}

func TypeFrom(ptr unsafe.Pointer) Type {
	return Type{
		Object: objc.ObjectFrom(ptr),
	}
}

func (tc _TypeClass) TypeWithIdentifier(identifier string) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithIdentifier:"), identifier)
	return rv
}

// Creates a type based on an identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548218-typewithidentifier?language=objc
func Type_TypeWithIdentifier(identifier string) Type {
	return TypeClass.TypeWithIdentifier(identifier)
}

func (tc _TypeClass) TypeWithMIMEType(mimeType string) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithMIMEType:"), mimeType)
	return rv
}

// Creates a type based on a MIME type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548219-typewithmimetype?language=objc
func Type_TypeWithMIMEType(mimeType string) Type {
	return TypeClass.TypeWithMIMEType(mimeType)
}

func (tc _TypeClass) TypeWithTagTagClassConformingToType(tag string, tagClass string, supertype IType) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithTag:tagClass:conformingToType:"), tag, tagClass, supertype)
	return rv
}

// Creates a type that represents the specified tag and tag class and which conforms to an existing type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548221-typewithtag?language=objc
func Type_TypeWithTagTagClassConformingToType(tag string, tagClass string, supertype IType) Type {
	return TypeClass.TypeWithTagTagClassConformingToType(tag, tagClass, supertype)
}

func (tc _TypeClass) TypeWithMIMETypeConformingToType(mimeType string, supertype IType) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithMIMEType:conformingToType:"), mimeType, supertype)
	return rv
}

// Creates a type based on a MIME type and a supertype that it conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548220-typewithmimetype?language=objc
func Type_TypeWithMIMETypeConformingToType(mimeType string, supertype IType) Type {
	return TypeClass.TypeWithMIMETypeConformingToType(mimeType, supertype)
}

func (tc _TypeClass) TypeWithFilenameExtensionConformingToType(filenameExtension string, supertype IType) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithFilenameExtension:conformingToType:"), filenameExtension, supertype)
	return rv
}

// Creates a type that represents the specified filename extension and conforms to an existing type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548217-typewithfilenameextension?language=objc
func Type_TypeWithFilenameExtensionConformingToType(filenameExtension string, supertype IType) Type {
	return TypeClass.TypeWithFilenameExtensionConformingToType(filenameExtension, supertype)
}

func (tc _TypeClass) TypeWithFilenameExtension(filenameExtension string) Type {
	rv := objc.Call[Type](tc, objc.Sel("typeWithFilenameExtension:"), filenameExtension)
	return rv
}

// Creates a type that represents the specified filename extension. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548216-typewithfilenameextension?language=objc
func Type_TypeWithFilenameExtension(filenameExtension string) Type {
	return TypeClass.TypeWithFilenameExtension(filenameExtension)
}

func (tc _TypeClass) Alloc() Type {
	rv := objc.Call[Type](tc, objc.Sel("alloc"))
	return rv
}

func (tc _TypeClass) New() Type {
	rv := objc.Call[Type](tc, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewType() Type {
	return TypeClass.New()
}

func (t_ Type) Init() Type {
	rv := objc.Call[Type](t_, objc.Sel("init"))
	return rv
}

// Creates a type your app owns based on an identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600608-exportedtypewithidentifier?language=objc
func (tc _TypeClass) ExportedTypeWithIdentifier(identifier string) Type {
	rv := objc.Call[Type](tc, objc.Sel("exportedTypeWithIdentifier:"), identifier)
	return rv
}

// Creates a type your app owns based on an identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600608-exportedtypewithidentifier?language=objc
func Type_ExportedTypeWithIdentifier(identifier string) Type {
	return TypeClass.ExportedTypeWithIdentifier(identifier)
}

// Returns a Boolean value that indicates whether a type conforms to the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548203-conformstotype?language=objc
func (t_ Type) ConformsToType(type_ IType) bool {
	rv := objc.Call[bool](t_, objc.Sel("conformsToType:"), type_)
	return rv
}

// Creates a type your app uses, but doesn’t own, based on an identifier and a supertype that it conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600611-importedtypewithidentifier?language=objc
func (tc _TypeClass) ImportedTypeWithIdentifierConformingToType(identifier string, parentType IType) Type {
	rv := objc.Call[Type](tc, objc.Sel("importedTypeWithIdentifier:conformingToType:"), identifier, parentType)
	return rv
}

// Creates a type your app uses, but doesn’t own, based on an identifier and a supertype that it conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600611-importedtypewithidentifier?language=objc
func Type_ImportedTypeWithIdentifierConformingToType(identifier string, parentType IType) Type {
	return TypeClass.ImportedTypeWithIdentifierConformingToType(identifier, parentType)
}

// Creates a type your app owns based on an identifier and a supertype that it conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600609-exportedtypewithidentifier?language=objc
func (tc _TypeClass) ExportedTypeWithIdentifierConformingToType(identifier string, parentType IType) Type {
	rv := objc.Call[Type](tc, objc.Sel("exportedTypeWithIdentifier:conformingToType:"), identifier, parentType)
	return rv
}

// Creates a type your app owns based on an identifier and a supertype that it conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600609-exportedtypewithidentifier?language=objc
func Type_ExportedTypeWithIdentifierConformingToType(identifier string, parentType IType) Type {
	return TypeClass.ExportedTypeWithIdentifierConformingToType(identifier, parentType)
}

// Creates a type your app uses, but doesn’t own, based on an identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600610-importedtypewithidentifier?language=objc
func (tc _TypeClass) ImportedTypeWithIdentifier(identifier string) Type {
	rv := objc.Call[Type](tc, objc.Sel("importedTypeWithIdentifier:"), identifier)
	return rv
}

// Creates a type your app uses, but doesn’t own, based on an identifier. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3600610-importedtypewithidentifier?language=objc
func Type_ImportedTypeWithIdentifier(identifier string) Type {
	return TypeClass.ImportedTypeWithIdentifier(identifier)
}

// Returns a Boolean value that indicates whether a type is higher in a hierarchy than the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548207-issubtypeoftype?language=objc
func (t_ Type) IsSubtypeOfType(type_ IType) bool {
	rv := objc.Call[bool](t_, objc.Sel("isSubtypeOfType:"), type_)
	return rv
}

// Returns a Boolean value that indicates whether a type is lower in a hierarchy than the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548208-issupertypeoftype?language=objc
func (t_ Type) IsSupertypeOfType(type_ IType) bool {
	rv := objc.Call[bool](t_, objc.Sel("isSupertypeOfType:"), type_)
	return rv
}

// Returns an array of types from the provided tag and tag class. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548222-typeswithtag?language=objc
func (tc _TypeClass) TypesWithTagTagClassConformingToType(tag string, tagClass string, supertype IType) []Type {
	rv := objc.Call[[]Type](tc, objc.Sel("typesWithTag:tagClass:conformingToType:"), tag, tagClass, supertype)
	return rv
}

// Returns an array of types from the provided tag and tag class. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548222-typeswithtag?language=objc
func Type_TypesWithTagTagClassConformingToType(tag string, tagClass string, supertype IType) []Type {
	return TypeClass.TypesWithTagTagClassConformingToType(tag, tagClass, supertype)
}

// The string that represents the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548206-identifier?language=objc
func (t_ Type) Identifier() string {
	rv := objc.Call[string](t_, objc.Sel("identifier"))
	return rv
}

// A Boolean value that indicates whether the type is in the public domain. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548212-publictype?language=objc
func (t_ Type) IsPublicType() bool {
	rv := objc.Call[bool](t_, objc.Sel("isPublicType"))
	return rv
}

// The set of types the type directly or indirectly conforms to. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548214-supertypes?language=objc
func (t_ Type) Supertypes() foundation.Set {
	rv := objc.Call[foundation.Set](t_, objc.Sel("supertypes"))
	return rv
}

// The type’s version, if available. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548223-version?language=objc
func (t_ Type) Version() foundation.Number {
	rv := objc.Call[foundation.Number](t_, objc.Sel("version"))
	return rv
}

// The reference URL for the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548213-referenceurl?language=objc
func (t_ Type) ReferenceURL() foundation.URL {
	rv := objc.Call[foundation.URL](t_, objc.Sel("referenceURL"))
	return rv
}

// The tag specification dictionary of the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548215-tags?language=objc
func (t_ Type) Tags() map[string][]string {
	rv := objc.Call[map[string][]string](t_, objc.Sel("tags"))
	return rv
}

// The preferred filename extension for the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548210-preferredfilenameextension?language=objc
func (t_ Type) PreferredFilenameExtension() string {
	rv := objc.Call[string](t_, objc.Sel("preferredFilenameExtension"))
	return rv
}

// The preferred MIME type for the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548211-preferredmimetype?language=objc
func (t_ Type) PreferredMIMEType() string {
	rv := objc.Call[string](t_, objc.Sel("preferredMIMEType"))
	return rv
}

// A Boolean value that indicates whether the system generates the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548205-dynamic?language=objc
func (t_ Type) IsDynamic() bool {
	rv := objc.Call[bool](t_, objc.Sel("isDynamic"))
	return rv
}

// A localized description of the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548209-localizeddescription?language=objc
func (t_ Type) LocalizedDescription() string {
	rv := objc.Call[string](t_, objc.Sel("localizedDescription"))
	return rv
}

// A Boolean value that indicates whether the system declares the type. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/uniformtypeidentifiers/uttype/3548204-declared?language=objc
func (t_ Type) IsDeclared() bool {
	rv := objc.Call[bool](t_, objc.Sel("isDeclared"))
	return rv
}