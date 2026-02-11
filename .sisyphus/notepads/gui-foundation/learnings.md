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
