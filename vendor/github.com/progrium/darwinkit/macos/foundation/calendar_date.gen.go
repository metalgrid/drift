// Code generated by DarwinKit. DO NOT EDIT.

package foundation

import (
	"unsafe"

	"github.com/progrium/darwinkit/objc"
)

// The class instance for the [CalendarDate] class.
var CalendarDateClass = _CalendarDateClass{objc.GetClass("NSCalendarDate")}

type _CalendarDateClass struct {
	objc.Class
}

// An interface definition for the [CalendarDate] class.
type ICalendarDate interface {
	IDate
}

// A specialized date object with embedded calendar information. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nscalendardate?language=objc
type CalendarDate struct {
	Date
}

func CalendarDateFrom(ptr unsafe.Pointer) CalendarDate {
	return CalendarDate{
		Date: DateFrom(ptr),
	}
}

func (cc _CalendarDateClass) Alloc() CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("alloc"))
	return rv
}

func (cc _CalendarDateClass) New() CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("new"))
	rv.Autorelease()
	return rv
}

func NewCalendarDate() CalendarDate {
	return CalendarDateClass.New()
}

func (c_ CalendarDate) Init() CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("init"))
	return rv
}

func (cc _CalendarDateClass) DateWithTimeIntervalSince1970(secs TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("dateWithTimeIntervalSince1970:"), secs)
	return rv
}

// Creates and returns a date object set to the given number of seconds from 00:00:00 UTC on 1 January 1970. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1591576-datewithtimeintervalsince1970?language=objc
func CalendarDate_DateWithTimeIntervalSince1970(secs TimeInterval) CalendarDate {
	return CalendarDateClass.DateWithTimeIntervalSince1970(secs)
}

func (cc _CalendarDateClass) Date() CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("date"))
	return rv
}

// Creates and returns a new date object set to the current date and time. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1591574-date?language=objc
func CalendarDate_Date() CalendarDate {
	return CalendarDateClass.Date()
}

func (c_ CalendarDate) InitWithTimeIntervalSince1970(secs TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("initWithTimeIntervalSince1970:"), secs)
	return rv
}

// Returns a date object initialized relative to 00:00:00 UTC on 1 January 1970 by a given number of seconds. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1416453-initwithtimeintervalsince1970?language=objc
func NewCalendarDateWithTimeIntervalSince1970(secs TimeInterval) CalendarDate {
	instance := CalendarDateClass.Alloc().InitWithTimeIntervalSince1970(secs)
	instance.Autorelease()
	return instance
}

func (cc _CalendarDateClass) DateWithTimeIntervalSinceDate(secsToBeAdded TimeInterval, date IDate) CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("dateWithTimeInterval:sinceDate:"), secsToBeAdded, date)
	return rv
}

// Creates and returns a date object set to a given number of seconds from the specified date. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1591578-datewithtimeinterval?language=objc
func CalendarDate_DateWithTimeIntervalSinceDate(secsToBeAdded TimeInterval, date IDate) CalendarDate {
	return CalendarDateClass.DateWithTimeIntervalSinceDate(secsToBeAdded, date)
}

func (c_ CalendarDate) InitWithTimeIntervalSinceReferenceDate(ti TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("initWithTimeIntervalSinceReferenceDate:"), ti)
	return rv
}

// Returns a date object initialized relative to 00:00:00 UTC on 1 January 2001 by a given number of seconds. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1409769-initwithtimeintervalsincereferen?language=objc
func NewCalendarDateWithTimeIntervalSinceReferenceDate(ti TimeInterval) CalendarDate {
	instance := CalendarDateClass.Alloc().InitWithTimeIntervalSinceReferenceDate(ti)
	instance.Autorelease()
	return instance
}

func (c_ CalendarDate) InitWithTimeIntervalSinceNow(secs TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("initWithTimeIntervalSinceNow:"), secs)
	return rv
}

// Returns a date object initialized relative to the current date and time by a given number of seconds. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1411701-initwithtimeintervalsincenow?language=objc
func NewCalendarDateWithTimeIntervalSinceNow(secs TimeInterval) CalendarDate {
	instance := CalendarDateClass.Alloc().InitWithTimeIntervalSinceNow(secs)
	instance.Autorelease()
	return instance
}

func (c_ CalendarDate) InitWithTimeIntervalSinceDate(secsToBeAdded TimeInterval, date IDate) CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("initWithTimeInterval:sinceDate:"), secsToBeAdded, date)
	return rv
}

// Returns a date object initialized relative to another given date by a given number of seconds. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1414201-initwithtimeinterval?language=objc
func NewCalendarDateWithTimeIntervalSinceDate(secsToBeAdded TimeInterval, date IDate) CalendarDate {
	instance := CalendarDateClass.Alloc().InitWithTimeIntervalSinceDate(secsToBeAdded, date)
	instance.Autorelease()
	return instance
}

func (c_ CalendarDate) DateByAddingTimeInterval(ti TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](c_, objc.Sel("dateByAddingTimeInterval:"), ti)
	return rv
}

// Returns a new date object that is set to a given number of seconds relative to the receiver. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1408823-datebyaddingtimeinterval?language=objc
func CalendarDate_DateByAddingTimeInterval(ti TimeInterval) CalendarDate {
	instance := CalendarDateClass.Alloc().DateByAddingTimeInterval(ti)
	instance.Autorelease()
	return instance
}

func (cc _CalendarDateClass) DateWithTimeIntervalSinceNow(secs TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("dateWithTimeIntervalSinceNow:"), secs)
	return rv
}

// Creates and returns a date object set to a given number of seconds from the current date and time. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1591579-datewithtimeintervalsincenow?language=objc
func CalendarDate_DateWithTimeIntervalSinceNow(secs TimeInterval) CalendarDate {
	return CalendarDateClass.DateWithTimeIntervalSinceNow(secs)
}

func (cc _CalendarDateClass) DateWithTimeIntervalSinceReferenceDate(ti TimeInterval) CalendarDate {
	rv := objc.Call[CalendarDate](cc, objc.Sel("dateWithTimeIntervalSinceReferenceDate:"), ti)
	return rv
}

// Creates and returns a date object set to a given number of seconds from 00:00:00 UTC on 1 January 2001. [Full Topic]
//
// [Full Topic]: https://developer.apple.com/documentation/foundation/nsdate/1591577-datewithtimeintervalsincereferen?language=objc
func CalendarDate_DateWithTimeIntervalSinceReferenceDate(ti TimeInterval) CalendarDate {
	return CalendarDateClass.DateWithTimeIntervalSinceReferenceDate(ti)
}
