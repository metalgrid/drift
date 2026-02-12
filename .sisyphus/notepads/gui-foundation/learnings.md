## Observer Pattern Implementation (2026-02-12)

### Pattern Used
- Simple `[]func()` callback slice on Peers struct
- No arguments passed to callbacks (callers use `All()` to get current state)
- Post-construction registration via `OnChange(fn func())`

### Concurrency Safety
- Copy observers slice before unlocking mutex
- Call observers WITHOUT holding lock to prevent deadlocks
- Pattern: `observers := p.observers; p.mu.Unlock(); for _, obs := range observers { obs() }`
- Critical: If observer calls back into Peers methods, it won't deadlock

### Test Coverage
- Observer fires on add
- Observer fires on remove  
- Multiple observers all fire
- All() returns snapshot (baseline test)

### Code Location
- `internal/zeroconf/zeroconf.go:58-62` — Peers struct with observers field
- `internal/zeroconf/zeroconf.go:68-72` — OnChange registration method
- `internal/zeroconf/zeroconf.go:113-121` — add() with observer notification
- `internal/zeroconf/zeroconf.go:123-131` — remove() with observer notification
- `internal/zeroconf/zeroconf_test.go` — 4 test cases (all passing)

## [2026-02-12] Task: 5 - Progress Tracking Wrappers

### Implementation Pattern
- Standard Go io.Writer/Reader wrapper pattern
- `ProgressFunc` type: `func(bytesTransferred int64, totalBytes int64)`
- Both wrappers track cumulative bytes and call callback after each operation
- Nil callback support: wrappers work without progress tracking (no-op)

### Key Design Decisions
- Made ProgressFunc parameter optional (nil = no tracking)
- ProgressWriter wraps io.Writer, calls callback after Write()
- ProgressReader wraps io.Reader, calls callback after Read()
- Both maintain cumulative bytesTransferred counter
- totalBytes passed to callback for progress percentage calculation

### Integration Points
- `sendFile()`: Gets file size via Stat(), wraps writer with ProgressWriter
- `storeFile()`: Wraps io.LimitReader with ProgressReader
- Both functions accept optional ProgressFunc parameter (nil in current calls)

### Test Coverage
- 7 tests total (all passing)
- TestProgressWriterReportsBytesCorrectly: Verifies cumulative byte tracking
- TestProgressReaderReportsBytesCorrectly: Verifies cumulative byte tracking
- TestProgressWriterNilCallback: Graceful handling of nil callback
- TestProgressReaderNilCallback: Graceful handling of nil callback
- TestProgressMonotonicallyIncreasing: Bytes never decrease
- TestProgressWriterTotalBytesParameter: totalBytes passed correctly
- TestProgressReaderTotalBytesParameter: totalBytes passed correctly

### Code Locations
- `internal/transport/progress.go` — ProgressWriter/Reader implementations
- `internal/transport/progress_test.go` — 7 test cases
- `internal/transport/connection.go:74-98` — storeFile with ProgressReader
- `internal/transport/connection.go:100-113` — sendFile with ProgressWriter
- `internal/transport/connection.go:54` — storeFile call with nil progress
- `internal/transport/connection.go:63` — sendFile call with nil progress

### Verification
- `go test ./internal/transport/ -v -count=1 -mod=mod` → ALL PASS (21 tests)
- `go vet ./internal/transport/ -mod=mod` → exit 0
- `git status --short` → 4 files (progress.go, progress_test.go, connection.go, connection_test.go)

## [2026-02-12] Task: 6 - BATCH_OFFER Wire Protocol Extension

### Implementation Summary
Extended wire protocol with BATCH_OFFER message type for multi-file transfers.

### Key Patterns Followed
1. **TDD Approach**: Wrote 6 comprehensive tests before implementation
   - TestBatchOfferMarshal: Wire format validation
   - TestBatchOfferUnmarshalRoundTrip: Serialization integrity
   - TestBatchOfferSingleFile: Edge case (single file batch)
   - TestBatchOfferInvalidCount: Error handling (count mismatch)
   - TestBatchOfferEmptyBatch: Error handling (empty batch)
   - TestExistingOfferStillWorks: Backward compatibility

2. **Wire Format Design**:
   - Format: `BATCH_OFFER|count|filename1|mimetype1|size1|filename2|mimetype2|size2\n`
   - Count field enables validation before parsing entries
   - Consistent with existing pipe-delimited message format

3. **Critical Ordering**: BATCH_OFFER case MUST precede OFFER case in UnmarshalMessage
   - Reason: Prefix matching - "BATCH_OFFER" starts with "OFFER"
   - Prevents misclassification of batch messages as single offers

4. **Error Handling**:
   - Validate count matches actual entries (prevents partial parsing)
   - Reject empty batches (count <= 0)
   - Return errors for malformed messages (maintains protocol integrity)

5. **Batch Receive Pattern**:
   - All-or-nothing acceptance (single Ask() for entire batch)
   - Sequential file storage (loop over Files, call storeFile for each)
   - Aggregate size display for user decision

### Code Structure
- **Types**: FileEntry (reusable), BatchOffer (contains []FileEntry)
- **Functions**: MakeBatchOffer (validates files, builds offer), SendBatch (sends batch offer)
- **Handler**: BatchOffer case in HandleConnection (before Offer case)

### Test Results
All 17 tests pass (11 existing + 6 new batch tests)

### Files Modified
- internal/transport/message.go (types, marshal, unmarshal, MakeBatchOffer)
- internal/transport/connection.go (HandleConnection batch case, SendBatch)
- internal/transport/message_test.go (6 new tests)

## [2026-02-12] Task: 7 - Gateway Interface Evolution for Batch Transfers

### Changes Implemented
- Added `FileInfo` struct and `BatchGateway` interface to `internal/platform/gateway.go`
- Evolved `Request` struct from `File string` to `Files []string` for multi-file support
- Updated all platform gateway implementations (Linux/DBus, Windows, macOS) to use Files slice
- Added type assertion pattern in `connection.go` to detect BatchGateway capability
- Updated outbound processors in both `run.go` and `main.go` to handle single vs batch transfers

### Key Patterns
- **Backwards Compatibility**: Kept `Ask(string) string` method; BatchGateway extends Gateway
- **Type Assertion Pattern**: `if bg, ok := gw.(BatchGateway); ok { ... }` for optional interface
- **Slice Wrapping**: Single file requests wrap as `[]string{file}` for uniform handling
- **Conditional Dispatch**: Check `len(request.Files)` to route to SendFile vs SendBatch

### Platform-Specific Notes
- **Linux (DBus)**: Uses dbus_gateway.go, requires godbus/dbus/v5 in vendor
- **Windows**: Uses walk library for GUI, gateway_windows.go
- **macOS**: Uses darwinkit, gateway_macosx.go (build constraints exclude some packages on Linux)
- All platforms now support Files slice in Request struct

### Build Verification
- `go build ./cmd/drift` succeeds on Linux and Windows
- Darwin build shows expected vendor constraint warnings (platform-specific)
- IPv6 format warnings in vet are pre-existing, not introduced by this change
- Required `go mod tidy && go mod vendor` to add godbus dependency

### Files Modified
1. internal/platform/gateway.go - Added BatchGateway interface, FileInfo struct, updated Request
2. internal/platform/dbus_gateway.go - Updated NewRequest to use Files slice
3. internal/platform/gateway_windows.go - Updated NewRequest to use Files slice
4. internal/platform/gateway_macosx.go - No changes needed (stub implementation)
5. internal/transport/connection.go - Added BatchGateway type assertion for BatchOffer handling
6. internal/app/run.go - Updated outbound processor for Files iteration
7. cmd/drift/main.go - Updated outbound processor for Files iteration
8. go.mod, go.sum, vendor/ - Added godbus dependency

### Gotchas
- Must run `go mod vendor` after adding new dependencies
- Platform-specific build tags mean not all gateway files compile on all platforms
- Type assertion provides graceful fallback for gateways that don't implement BatchGateway

## [2026-02-12] Task: 8 - Build Tag Restructuring for TUI/GUI Coexistence

### Implementation Summary
Restructured Linux gateway to support both TUI (default) and GUI (optional) builds using Go build tags.

### Key Changes
1. **File Rename**: `dbus_gateway.go` → `gateway_linux_tui.go` (preserves git history via `git mv`)
2. **TUI Build Tags**: `//go:build linux && !gui` / `// +build linux,!gui`
   - Ensures TUI gateway only compiles when gui tag is NOT set
   - Default build (no tags) uses TUI implementation
3. **GUI Skeleton**: New `gateway_linux_gui.go` with `//go:build linux && gui` / `// +build linux,gui`
   - Implements both `Gateway` and `BatchGateway` interfaces
   - All methods have stub implementations that log "not implemented"
   - Matches `newGateway()` constructor signature for seamless switching

### Build Tag Syntax Pattern
```go
//go:build linux && !gui
// +build linux,!gui
```
- Both directives required (modern + legacy format)
- Comma in `+build` means AND (not OR)
- Exclamation mark negates condition

### Stub Implementation Pattern
- `guiGateway` struct with same fields as TUI (mu, peers, reqch)
- Each method prints "GUI gateway: [method]() not implemented"
- Returns sensible defaults (nil errors, "DECLINE" for Ask/AskBatch)
- Compiles cleanly without external dependencies

### Asset Structure
- Created `internal/platform/assets/` directory
- Placeholder SVG icon: `drift-icon.svg` (simple gradient circle with arrow)
- Valid XML structure for future integration with gotk4

### Build Verification
- `go build ./cmd/drift` → exit 0 (TUI, default)
- `go build -tags gui ./cmd/drift` → exit 0 (GUI skeleton)
- Both paths compile without errors
- Pre-existing IPv6 vet warnings unrelated to this change

### Key Insight: Build Tag Ordering
When multiple gateway files exist for same platform:
- TUI: `linux && !gui` (matches when gui tag absent)
- GUI: `linux && gui` (matches when gui tag present)
- Exactly one implementation compiles per build
- No runtime selection needed - compile-time choice

### Files Created/Modified
1. `internal/platform/gateway_linux_tui.go` (renamed from dbus_gateway.go)
   - Only change: build tags (line 1-2)
   - All functionality preserved
2. `internal/platform/gateway_linux_gui.go` (new)
   - 50 lines: struct, 6 method stubs, constructor
3. `internal/platform/assets/drift-icon.svg` (new)
   - Simple SVG with gradient and arrow shape

### Next Steps (Not in This Task)
- Import gotk4 in GUI implementation
- Implement actual GTK window and event loop
- Add icon resource embedding
- Platform-specific build instructions for GUI variant

## [2026-02-12] Task: 11 - Linux Desktop Notifications via DBus

### Implementation Pattern
- **File**: `internal/platform/notify_linux.go` with build tags `//go:build linux && gui`
- **DBus Method**: `org.freedesktop.Notifications.Notify` via `github.com/godbus/dbus/v5`
- **Parameters**: app_name="Drift", replaces_id=0, app_icon=path, summary, body, actions=[], hints={}, expire_timeout=5000ms

### Key Learnings
1. **DBus Connection**: Use `dbus.SessionBus()` for user-level notifications (not system bus)
2. **Object Path**: Notifications service uses `/org/freedesktop/Notifications` path
3. **Error Handling**: Gracefully handle DBus connection errors by returning error from SendNotification
4. **Icon Path**: Can pass relative or absolute path to SVG icon file
5. **Non-blocking**: Notify() method doesn't wait for response, just sends and ignores errors (pattern: `_ = SendNotification(...)`)

### Integration
- `gateway_linux_gui.go` Notify() method calls SendNotification with icon path
- Icon path: `internal/platform/assets/drift-icon.svg` (created in Task 8)
- Build succeeds with `-tags gui` flag

### DBus Spec Reference
- Method signature: `Notify(app_name, replaces_id, app_icon, summary, body, actions, hints, expire_timeout) -> uint32`
- Return value (notification ID) is ignored in this implementation
- Hints map can be empty for basic notifications (no urgency, sound, etc.)

## [2026-02-12] Task: 9a - GTK4 Dependency + App Bootstrap

### Implementation Summary
Added gotk4 dependency and implemented GTK4 application bootstrap with window creation and graceful shutdown.

### Key Implementation Details

**Imports Added**:
```go
import (
    "context"
    "fmt"
    "runtime"
    "sync"
    
    "github.com/diamondburned/gotk4/pkg/core/glib"
    gio "github.com/diamondburned/gotk4/pkg/gio/v2"
    gtk "github.com/diamondburned/gotk4/pkg/gtk/v4"
    
    "github.com/metalgrid/drift/internal/zeroconf"
)
```

**Struct Fields Added**:
- `app *gtk.Application` — GTK application instance
- `window *gtk.ApplicationWindow` — Main application window

**Run() Implementation Pattern**:
1. `runtime.LockOSThread()` — Required by GTK4 to ensure main thread execution
2. `gtk.NewApplication("com.github.metalgrid.drift", gio.ApplicationFlagsNone)` — Create app with reverse-domain ID
3. `app.ConnectActivate(func() { ... })` — Register activation callback for window creation
4. Window creation: 400x500 size, HeaderBar with title "Drift", empty vertical Box placeholder
5. Context watcher goroutine: `<-ctx.Done()` → `glib.IdleAdd(func() { app.Quit() })`
6. `app.Run(nil)` — Blocks until app quits

**Shutdown() Implementation**:
- Uses `glib.IdleAdd()` to schedule quit on GTK main thread
- Nil-safe check: `if g.app != nil { g.app.Quit() }`
- Closes request channel after quit scheduled

### Critical Patterns

**Thread Safety with glib.IdleAdd**:
- GTK4 requires all UI operations on main thread
- Cross-goroutine updates use `glib.IdleAdd(func() { /* UI code */ })`
- Context cancellation watcher runs in separate goroutine, schedules Quit via IdleAdd
- Pattern prevents deadlocks and race conditions

**Build Tag Coexistence**:
- TUI: `//go:build linux && !gui` (default)
- GUI: `//go:build linux && gui` (with -tags gui flag)
- Exactly one implementation compiles per build

### Build Verification Results
- `go mod tidy` → SUCCESS
- `go mod vendor` → SUCCESS (synced gotk4 and dependencies)
- `go build -tags gui -o drift-gui ./cmd/drift` → SUCCESS (7.0M binary)
- Binary is valid ELF 64-bit executable with debug info
- `go test ./... -count=1` → ALL PASS (no regressions in non-GUI tests)

### Files Modified
1. `internal/platform/gateway_linux_gui.go`
   - Added gotk4 imports (gtk/v4, gio/v2, core/glib)
   - Added app and window fields to guiGateway struct
   - Implemented Run() with full GTK4 bootstrap
   - Implemented Shutdown() with glib.IdleAdd
2. `go.mod` / `go.sum` — Updated by go mod tidy
3. `vendor/modules.txt` — Updated by go mod vendor

### Key Learnings
1. **gotk4 Import Paths**: Must use exact paths (gtk/v4, gio/v2, core/glib)
2. **Thread Locking**: `runtime.LockOSThread()` is non-negotiable for GTK4
3. **Application ID**: Reverse-domain format (com.github.metalgrid.drift) is convention
4. **ConnectActivate**: Window creation happens in activate callback, not main()
5. **Context Cancellation**: Watcher goroutine + glib.IdleAdd pattern for graceful shutdown
6. **Build Compilation**: CGo compilation takes time (warnings about free() are harmless)

### Verification Checklist
- [x] go mod tidy runs successfully
- [x] go build -tags gui succeeds with valid ELF binary
- [x] Binary is 7.0M (reasonable size for GTK4 app)
- [x] go test ./... passes (no regressions)
- [x] LSP diagnostics clean (no errors)
- [x] Only expected files modified (gateway_linux_gui.go, go.mod, go.sum, vendor/)

### Next Steps (Task 9b onwards)
- Add peer list UI widget
- Implement drag-and-drop for file transfers
- Add file picker dialog
- Implement system tray integration

## [2026-02-12] Task: 9b - Peer List UI with Static Rendering + OnChange Refresh

### Implementation Summary
Added peer list UI widget with automatic refresh on peer changes. Implemented buildPeerList() method that creates a scrollable list of peers with name, OS badge, and IP address.

### Key Implementation Details

**buildPeerList() Method**:
- Creates gtk.ListBox with SelectionNone mode (no row selection)
- Iterates over g.peers.All() snapshot
- For each peer, creates horizontal gtk.Box with:
  - Peer name label (bold via SetMarkup("<b>text</b>"))
  - OS badge from peer.GetRecord("os")
  - First IP address from peer.Addresses[0].String()
- Sets margins (5px top/bottom, 10px start/end) for spacing
- Uses SetHExpand(true) + SetXAlign(0) for left-aligned name label

**Activate Callback Updates**:
- Wraps peer list in gtk.ScrolledWindow with:
  - SetPolicy(PolicyNever, PolicyAutomatic) — vertical scroll only
  - SetVExpand(true) — fills available space
- Appends scrolled window to main vertical box
- Registers peers.OnChange() callback with glib.IdleAdd for thread-safe refresh

**Observer Pattern Integration**:
- peers.OnChange(func() { ... }) registers callback
- Callback wrapped in glib.IdleAdd(func() { ... }) for GTK main thread execution
- On peer change: rebuilds entire list via buildPeerList() and updates scrolled window child
- Thread-safe: observer callback runs on GTK main thread via IdleAdd

### Key Patterns Verified
- gtk.NewScrolledWindow() with SetPolicy(never, automatic) for vertical scrolling
- SetMarkup("<b>text</b>") for bold text rendering
- SetHExpand(true) + SetXAlign(0) for left-aligned labels
- glib.IdleAdd for cross-goroutine UI updates from observer callbacks
- SetSelectionMode(SelectionNone) to disable row selection
- SetMarginTop/Bottom/Start/End for consistent spacing

### Build Verification Results
- `go build -tags gui -o drift-gui ./cmd/drift` → SUCCESS (7.0M binary created)
- `go test ./... -count=1` → ALL PASS (no regressions)
  - config: 0.005s
  - transport: 0.004s
  - zeroconf: 0.003s
- `go vet -tags gui ./internal/platform/` → No errors (CGo warnings are pre-existing)

### Files Modified
- `internal/platform/gateway_linux_gui.go`
  - Added buildPeerList() method (lines 27-60)
  - Updated ConnectActivate callback (lines 67-101)
  - Registered peers.OnChange() with glib.IdleAdd (lines 93-100)

### Code Review Checklist
- [x] buildPeerList() exists and returns gtk.ListBox
- [x] Each peer row: horizontal box with name (bold), OS badge, IP
- [x] ScrolledWindow wraps peer list with vertical scroll policy
- [x] OnChange callback registered with glib.IdleAdd
- [x] Thread-safe UI updates via IdleAdd pattern
- [x] Build succeeds with -tags gui flag
- [x] Tests pass (no regressions)
- [x] Only gateway_linux_gui.go modified

### Next Steps
- Task 9c: Add "Send File" button per peer row
- Task 9d: Add drag-and-drop file sending
- Task 10: System tray integration

## [2026-02-12] Task: 9c - File Picker + Send Button Per Peer

### Implementation Summary
Added "Send File" button to each peer row in buildPeerList(). Button click opens native file chooser dialog with multi-select enabled. Selected files are extracted and sent to reqch channel as Request{To: peer, Files: filePaths}.

### Key Implementation Details

**Button Addition to Peer Row**:
- Created gtk.NewButton() with label "Send File"
- Captured peerInstance via closure before ConnectClicked handler
- Appended button to row after IP label

**File Chooser Dialog Pattern**:
```go
dialog := gtk.NewFileChooserNative(
    "Select Files to Send",
    &g.window.Window,
    gtk.FileChooserActionOpen,
    "Send",
    "Cancel",
)
dialog.SetSelectMultiple(true)
dialog.Show()
```

**Response Handler Pattern**:
- ConnectResponse callback checks responseID == int(gtk.ResponseAccept)
- dialog.GetFiles() returns gio.ListModel
- Iterate with NItems() and Item(i) → cast to *gio.File
- (*gio.File).Path() extracts file system path
- Send Request{To: peerInstance, Files: filePaths} directly to g.reqch

**NewRequest() Update**:
- Removed fmt.Println stub line
- Kept fmt import (still used by Ask() and AskBatch() stubs)
- Direct channel send: g.reqch <- Request{To: to, Files: []string{file}}

### Key Patterns Verified
- gtk.NewFileChooserNative(title, parent, action, acceptLabel, cancelLabel)
- SetSelectMultiple(true) for multi-file selection
- ConnectResponse(func(responseID int) { ... }) for dialog response handling
- dialog.GetFiles() → gio.ListModel → NItems() / Item(i) → cast to *gio.File
- (*gio.File).Path() extracts file system path string
- Closure capture of peerInstance before button click handler
- Direct channel send (no async spawning needed)

### Build Verification Results
- `go build -tags gui -o drift-gui ./cmd/drift` → SUCCESS (7.0M binary)
- `go test ./... -count=1` → ALL PASS (no regressions)
  - config: 0.003s
  - transport: 0.003s
  - zeroconf: 0.003s
- `go vet -tags gui ./internal/platform/` → No errors (CGo warnings pre-existing)
- LSP diagnostics: No errors

### Files Modified
- `internal/platform/gateway_linux_gui.go`
  - Added "Send File" button to buildPeerList() (lines 55-95)
  - Updated NewRequest() to remove fmt.Println stub (lines 165-168)
  - Kept fmt import for Ask() and AskBatch() stubs

### Code Review Checklist
- [x] "Send File" button exists in each peer row
- [x] Button appended after IP label, before listBox.Append(row)
- [x] gtk.NewFileChooserNative with SelectMultiple(true)
- [x] ConnectResponse callback extracts file paths correctly
- [x] Request sent to reqch with To and Files fields
- [x] NewRequest() no longer has fmt.Println stub
- [x] Build succeeds with -tags gui flag
- [x] Tests pass (no regressions)
- [x] Only gateway_linux_gui.go modified (plus binary timestamp)

### Next Steps
- Task 9d: Add drag-and-drop file sending
- Task 10: System tray integration

## [2026-02-12] Task: 9d - Drag-and-Drop File Sending

### Implementation Summary
Added gtk.DropTarget to each peer row in buildPeerList(). DropTarget accepts file:// URIs from file manager drag-and-drop operations, parses the path, and sends Request{To: peer, Files: [path]} to reqch channel.

### Key Implementation Details

**Imports Added**:
```go
import (
    "strings"
    "github.com/diamondburned/gotk4/pkg/gdk/v4"
)
```

**DropTarget Creation Pattern** (in buildPeerList(), after "Send File" button):
```go
// Capture peer instance before drop handler (closure safety)
peerInstanceForDrop := peer.GetInstance()

// Create drop target accepting string data with copy action
drop := gtk.NewDropTarget(glib.TypeString, gdk.ActionCopy)

// Connect drop signal handler
drop.ConnectDrop(func(drop *gtk.DropTarget, val *glib.Value, x, y float64) bool {
    // Extract string from GLib Value
    str, ok := val.GoValue().(string)
    if !ok {
        return false
    }
    
    // Validate file:// URI prefix
    if !strings.HasPrefix(str, "file://") {
        return false
    }
    
    // Parse URI to file path
    path := strings.TrimPrefix(str, "file://")
    path = strings.TrimSpace(path)
    
    if path == "" {
        return false
    }
    
    // Send request to reqch
    g.reqch <- Request{To: peerInstanceForDrop, Files: []string{path}}
    return true  // Accept drop
})

// Attach drop target to row widget
row.AddController(drop)
```

### Key Patterns Verified
- `gtk.NewDropTarget(glib.TypeString, gdk.ActionCopy)` creates drop target for string data
- `ConnectDrop(func(drop, val, x, y) bool { ... })` connects drop signal handler
- `val.GoValue().(string)` extracts Go string from GLib Value with type assertion
- `strings.TrimPrefix(str, "file://")` removes URI scheme
- `strings.TrimSpace(path)` removes whitespace from parsed path
- Return `true` from handler to accept drop, `false` to reject
- `row.AddController(drop)` attaches drop target to widget
- Capture `peerInstanceForDrop` before handler to avoid closure issues with loop variable

### GTK4 DnD Behavior
- DropTarget with TypeString handles single URI per drop event
- File manager sends multiple drop events for multiple files (sequential)
- Each drop event triggers separate ConnectDrop callback
- This is GTK4 limitation, not implementation issue

### Build Verification Results
- `go build -tags gui -o drift-gui ./cmd/drift` → SUCCESS (7.0M binary)
- `go test ./... -count=1` → ALL PASS (no regressions)
  - config: 0.002s
  - transport: 0.005s
  - zeroconf: 0.004s
- `go vet -tags gui ./internal/platform/` → No errors (CGo warnings pre-existing)
- LSP diagnostics: No errors

### Files Modified
- `internal/platform/gateway_linux_gui.go`
  - Added "strings" import (line 10)
  - Added "github.com/diamondburned/gotk4/pkg/gdk/v4" import (line 14)
  - Added DropTarget creation and handler in buildPeerList() (lines 100-128)
  - Attached drop target via row.AddController(drop) (line 128)

### Code Review Checklist
- [x] gdk import added with correct path
- [x] strings import added for TrimPrefix/TrimSpace
- [x] DropTarget created with glib.TypeString and gdk.ActionCopy
- [x] ConnectDrop handler parses file:// URI correctly
- [x] peerInstanceForDrop captured before handler (closure safety)
- [x] Request sent to reqch with peer instance and file path
- [x] Drop target attached via row.AddController(drop)
- [x] Build succeeds with -tags gui flag
- [x] Tests pass (no regressions)
- [x] Only gateway_linux_gui.go modified (plus binary timestamp)

### Integration with Existing Features
- Works alongside "Send File" button (Task 9c) - both methods send to same reqch
- Peer list refresh (Task 9b) rebuilds rows with new drop targets
- Observer pattern (Task 9a) triggers list rebuild on peer changes
- Thread-safe: drop handler runs on GTK main thread (no glib.IdleAdd needed)

### Next Steps
- Task 10: System tray integration (DBus StatusNotifierItem)
- Task 12: Transfer progress UI
- Task 13a/13b: Incoming transfer dialog

## [2026-02-12 01:02] Task: 10 - DBus StatusNotifierItem System Tray

### Implementation
- Created tray_linux.go with DBus SNI implementation
- Exported object on /StatusNotifierItem path
- Registered with org.kde.StatusNotifierWatcher
- Implemented required properties: Category, Id, Title, IconName, ItemIsMenu
- Implemented Activate (toggle window) and SecondaryActivate (quit)
- Integrated into guiGateway.Run() with callbacks
- Window starts hidden, tray click shows/hides
- Close button hides to tray instead of quitting

### Key Patterns
- dbus.SessionBus() for user-level DBus
- conn.Export(obj, path, interface) exports object
- Property getters: func() (T, *dbus.Error)
- Method handlers: func(args...) *dbus.Error
- glib.IdleAdd for GTK calls from DBus callbacks
- ConnectCloseRequest returns bool (true = prevent close)
- window.IsVisible() to check visibility state
- introspect.NewIntrospectable() for DBus introspection support

### DBus SNI Spec
- Interface: org.kde.StatusNotifierItem
- Watcher: org.kde.StatusNotifierWatcher at /StatusNotifierWatcher
- Activate = left-click, SecondaryActivate = right-click
- ItemIsMenu = false for click activation (not menu)
- Required properties: Category, Id, Title, IconName, ItemIsMenu
- Registration: Call RegisterStatusNotifierItem on watcher with object path

### Critical Implementation Details
- Export introspection data via introspect.NewIntrospectable()
- Pass object path "/StatusNotifierItem" to watcher (not full bus name)
- Close DBus connection in Shutdown() before quitting app
- Error handling: gracefully continue without tray if registration fails
- Window.SetVisible(false) starts hidden (tray shows it)
- ConnectCloseRequest(func() bool) - return true prevents default close

### Verification
- Build with -tags gui: SUCCESS (7.0M binary)
- Tests: ALL PASS (no regressions)
- Vet: No errors (CGo warnings pre-existing)
- Files: tray_linux.go (new), gateway_linux_gui.go (modified)

### Next Steps
- Task 12: Transfer progress UI
- Task 13: Incoming transfer dialog

## [2026-02-12 01:15] Task: 12 - Transfer Progress UI

### Implementation Summary
Added transfer progress tracking UI to Linux GUI with real-time speed calculation and thread-safe updates. Implemented TransferState struct, buildTransferList() method, and UpdateTransfer() method with glib.IdleAdd for cross-goroutine UI updates.

### Key Implementation Details

**TransferState Struct**:
```go
type TransferState struct {
    ID         string
    PeerName   string
    Filename   string
    Direction  string    // "↑" for upload, "↓" for download
    Total      int64
    Current    int64
    Speed      float64   // bytes per second
    LastUpdate time.Time
    Status     string    // "active", "complete", "failed"
}
```

**guiGateway Fields Added**:
- `transfers map[string]*TransferState` — keyed by transfer ID
- `transferList *gtk.ListBox` — widget for displaying transfers
- `transferBox *gtk.Box` — container for transfer section

**buildTransferList() Method**:
- Creates gtk.ListBox with SelectionNone mode
- Locks mutex and iterates over transfers map
- Only shows transfers with Status == "active"
- Each row: horizontal box with:
  - Direction + filename label (left-aligned, expandable)
  - Progress bar (200px width, SetFraction for 0.0-1.0 range)
  - Speed + percentage label (e.g., "1.2 MB/s - 45%")
- Margins: 5px top/bottom, 10px start/end

**UpdateTransfer() Method**:
- Locks mutex, finds transfer by ID
- Calculates speed: `delta_bytes / delta_time` (bytes per second)
- Updates Current and LastUpdate fields
- Unlocks mutex
- Uses `glib.IdleAdd()` to rebuild list on GTK main thread
- Removes old child from transferBox, appends new list

**Integration into Run() activate callback**:
- Added after peer list scrolled window
- Creates bold "Active Transfers" label
- Creates transferBox (vertical box) and builds initial transfer list
- Appends both label and box to main window

**Initialization in newGateway()**:
- Changed from positional struct initialization to named fields
- Added `transfers: make(map[string]*TransferState)`

### Key Patterns Verified

**Speed Calculation**:
- Track LastUpdate timestamp on each UpdateTransfer call
- Calculate elapsed time: `now.Sub(transfer.LastUpdate).Seconds()`
- Calculate delta bytes: `current - transfer.Current`
- Speed = delta / elapsed (bytes per second)
- Display: MB/s = speed / 1024 / 1024

**Progress Bar API**:
- `gtk.NewProgressBar()` creates widget
- `SetFraction(float64)` sets progress (0.0 to 1.0)
- `SetSizeRequest(width, -1)` sets minimum width
- Percentage calculation: `(current / total) * 100`

**Thread Safety**:
- Mutex protects transfers map access
- All widget updates via `glib.IdleAdd()` on GTK main thread
- Pattern: lock → read/modify → unlock → IdleAdd for UI update

**List Rebuild Strategy**:
- Simple MVP approach: rebuild entire list on each update
- Remove old child via `FirstChild()` and `Remove()`
- Append new list via `Append(newList)`
- Avoids complex row update logic

### Build Verification Results
- `go build -tags gui -o drift-gui ./cmd/drift` → SUCCESS (7.0M binary)
- `go test ./... -count=1` → ALL PASS (no regressions)
  - config: 0.003s
  - transport: 0.003s
  - zeroconf: 0.004s
- `go vet -tags gui ./internal/platform/` → No errors (CGo warnings pre-existing)
- `git status --short` → ONLY gateway_linux_gui.go modified (plus binary timestamp)

### Files Modified
- `internal/platform/gateway_linux_gui.go`
  - Added time import (line 12)
  - Added TransferState struct (lines 22-32)
  - Added transfers, transferList, transferBox fields to guiGateway (lines 41-43)
  - Added buildTransferList() method (lines 153-196)
  - Added UpdateTransfer() method (lines 198-227)
  - Integrated transfer section into Run() activate callback (lines 257-268)
  - Updated newGateway() to initialize transfers map (lines 361-368)

### Code Review Checklist
- [x] TransferState struct defined with all required fields
- [x] transfers map and transferList/transferBox fields added to guiGateway
- [x] buildTransferList() creates progress bars with direction, filename, speed
- [x] UpdateTransfer() calculates speed and updates UI via glib.IdleAdd
- [x] Transfer section added to window below peer list
- [x] transfers map initialized in newGateway()
- [x] time import added
- [x] Build succeeds with -tags gui flag
- [x] Tests pass (no regressions)
- [x] Only gateway_linux_gui.go modified

### Integration with Existing Features
- Works alongside peer list (Task 9b) and system tray (Task 10)
- Observer pattern (Task 9a) triggers peer list rebuild on changes
- Transfer section appears below peer list in window layout
- Thread-safe: all UI updates via glib.IdleAdd pattern

### Next Steps
- Task 13: Incoming transfer dialog (Ask/AskBatch implementation)
- Task 14b: Wire progress callbacks to UpdateTransfer()
- Task 14c: Add transfer completion/error handling

### Key Learnings
1. **Speed Calculation**: Must track LastUpdate timestamp and calculate delta bytes / delta time
2. **Progress Bar**: SetFraction expects 0.0-1.0 range, not percentage
3. **List Rebuild**: Simple approach (rebuild entire list) works well for MVP
4. **Thread Safety**: Mutex for data access, glib.IdleAdd for UI updates
5. **GTK4 Patterns**: SetMarkup for bold text, SetHExpand/SetXAlign for layout
6. **Closure Safety**: Capture variables before loop to avoid closure issues
