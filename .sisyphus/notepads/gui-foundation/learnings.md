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
