// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package walk

import (
	"fmt"
	"math"
	"syscall"
	"unsafe"

	"github.com/tailscale/walk/idalloc"
	"github.com/tailscale/win"
)

type CloseReason byte

const (
	CloseReasonUnknown CloseReason = iota
	CloseReasonUser
)

var (
	taskbarButtonCreatedMsgId uint32

	activeForm *FormBase
)

func init() {
	AppendToWalkInit(func() {
		taskbarButtonCreatedMsgId = win.RegisterWindowMessage(syscall.StringToUTF16Ptr("TaskbarButtonCreated"))
	})
}

type Form interface {
	Container
	AsFormBase() *FormBase
	Starting() *Event
	Closing() *CloseEvent
	Activating() *Event
	Deactivating() *Event
	Activate() error
	Show()
	Hide()
	Title() string
	SetTitle(title string) error
	TitleChanged() *Event
	Icon() Image
	SetIcon(icon Image) error
	IconChanged() *Event
	Owner() Form
	SetOwner(owner Form) error
	ProgressIndicator() *ProgressIndicator

	// RightToLeftLayout returns whether coordinates on the x axis of the
	// Form increase from right to left.
	RightToLeftLayout() bool

	// SetRightToLeftLayout sets whether coordinates on the x axis of the
	// Form increase from right to left.
	SetRightToLeftLayout(rtl bool) error

	// SetFocusToWindow sets keyboard focus to the Window specified by w.
	SetFocusToWindow(w Window)
}

type FormBase struct {
	WindowBase
	clientComposite             *Composite
	owner                       Form
	stopwatch                   *stopwatch
	performLayout               chan layoutStartInfo
	layoutResults               chan layoutResultsWithCompletionFuncs
	inSizeLoop                  chan bool
	updateStopwatch             chan *stopwatch
	quitLayoutPerformer         chan struct{}
	closingPublisher            CloseEventPublisher
	activatingPublisher         EventPublisher
	deactivatingPublisher       EventPublisher
	startingPublisher           EventPublisher
	titleChangedPublisher       EventPublisher
	iconChangedPublisher        EventPublisher
	progressIndicator           *ProgressIndicator
	icon                        Image
	prevFocusHWnd               win.HWND
	proposedSize                Size // in native pixels
	closeReason                 CloseReason
	inSizingLoop                bool
	startingLayoutViaSizingLoop bool
	isInRestoreState            bool
	started                     bool
	layoutScheduled             bool
	layoutCompletionFuncs       []func() // used for scheduled layouts only
	ctrlIDs                     idalloc.IDAllocator
}

const maxPredefinedCtrlID = win.IDCONTINUE

func makeControlIDAllocator() idalloc.IDAllocator {
	alloc := idalloc.New(1 << 16)
	for i := 0; i <= maxPredefinedCtrlID; i++ {
		alloc.Allocate()
	}
	return alloc
}

type ctrlIDAllocator interface {
	allocCtrlID() uint16
	freeCtrlID(id uint16)
}

func (fb *FormBase) allocCtrlID() uint16 {
	id, err := fb.ctrlIDs.Allocate()
	if err != nil {
		panic(err)
	}
	return uint16(id)
}

func (fb *FormBase) freeCtrlID(id uint16) {
	if id > maxPredefinedCtrlID {
		fb.ctrlIDs.Free(uint32(id))
	}
}

func (fb *FormBase) assignCtrlIDs(child *WidgetBase) {
	walkDescendants(child, func(w Window) bool {
		if wgt, ok := w.(Widget); ok {
			wgt.setCtrlID(fb)
		}
		return true
	})
}

func (fb *FormBase) revokeCtrlIDs(child *WidgetBase) {
	walkDescendants(child, func(w Window) bool {
		if wgt, ok := w.(Widget); ok {
			wgt.clearCtrlID(fb)
		}
		return true
	})
}

func (fb *FormBase) init(form Form) error {
	fb.ctrlIDs = makeControlIDAllocator()

	var err error
	if fb.clientComposite, err = NewComposite(form); err != nil {
		return err
	}
	fb.clientComposite.SetName("clientComposite")
	fb.clientComposite.background = nil

	fb.clientComposite.children.observer = form.AsFormBase()

	fb.MustRegisterProperty("Icon", NewProperty(
		func() interface{} {
			return fb.Icon()
		},
		func(v interface{}) error {
			icon, err := IconFrom(v, fb.DPI())
			if err != nil {
				return err
			}

			var img Image
			if icon != nil {
				img = icon
			}

			fb.SetIcon(img)

			return nil
		},
		fb.iconChangedPublisher.Event()))

	fb.MustRegisterProperty("Title", NewProperty(
		func() interface{} {
			return fb.Title()
		},
		func(v interface{}) error {
			return fb.SetTitle(assertStringOr(v, ""))
		},
		fb.titleChangedPublisher.Event()))

	version := win.GetVersion()
	if (version&0xFF) > 6 || ((version&0xFF) == 6 && (version&0xFF00>>8) > 0) {
		win.ChangeWindowMessageFilterEx(fb.hWnd, taskbarButtonCreatedMsgId, win.MSGFLT_ALLOW, nil)
	}

	fb.performLayout, fb.layoutResults, fb.inSizeLoop, fb.updateStopwatch, fb.quitLayoutPerformer = startLayoutPerformer(fb)
	return nil
}

func (fb *FormBase) start() {
	if fb.owner != nil {
		invalidateDescendentBorders := func() {
			walkDescendants(fb.owner, func(wnd Window) bool {
				if widget, ok := wnd.(Widget); ok {
					widget.AsWidgetBase().invalidateBorderInParent()
				}

				return true
			})
		}

		invalidateDescendentBorders()
		defer invalidateDescendentBorders()
	}

	fb.started = true
	fb.startingPublisher.Publish()

	fb.SetBoundsPixels(fb.BoundsPixels())

	if fb.proposedSize == (Size{}) {
		fb.proposedSize = maxSize(SizeFrom96DPI(fb.minSize96dpi, fb.DPI()), fb.SizePixels())
		if !fb.Suspended() {
			fb.startLayout(nil)
		}
	}

	fb.SetSuspended(false)
}

func (fb *FormBase) Dispose() {
	if fb.hWnd != 0 && fb.quitLayoutPerformer != nil {
		fb.quitLayoutPerformer <- struct{}{}
	}

	fb.WindowBase.Dispose()
}

func (fb *FormBase) AsContainerBase() *ContainerBase {
	if fb.clientComposite == nil {
		return nil
	}

	return fb.clientComposite.AsContainerBase()
}

func (fb *FormBase) AsFormBase() *FormBase {
	return fb
}

func (fb *FormBase) Children() *WidgetList {
	if fb.clientComposite == nil {
		return nil
	}

	return fb.clientComposite.Children()
}

func (fb *FormBase) Layout() Layout {
	if fb.clientComposite == nil {
		return nil
	}

	return fb.clientComposite.Layout()
}

func (fb *FormBase) SetLayout(value Layout) error {
	if fb.clientComposite == nil {
		return newError("clientComposite not initialized")
	}

	return fb.clientComposite.SetLayout(value)
}

func (fb *FormBase) SetBoundsPixels(bounds Rectangle) error {
	if layout := fb.Layout(); layout != nil {
		layoutItem := CreateLayoutItemsForContainer(fb)
		minSize := fb.sizeFromClientSizePixels(layoutItem.MinSizeForSize(bounds.Size()))
		minSize = fb.sizeFromClientSizePixels(layoutItem.MinSizeForSize(minSize))

		if bounds.Width < minSize.Width {
			bounds.Width = minSize.Width
		}
		if bounds.Height < minSize.Height {
			bounds.Height = minSize.Height
		}
	}

	if err := fb.WindowBase.SetBoundsPixels(bounds); err != nil {
		return err
	}

	fb.proposedSize = bounds.Size()

	return nil
}

func (fb *FormBase) fixedSize() bool {
	return !fb.hasStyleBits(win.WS_THICKFRAME)
}

func (fb *FormBase) DataBinder() *DataBinder {
	return fb.clientComposite.DataBinder()
}

func (fb *FormBase) SetDataBinder(db *DataBinder) {
	fb.clientComposite.SetDataBinder(db)
}

func (fb *FormBase) SetSuspended(suspended bool) {
	if suspended == fb.suspended {
		return
	}

	fb.suspended = suspended

	if fb.clientComposite != nil {
		fb.clientComposite.SetSuspended(suspended)
	}
}

func (fb *FormBase) MouseDown() *MouseEvent {
	return fb.clientComposite.MouseDown()
}

func (fb *FormBase) MouseMove() *MouseEvent {
	return fb.clientComposite.MouseMove()
}

func (fb *FormBase) MouseUp() *MouseEvent {
	return fb.clientComposite.MouseUp()
}

func (fb *FormBase) onInsertingWidget(index int, widget Widget) error {
	return fb.clientComposite.onInsertingWidget(index, widget)
}

func (fb *FormBase) onInsertedWidget(index int, widget Widget) error {
	return fb.clientComposite.onInsertedWidget(index, widget)
}

func (fb *FormBase) onRemovingWidget(index int, widget Widget) error {
	return fb.clientComposite.onRemovingWidget(index, widget)
}

func (fb *FormBase) onRemovedWidget(index int, widget Widget) error {
	return fb.clientComposite.onRemovedWidget(index, widget)
}

func (fb *FormBase) onClearingWidgets() error {
	return fb.clientComposite.onClearingWidgets()
}

func (fb *FormBase) onClearedWidgets() error {
	return fb.clientComposite.onClearedWidgets()
}

func (fb *FormBase) ContextMenu() *Menu {
	return fb.clientComposite.ContextMenu()
}

func (fb *FormBase) SetContextMenu(contextMenu *Menu) {
	fb.clientComposite.SetContextMenu(contextMenu)
}

func (fb *FormBase) ContextMenuLocation() Point {
	return fb.clientComposite.ContextMenuLocation()
}

func (fb *FormBase) applyEnabled(enabled bool) {
	fb.WindowBase.applyEnabled(enabled)

	fb.clientComposite.applyEnabled(enabled)
}

func (fb *FormBase) applyFont(font *Font) {
	fb.WindowBase.applyFont(font)

	fb.clientComposite.applyFont(font)
}

func (fb *FormBase) ApplySysColors() {
	fb.WindowBase.ApplySysColors()
	fb.clientComposite.ApplySysColors()
}

func (fb *FormBase) Background() Brush {
	return fb.clientComposite.Background()
}

func (fb *FormBase) SetBackground(background Brush) {
	fb.clientComposite.SetBackground(background)
}

func (fb *FormBase) Title() string {
	return fb.text()
}

func (fb *FormBase) SetTitle(value string) error {
	return fb.setText(value)
}

func (fb *FormBase) TitleChanged() *Event {
	return fb.titleChangedPublisher.Event()
}

// RightToLeftLayout returns whether coordinates on the x axis of the
// FormBase increase from right to left.
func (fb *FormBase) RightToLeftLayout() bool {
	return fb.hasExtendedStyleBits(win.WS_EX_LAYOUTRTL)
}

// SetRightToLeftLayout sets whether coordinates on the x axis of the
// FormBase increase from right to left.
func (fb *FormBase) SetRightToLeftLayout(rtl bool) error {
	return fb.ensureExtendedStyleBits(win.WS_EX_LAYOUTRTL, rtl)
}

func (fb *FormBase) HandleKeyDown(msg *win.MSG) bool {
	ret := false

	key, mods := Key(msg.WParam), ModifiersDown()

	// Tabbing
	if key == KeyTab && (mods&ModControl) != 0 {
		doTabbing := func(tw *TabWidget) {
			index := tw.CurrentIndex()
			if (mods & ModShift) != 0 {
				index--
				if index < 0 {
					index = tw.Pages().Len() - 1
				}
			} else {
				index++
				if index >= tw.Pages().Len() {
					index = 0
				}
			}
			tw.SetCurrentIndex(index)
		}

		hwnd := win.GetFocus()

	LOOP:
		for hwnd != 0 {
			window := windowFromHandle(hwnd)

			switch widget := window.(type) {
			case nil:

			case *TabWidget:
				doTabbing(widget)
				return true

			case Widget:

			default:
				break LOOP
			}

			hwnd = win.GetParent(hwnd)
		}

		walkDescendants(fb.window, func(w Window) bool {
			if tw, ok := w.(*TabWidget); ok {
				doTabbing(tw)
				ret = true
				return false
			}
			return true
		})
		if ret {
			return true
		}
	}

	// Shortcut actions
	hwnd := msg.HWnd
	for hwnd != 0 {
		if window := windowFromHandle(hwnd); window != nil {
			wb := window.AsWindowBase()

			if wb.shortcutActions != nil {
				for _, action := range wb.shortcutActions.actions {
					if action.shortcut.Key == key && action.shortcut.Modifiers == mods && action.Visible() && action.Enabled() {
						action.raiseTriggered()
						return true
					}
				}
			}
		}

		hwnd = win.GetParent(hwnd)
	}

	// WebView
	walkDescendants(fb.window, func(w Window) bool {
		if webView, ok := w.(*WebView); ok {
			webViewHWnd := webView.Handle()
			if webViewHWnd == msg.HWnd || win.IsChild(webViewHWnd, msg.HWnd) {
				_ret := webView.translateAccelerator(msg)
				if _ret {
					ret = _ret
				}
			}
		}
		return true
	})
	return ret
}

func (fb *FormBase) Starting() *Event {
	return fb.startingPublisher.Event()
}

func (fb *FormBase) Activating() *Event {
	return fb.activatingPublisher.Event()
}

func (fb *FormBase) Deactivating() *Event {
	return fb.deactivatingPublisher.Event()
}

func (fb *FormBase) Activate() error {
	if hwndPrevActive := win.SetActiveWindow(fb.hWnd); hwndPrevActive == 0 {
		return lastError("SetActiveWindow")
	}

	return nil
}

func (fb *FormBase) Owner() Form {
	return fb.owner
}

func (fb *FormBase) SetOwner(value Form) error {
	fb.owner = value

	var ownerHWnd win.HWND
	if value != nil {
		ownerHWnd = value.Handle()
	}

	win.SetLastError(0)
	if 0 == win.SetWindowLong(
		fb.hWnd,
		win.GWL_HWNDPARENT,
		int32(ownerHWnd)) && win.GetLastError() != 0 {

		return lastError("SetWindowLong")
	}

	return nil
}

func (fb *FormBase) OwnerWindow() Window {
	return fb.owner
}

func (fb *FormBase) Icon() Image {
	return fb.icon
}

func (fb *FormBase) SetIcon(icon Image) error {
	var hIconSmall, hIconBig uintptr

	if icon != nil {
		dpi := fb.DPI()
		size96dpi := icon.Size()

		smallHeight := int(win.GetSystemMetricsForDpi(win.SM_CYSMICON, uint32(dpi)))
		smallDPI := int(math.Round(float64(smallHeight) / float64(size96dpi.Height) * 96.0))
		smallIcon, err := iconCache.Icon(icon, smallDPI)
		if err != nil {
			return err
		}
		hIconSmall = uintptr(smallIcon.handleForDPI(smallDPI))

		bigHeight := int(win.GetSystemMetricsForDpi(win.SM_CYICON, uint32(dpi)))
		bigDPI := int(math.Round(float64(bigHeight) / float64(size96dpi.Height) * 96.0))
		bigIcon, err := iconCache.Icon(icon, bigDPI)
		if err != nil {
			return err
		}
		hIconBig = uintptr(bigIcon.handleForDPI(bigDPI))
	}

	fb.SendMessage(win.WM_SETICON, 0, hIconSmall)
	fb.SendMessage(win.WM_SETICON, 1, hIconBig)

	fb.icon = icon

	fb.iconChangedPublisher.Publish()

	return nil
}

func (fb *FormBase) IconChanged() *Event {
	return fb.iconChangedPublisher.Event()
}

func (fb *FormBase) Hide() {
	fb.window.SetVisible(false)
}

func (fb *FormBase) Show() {
	fb.proposedSize = maxSize(SizeFrom96DPI(fb.minSize96dpi, fb.DPI()), fb.SizePixels())

	if p, ok := fb.window.(Persistable); ok && p.Persistent() && App().Settings() != nil {
		p.RestoreState()
	}

	fb.start()

	fb.window.SetVisible(true)
}

func (fb *FormBase) Close() error {
	fb.SendMessage(win.WM_CLOSE, 0, 0)

	return nil
}

func (fb *FormBase) Persistent() bool {
	return fb.clientComposite.persistent
}

func (fb *FormBase) SetPersistent(value bool) {
	fb.clientComposite.persistent = value
}

func (fb *FormBase) SaveState() error {
	if err := fb.clientComposite.SaveState(); err != nil {
		return err
	}

	var wp win.WINDOWPLACEMENT

	wp.Length = uint32(unsafe.Sizeof(wp))

	if !win.GetWindowPlacement(fb.hWnd, &wp) {
		return lastError("GetWindowPlacement")
	}

	state := fmt.Sprint(
		wp.Flags, wp.ShowCmd,
		wp.PtMinPosition.X, wp.PtMinPosition.Y,
		wp.PtMaxPosition.X, wp.PtMaxPosition.Y,
		wp.RcNormalPosition.Left, wp.RcNormalPosition.Top,
		wp.RcNormalPosition.Right, wp.RcNormalPosition.Bottom)

	if err := fb.WriteState(state); err != nil {
		return err
	}

	return nil
}

func (fb *FormBase) RestoreState() error {
	if fb.isInRestoreState {
		return nil
	}
	fb.isInRestoreState = true
	defer func() {
		fb.isInRestoreState = false
	}()

	state, err := fb.ReadState()
	if err != nil {
		return err
	}
	if state == "" {
		return nil
	}

	var wp win.WINDOWPLACEMENT

	if _, err := fmt.Sscan(state,
		&wp.Flags, &wp.ShowCmd,
		&wp.PtMinPosition.X, &wp.PtMinPosition.Y,
		&wp.PtMaxPosition.X, &wp.PtMaxPosition.Y,
		&wp.RcNormalPosition.Left, &wp.RcNormalPosition.Top,
		&wp.RcNormalPosition.Right, &wp.RcNormalPosition.Bottom); err != nil {
		return err
	}

	wp.Length = uint32(unsafe.Sizeof(wp))

	if layout := fb.Layout(); layout != nil && fb.fixedSize() {
		layoutItem := CreateLayoutItemsForContainer(fb)
		minSize := fb.sizeFromClientSizePixels(layoutItem.MinSize())

		wp.RcNormalPosition.Right = wp.RcNormalPosition.Left + int32(minSize.Width) - 1
		wp.RcNormalPosition.Bottom = wp.RcNormalPosition.Top + int32(minSize.Height) - 1
	}

	if !win.SetWindowPlacement(fb.hWnd, &wp) {
		return lastError("SetWindowPlacement")
	}

	return fb.clientComposite.RestoreState()
}

func (fb *FormBase) Closing() *CloseEvent {
	return fb.closingPublisher.Event()
}

func (fb *FormBase) ProgressIndicator() *ProgressIndicator {
	return fb.progressIndicator
}

func (fb *FormBase) setStopwatch(sw *stopwatch) {
	fb.stopwatch = sw

	fb.updateStopwatch <- sw
}

func (fb *FormBase) startLayout(completionFuncs []func()) bool {
	if fb.performLayout == nil || fb.inSizingLoop && !fb.startingLayoutViaSizingLoop {
		return false
	}

	cs := fb.clientSizeFromSizePixels(fb.proposedSize)
	min := CreateLayoutItemsForContainer(fb.clientComposite).MinSizeForSize(fb.proposedSize)

	if cs.Width < min.Width || cs.Height < min.Height {
		cs = maxSize(cs, min)
		size := fb.sizeFromClientSizePixels(cs)
		fb.SetSizePixels(size)
		fb.Invalidate()
	}

	cbp := fb.window.ClientBoundsPixels()

	fb.clientComposite.SetBoundsPixels(Rectangle{Y: cbp.Y, Width: cs.Width, Height: cs.Height})

	cli := CreateLayoutItemsForContainer(fb)
	cli.Geometry().ClientSize = cs

	fb.performLayout <- layoutStartInfo{item: cli, completionFuncs: completionFuncs}

	return true
}

func (fb *FormBase) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_ACTIVATE:
		switch win.LOWORD(uint32(wParam)) {
		case win.WA_ACTIVE, win.WA_CLICKACTIVE:
			if fb.prevFocusHWnd != 0 {
				win.SetFocus(fb.prevFocusHWnd)
			}

			activeForm = fb

			fb.activatingPublisher.Publish()

		case win.WA_INACTIVE:
			fb.prevFocusHWnd = win.GetFocus()

			activeForm = nil

			fb.deactivatingPublisher.Publish()
		}

		return 0

	case win.WM_CLOSE:
		fb.closeReason = CloseReasonUnknown
		var canceled bool
		fb.closingPublisher.Publish(&canceled, fb.closeReason)
		if !canceled {
			if p, ok := fb.window.(Persistable); ok && p.Persistent() && App().Settings() != nil {
				p.SaveState()
			}
			fb.started = false
			fb.Dispose()
		}
		return 0

	case win.WM_COMMAND:
		return fb.clientComposite.WndProc(hwnd, msg, wParam, lParam)

	case win.WM_GETMINMAXINFO:
		if fb.Suspended() || fb.proposedSize == (Size{}) {
			break
		}

		mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))

		var min Size
		if layout := fb.clientComposite.layout; layout != nil {
			size := fb.clientSizeFromSizePixels(fb.proposedSize)
			layoutItem := CreateLayoutItemsForContainer(fb)
			min = fb.sizeFromClientSizePixels(layoutItem.MinSizeForSize(size))

			if fb.proposedSize.Width < min.Width {
				min = fb.sizeFromClientSizePixels(layoutItem.MinSizeForSize(min))
			}
		}

		minSize := SizeFrom96DPI(fb.minSize96dpi, fb.DPI())

		mmi.PtMinTrackSize = Point{
			maxi(min.Width, minSize.Width),
			maxi(min.Height, minSize.Height),
		}.toPOINT()
		return 0

	case win.WM_NOTIFY:
		return fb.clientComposite.WndProc(hwnd, msg, wParam, lParam)

	case win.WM_SETTEXT:
		fb.titleChangedPublisher.Publish()

	case win.WM_ENTERSIZEMOVE:
		fb.inSizingLoop = true
		fb.inSizeLoop <- true

	case win.WM_EXITSIZEMOVE:
		fb.inSizingLoop = false
		fb.inSizeLoop <- false

	case win.WM_WINDOWPOSCHANGED:
		wp := (*win.WINDOWPOS)(unsafe.Pointer(lParam))

		if wp.Flags&win.SWP_SHOWWINDOW != 0 {
			fb.startLayout(nil)
		}

		if wp.Flags&win.SWP_NOSIZE != 0 || fb.Layout() == nil || fb.Suspended() {
			break
		}

		fb.proposedSize = Size{int(wp.Cx), int(wp.Cy)}

		const performingLayoutSubject = "*FormBase.WndProc - WM_WINDOWPOSCHANGED - full layout from sizing loop"

		if fb.inSizingLoop {
			fb.startingLayoutViaSizingLoop = true

			if fb.stopwatch != nil {
				fb.stopwatch.Start(performingLayoutSubject)
			}
		}

		if fb.startLayout(nil) {
			if fb.inSizingLoop {
				fb.startingLayoutViaSizingLoop = false

				results := <-fb.layoutResults
				applyLayoutResults(results.results, fb.stopwatch)

				if fb.stopwatch != nil {
					fb.stopwatch.Stop(performingLayoutSubject)
				}
			}
		}

	case win.WM_SYSCOLORCHANGE:
		fb.ApplySysColors()

	case win.WM_DPICHANGED:
		wasSuspended := fb.Suspended()
		fb.SetSuspended(true)
		defer fb.SetSuspended(wasSuspended)

		dpi := int(win.HIWORD(uint32(wParam)))

		seenInApplyFontToDescendantsDuringDPIChange = make(map[*WindowBase]bool)
		seenInApplyDPIToDescendantsDuringDPIChange = make(map[*WindowBase]bool)
		defer func() {
			seenInApplyFontToDescendantsDuringDPIChange = nil
			seenInApplyDPIToDescendantsDuringDPIChange = nil
		}()

		fb.clientComposite.ApplyDPI(dpi)
		fb.ApplyDPI(dpi)
		if fb.progressIndicator != nil {
			fb.progressIndicator.SetOverlayIcon(fb.progressIndicator.overlayIcon, fb.progressIndicator.overlayIconDescription)
		}
		applyDPIToDescendants(fb.window, dpi)

		fb.SetSuspended(wasSuspended)

		rc := (*win.RECT)(unsafe.Pointer(lParam))
		bounds := rectangleFromRECT(*rc)
		fb.proposedSize = bounds.Size()
		fb.window.SetBoundsPixels(bounds)

		fb.SetIcon(fb.icon)

	case win.WM_SYSCOMMAND:
		if wParam == win.SC_CLOSE {
			fb.closeReason = CloseReasonUser
		}

	case taskbarButtonCreatedMsgId:
		version := win.GetVersion()
		major := version & 0xFF
		minor := version & 0xFF00 >> 8
		// Check that the OS is Win 7 or later (Win 7 is v6.1).
		if fb.progressIndicator == nil && (major > 6 || (major == 6 && minor > 0)) {
			fb.progressIndicator, _ = newTaskbarList3(fb.hWnd)
		}
	}

	return fb.WindowBase.WndProc(hwnd, msg, wParam, lParam)
}

func (fb *FormBase) EnterMode() {
}

func (fb *FormBase) Running() bool {
	return fb.started
}

func (fb *FormBase) OnPreTranslate(msg *win.MSG) bool {
	return DefaultModalPreTranslate(fb, msg)
}

func (fb *FormBase) OnPostDispatch() {
	if !fb.layoutScheduled {
		return
	}

	fb.layoutScheduled = false
	completionFuncs := fb.layoutCompletionFuncs
	fb.layoutCompletionFuncs = nil
	fb.startLayout(completionFuncs)
}

func (fb *FormBase) Window() Window {
	return fb.window
}

func (fb *FormBase) SetFocusToWindow(w Window) {
	if w != nil {
		win.SetFocus(w.Handle())
	}
}