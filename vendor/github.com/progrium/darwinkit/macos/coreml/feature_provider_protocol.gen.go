// Code generated by DarwinKit. DO NOT EDIT.

package coreml

import (
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

// An interface that represents a collection of values for either a model's input or its output. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreml/mlfeatureprovider?language=objc
type PFeatureProvider interface {
	// optional
	FeatureValueForName(featureName string) FeatureValue
	HasFeatureValueForName() bool

	// optional
	FeatureNames() foundation.Set
	HasFeatureNames() bool
}

// ensure impl type implements protocol interface
var _ PFeatureProvider = (*FeatureProviderObject)(nil)

// A concrete type for the [PFeatureProvider] protocol.
type FeatureProviderObject struct {
	objc.Object
}

func (f_ FeatureProviderObject) HasFeatureValueForName() bool {
	return f_.RespondsToSelector(objc.Sel("featureValueForName:"))
}

// Accesses the feature value given the feature's name. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreml/mlfeatureprovider/2879185-featurevalueforname?language=objc
func (f_ FeatureProviderObject) FeatureValueForName(featureName string) FeatureValue {
	rv := objc.Call[FeatureValue](f_, objc.Sel("featureValueForName:"), featureName)
	return rv
}

func (f_ FeatureProviderObject) HasFeatureNames() bool {
	return f_.RespondsToSelector(objc.Sel("featureNames"))
}

// The set of valid feature names. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/coreml/mlfeatureprovider/2879184-featurenames?language=objc
func (f_ FeatureProviderObject) FeatureNames() foundation.Set {
	rv := objc.Call[foundation.Set](f_, objc.Sel("featureNames"))
	return rv
}