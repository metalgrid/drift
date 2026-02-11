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
