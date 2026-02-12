# Linux GUI for Drift

## TL;DR

> **Quick Summary**: Build a GTK4-based Linux GUI for Drift's peer-to-peer file transfer, replacing the terminal REPL with a system tray icon + panel window that supports drag-and-drop sending, batch file transfers, live progress tracking, and incoming transfer notifications with countdown timers.
> 
> **Deliverables**:
> - GTK4 panel window with live peer list, drag-and-drop file sending, and transfer progress
> - DBus StatusNotifierItem (system tray icon) for unobtrusive presence
> - BATCH_OFFER protocol extension for multi-file transfers
> - Progress-tracking transport wrappers (io.Writer/Reader with callbacks)
> - Peer change observer (callback-based notifications from zeroconf.Peers)
> - XDG-compliant config system (~/.config/drift/config.toml)
> - Incoming transfer dialog with countdown timer + desktop notifications
> - Build tag coexistence: GUI (opt-in via `-tags gui`) alongside renamed TUI
> - Full TDD test infrastructure from scratch
> 
> **Estimated Effort**: XL
> **Parallel Execution**: YES - 6 waves
> **Critical Path**: Tests → Protocol Extension → Gateway Interface → Build Tags → GTK4 GUI → Integration

---

## Context

### Original Request
Build a Linux GUI for Drift that's unobtrusive and usable. Sending: show peer list, select peer + file picker OR drag-and-drop files onto peer. Receiving: notification with peer name, file count, size, configurable timeout (default 30s), accept/reject. Auto-reject on timeout.

### Interview Summary
**Key Discussions**:
- **Presence model**: Tray icon + panel window (tray for background presence, panel for peer list and operations)
- **Drag-and-drop**: From file manager onto specific peer entry in the peer list
- **Transfer progress**: Live progress in app window AND desktop notifications for start/complete
- **Multi-file**: Batch support via protocol extension (BATCH_OFFER message type)
- **Toolkit**: GTK4 via gotk4 — CGo required, AGPL-3.0 accepted, native Linux look, excellent DnD
- **Incoming UX**: Desktop notification fires first → click opens dialog with countdown, file list, Accept/Decline
- **REPL**: Preserved via build tags (not replaced)
- **Config**: XDG config file (~/.config/drift/config.toml) for download dir, timeout, identity
- **Tests**: Full TDD from scratch (project has zero test files)

**Research Findings**:
- Gateway interface: 5 methods (`Run`, `Shutdown`, `NewRequest`, `Ask`, `Notify`), `Ask` returns exactly `"ACCEPT"` or `"DECLINE"`
- Protocol: single-file OFFER (`OFFER|filename|mimetype|size\n`), ANSWER (`ANSWER|ACCEPT\n` or `ANSWER|DECLINE\n`)
- No progress callbacks exist — `file.WriteTo(conn)` and `io.LimitReader` block without intermediate feedback
- No peer change notifications — `Peers` struct has `All()`, `GetByInstance()`, `GetByAddr()` but no observer pattern
- `app/run.go` outbound processor has `return` instead of `continue` on errors (lines 128, 133, 143) — kills all outbound processing on first error
- GTK4's `DropTarget` API is mature for file DnD; gotk4 has working examples
- DBus StatusNotifierItem is the modern Linux tray standard (works on KDE natively, GNOME with AppIndicator extension, XFCE)
- `xdg` package (`github.com/adrg/xdg`) already in go.mod

### Metis Review
**Identified Gaps** (addressed):
- `Gateway` interface needs evolution for batch transfers — resolved: add `AskBatch` method, keep `Ask` for backwards compat
- `Peers` has no change notification — resolved: add callback registration (`OnChange`)
- Build tag scheme ambiguity — resolved: `gui` tag for opt-in GUI, TUI stays default
- `app/run.go` outbound processor bug (`return` → `continue`) — resolved: include fix in plan
- BATCH_OFFER backwards compatibility — resolved: require matching versions for v1 (no version negotiation)
- Partial batch acceptance — resolved: all-or-nothing for v1
- `context.WithValue` with string key `"filename"` — resolved: replace with typed key in protocol refactoring task
- Tray icon asset needed — resolved: include placeholder SVG icon task

---

## Work Objectives

### Core Objective
Replace Drift's Linux terminal REPL with a GTK4 GUI that provides system tray presence and a panel window for peer discovery, file sending (via picker or drag-and-drop), batch transfer management with progress tracking, and incoming transfer acceptance with configurable timeout.

### Concrete Deliverables
- `internal/config/config.go` — Config system with TOML parsing and XDG paths
- `internal/transport/progress.go` — Progress-tracking io.Writer/Reader wrappers
- `internal/transport/message.go` — Extended with `BatchOffer` type
- `internal/transport/connection.go` — Extended with batch send/receive and progress callbacks
- `internal/zeroconf/zeroconf.go` — Extended with peer change observer callbacks
- `internal/platform/gateway.go` — Extended with `BatchGateway` interface
- `internal/platform/gateway_linux_tui.go` — Renamed from `gateway_linux.go` with `nogui` build tag
- `internal/platform/gateway_linux_gui.go` — New GTK4 GUI gateway with `gui` build tag
- `internal/platform/tray_linux.go` — DBus StatusNotifierItem implementation
- `internal/platform/assets/drift-icon.svg` — Placeholder tray icon
- `internal/app/run.go` — Bug fix (return→continue) + config integration + batch support
- Test files for all new packages

### Definition of Done
- [ ] `go build -tags gui ./cmd/drift` produces a working binary with GTK4 GUI
- [ ] `go build ./cmd/drift` produces the TUI binary (backwards compatible)
- [ ] `go test ./...` passes all tests (no `gui` tag needed for non-GUI tests)
- [ ] `go vet ./...` reports no errors
- [ ] Batch file transfer works between two GUI instances on the same LAN
- [ ] Drag-and-drop from file manager onto peer in list initiates transfer
- [ ] Incoming transfer shows notification → dialog with countdown → accept/decline works
- [ ] Tray icon appears and panel window toggles on click

### Must Have
- GTK4 panel window with live peer list (auto-updates on peer discovery changes)
- File drag-and-drop from file manager onto peer entries in the list
- File picker dialog as alternative to drag-and-drop
- BATCH_OFFER protocol for multi-file transfers (all-or-nothing acceptance)
- Live transfer progress (percentage, speed, ETA) in the panel window
- Desktop notifications for incoming transfers and completion events
- Incoming transfer dialog with countdown timer, file list, total size, Accept/Decline
- Configurable timeout (default 30s) for incoming transfer acceptance
- Configurable download directory (default ~/Downloads/Drift/)
- System tray icon via DBus StatusNotifierItem
- Build tag coexistence with TUI (TUI as default, GUI as opt-in)
- Full test coverage for protocol, transport, and config layers

### Must NOT Have (Guardrails)
- **NO** changes to existing OFFER/ANSWER wire format — BATCH_OFFER is purely additive
- **NO** version negotiation in protocol (require matching versions for v1)
- **NO** selective batch acceptance (all-or-nothing for v1)
- **NO** tray icon animation or transfer progress in tray icon
- **NO** file type icons, thumbnails, or MIME sniffing in peer list
- **NO** drag-and-drop FROM app TO file manager (inbound DnD only)
- **NO** notification action buttons (notifications are informational; dialog handles accept/decline)
- **NO** compression, resumable transfers, or encryption layer changes
- **NO** refactoring of existing transport error handling beyond the specific `return`→`continue` bug fix
- **NO** config options beyond: download_dir, auto_accept_timeout, identity (keep config minimal)
- **NO** gotk4 imports in any file without `gui` build tag — strict isolation

---

## Verification Strategy

> **UNIVERSAL RULE: ZERO HUMAN INTERVENTION**
>
> ALL tasks in this plan MUST be verifiable WITHOUT any human action.
> Every criterion is verifiable by running a command or using a tool.

### Test Decision
- **Infrastructure exists**: NO (zero test files in project)
- **Automated tests**: YES (Full TDD)
- **Framework**: `go test` (stdlib, no external test framework — testify is in go.mod but unused)

### Test Infrastructure Setup
Task 1 establishes the testing foundation. Every subsequent task follows RED-GREEN-REFACTOR.

### Agent-Executed QA Scenarios (MANDATORY — ALL tasks)

**Verification Tool by Deliverable Type:**

| Type | Tool | How Agent Verifies |
|------|------|-------------------|
| Protocol/Transport (pure Go) | Bash (`go test`) | Run tests with `-v -run` flags, assert PASS |
| Config system (pure Go) | Bash (`go test`) | Run tests, also verify file I/O with temp dirs |
| GTK4 GUI | Bash (`go build -tags gui`) | Compile check; runtime via tmux + dbus-monitor |
| System tray (DBus) | Bash (dbus-send/dbus-monitor) | Introspect DBus to verify SNI registration |
| Desktop notifications | Bash (dbus-monitor) | Monitor org.freedesktop.Notifications calls |
| Build system | Bash (`go build`/`go vet`) | Multi-target build verification |

---

## Execution Strategy

### Parallel Execution Waves

```
Wave 1 (Start Immediately):
├── Task 1: Test infrastructure setup
├── Task 2: Config system (internal/config/)
└── Task 3: Bug fix — outbound processor return→continue

Wave 2 (After Wave 1):
├── Task 4: Peer change observer (zeroconf.Peers)
├── Task 5: Progress tracking wrappers (transport/progress.go)
└── Task 6: Protocol extension — BATCH_OFFER message type

Wave 3 (After Wave 2):
├── Task 7: Gateway interface evolution (BatchGateway)
└── Task 8: Build tag restructuring (TUI rename + GUI skeleton)

Wave 4 (After Wave 3):
├── Task 9: GTK4 panel window (peer list, DnD, file picker)
├── Task 10: DBus StatusNotifierItem (system tray)
└── Task 11: Desktop notifications (libnotify via DBus)

Wave 5 (After Wave 4):
├── Task 12: Transfer progress UI (progress bars in panel)
└── Task 13: Incoming transfer dialog (countdown, accept/decline)

Wave 6 (After Wave 5):
├── Task 14: Integration — wire GUI gateway to app + config
└── Task 15: Final verification & cross-platform build check
```

### Dependency Matrix

| Task | Depends On | Blocks | Can Parallelize With |
|------|------------|--------|---------------------|
| 1 | None | 4, 5, 6 | 2, 3 |
| 2 | None | 14 | 1, 3 |
| 3 | None | 14 | 1, 2 |
| 4 | 1 | 9 | 5, 6 |
| 5 | 1 | 12, 14 | 4, 6 |
| 6 | 1 | 7 | 4, 5 |
| 7 | 6 | 8, 9, 13 | — |
| 8 | 7 | 9, 10, 11 | — |
| 9 | 4, 8 | 12, 13 | 10, 11 |
| 10 | 8 | 14 | 9, 11 |
| 11 | 8 | 13 | 9, 10 |
| 12 | 5, 9 | 14 | 13 |
| 13 | 7, 9, 11 | 14 | 12 |
| 14 | 2, 3, 5, 10, 12, 13 | 15 | — |
| 15 | 14 | None | — |

### Agent Dispatch Summary

| Wave | Tasks | Recommended Dispatch |
|------|-------|---------------------|
| 1 | 1, 2, 3 | 3 parallel `task(category="quick")` |
| 2 | 4, 5, 6 | 3 parallel `task(category="unspecified-low")` |
| 3 | 7, 8 | Sequential: 7 then 8 (`task(category="unspecified-low")`) |
| 4 | 9, 10, 11 | 3 parallel `task(category="visual-engineering")` |
| 5 | 12, 13 | 2 parallel `task(category="visual-engineering")` |
| 6 | 14, 15 | Sequential: 14 then 15 (`task(category="deep")`) |

---

## TODOs

- [x] 1. Set Up Test Infrastructure + Baseline Protocol Tests

  **What to do**:
  - Create `internal/transport/message_test.go` with tests for existing message types:
    - `TestOfferMarshalMessage` — verify `OFFER|filename|mimetype|size\n` format
    - `TestOfferMarshalRoundTrip` — marshal → unmarshal → compare
    - `TestAnswerMarshalAccept` — verify `ANSWER|ACCEPT\n`
    - `TestAnswerMarshalDecline` — verify `ANSWER|DECLINE\n`
    - `TestUnmarshalInvalidMessage` — garbage input returns nil/error
    - `TestUnmarshalMalformedOffer` — wrong field count, bad size
    - `TestMakeOffer` — verify file stat integration (use temp file)
    - `TestFormatSize` — verify KiB/MiB/GiB/TiB formatting
  - Create `internal/transport/connection_test.go` with tests for `storeFile`:
    - `TestStoreFileSuccess` — write temp file, verify content and rename
    - `TestStoreFileZeroBytes` — zero-size file handled correctly
    - `TestStoreFileDirCreation` — target directory created if missing

  **Must NOT do**:
  - Do NOT add any external test framework (no testify) — use stdlib `testing` package
  - Do NOT write tests for platform-specific code (that needs GTK4/CGo)
  - Do NOT modify any existing source files in this task

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Straightforward test file creation following Go stdlib testing conventions
  - **Skills**: []
    - No special skills needed — standard Go testing
  - **Skills Evaluated but Omitted**:
    - `playwright`: No browser interaction
    - `frontend-ui-ux`: No UI work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 2, 3)
  - **Blocks**: Tasks 4, 5, 6 (they follow TDD and need test infrastructure patterns established)
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `internal/transport/message.go:16-39` — `Offer` struct and `MarshalMessage()` method — these are the exact types to test
  - `internal/transport/message.go:41-60` — `Answer` struct and marshaling — test both ACCEPT and DECLINE
  - `internal/transport/message.go:62-94` — `UnmarshalMessage()` function — switch on prefix, returns `any` (Offer, Answer, or error)
  - `internal/transport/message.go:96-108` — `MakeOffer()` — needs temp file for testing since it calls `os.Stat()`
  - `internal/transport/message.go:124-144` — `formatSize()` — unexported, test via `MakeOffer` or make a test in same package
  - `internal/transport/connection.go:88-112` — `storeFile()` — unexported, test in same package; uses temp file → atomic rename

  **Documentation References**:
  - Go testing stdlib: `go doc testing` — standard test file conventions

  **Acceptance Criteria**:

  **TDD (RED-GREEN-REFACTOR):**
  - [ ] Test files created: `internal/transport/message_test.go`, `internal/transport/connection_test.go`
  - [ ] `go test ./internal/transport/ -v` → PASS (all tests green)
  - [ ] Minimum 10 test cases covering marshal, unmarshal, roundtrip, edge cases, storeFile

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: All transport tests pass
    Tool: Bash
    Preconditions: None
    Steps:
      1. Run: go test ./internal/transport/ -v -count=1
      2. Assert: exit code 0
      3. Assert: output contains "PASS"
      4. Assert: output contains "TestOfferMarshalMessage"
      5. Assert: output contains "TestUnmarshalInvalidMessage"
      6. Assert: output contains "TestStoreFileSuccess"
    Expected Result: All tests pass, 0 failures
    Evidence: Test output captured to stdout

  Scenario: Test count verification
    Tool: Bash
    Preconditions: Tests from step 1 exist
    Steps:
      1. Run: go test ./internal/transport/ -v -count=1 2>&1 | grep -c "--- PASS"
      2. Assert: count >= 10
    Expected Result: At least 10 passing test cases
    Evidence: grep count output
  ```

  **Commit**: YES
  - Message: `test(transport): add baseline tests for message marshal/unmarshal and file storage`
  - Files: `internal/transport/message_test.go`, `internal/transport/connection_test.go`
  - Pre-commit: `go test ./internal/transport/ -v`

---

- [x] 2. Config System

  **What to do**:
  - Create `internal/config/config.go`:
    - `Config` struct with fields: `DownloadDir string`, `AcceptTimeout time.Duration`, `Identity string`
    - `DefaultConfig()` function returning sensible defaults: `DownloadDir = xdg.UserDirs.Download + "/Drift"`, `AcceptTimeout = 30 * time.Second`, `Identity = ""`
    - `Load(path string) (*Config, error)` — reads TOML file, merges with defaults
    - `DefaultPath()` — returns `~/.config/drift/config.toml` (via `xdg.ConfigHome`)
    - `EnsureConfigDir()` — creates `~/.config/drift/` with 0700 permissions if missing
  - Add `github.com/BurntSushi/toml` dependency (or `github.com/pelletier/go-toml/v2`)
  - Create `internal/config/config_test.go` (TDD):
    - RED: Write tests first for `DefaultConfig`, `Load` (valid file), `Load` (missing file → defaults), `Load` (corrupt file → defaults + error logged)
    - GREEN: Implement `config.go` to pass tests
    - REFACTOR: Clean up

  **Must NOT do**:
  - Do NOT add config options beyond `download_dir`, `accept_timeout`, `identity`
  - Do NOT make config required — missing file = use all defaults silently
  - Do NOT watch config file for changes — reload only on app restart

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: Small, self-contained package with clear interface
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: No UI work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 3)
  - **Blocks**: Task 14 (integration wires config into app)
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `internal/transport/connection.go:49` — Current hardcoded download path: `filepath.Join(xdg.UserDirs.Download, "Drift")` — config replaces this
  - `internal/platform/gateway_linux.go:162,169` — Current hardcoded 30-second timeout — config replaces this
  - `go.mod:6` — `github.com/adrg/xdg v0.5.3` already available for `xdg.ConfigHome`

  **API/Type References**:
  - `github.com/adrg/xdg` — `xdg.ConfigHome` for config directory, `xdg.UserDirs.Download` for default download dir

  **External References**:
  - BurntSushi/toml: `https://github.com/BurntSushi/toml` — TOML parser for Go
  - XDG Base Directory Spec: config file goes in `$XDG_CONFIG_HOME/drift/config.toml`

  **Acceptance Criteria**:

  **TDD:**
  - [ ] Test file: `internal/config/config_test.go`
  - [ ] `go test ./internal/config/ -v` → PASS
  - [ ] Tests cover: default config, load from file, missing file fallback, corrupt file fallback, config dir creation

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Config tests pass
    Tool: Bash
    Preconditions: None
    Steps:
      1. Run: go test ./internal/config/ -v -count=1
      2. Assert: exit code 0
      3. Assert: output contains "PASS"
    Expected Result: All config tests pass
    Evidence: Test output captured

  Scenario: Default config has expected values
    Tool: Bash
    Preconditions: Config package exists
    Steps:
      1. Run: go test ./internal/config/ -v -run TestDefaultConfig
      2. Assert: exit code 0
      3. Assert: test validates download_dir ends with "/Drift"
      4. Assert: test validates accept_timeout is 30s
    Expected Result: Defaults match specification
    Evidence: Test output captured
  ```

  **Commit**: YES
  - Message: `feat(config): add XDG-compliant TOML config system with download dir, timeout, and identity`
  - Files: `internal/config/config.go`, `internal/config/config_test.go`, `go.mod`, `go.sum`
  - Pre-commit: `go test ./internal/config/ -v`

---

- [x] 3. Bug Fix — Outbound Processor `return` → `continue`

  **What to do**:
  - In `internal/app/run.go`, fix the outbound connection processor goroutine (lines 110-151):
    - Line 128: `return` after dial failure → change to `continue`
    - Line 133: `return` after public key decode failure → change to `continue` (also close `conn`)
    - Line 143: `return` after secure connection failure → change to `continue` (also close `conn`)
  - Ensure `conn.Close()` is called before `continue` on error paths where a connection was established (lines 133, 143)
  - Also replace the `context.WithValue` string key `"filename"` (line 146) with a typed context key:
    - Add `type contextKey string` and `const filenameKey contextKey = "filename"` in `internal/app/run.go` (or a shared location)
    - Update line 146: `context.WithValue(ctx, filenameKey, request.File)`
    - Update `internal/transport/connection.go:58`: `ctx.Value(filenameKey)` — BUT this creates a cross-package dependency. Better: move the context key to `internal/transport/` package and export it.

  **Must NOT do**:
  - Do NOT refactor the entire outbound processor structure
  - Do NOT change the channel-based architecture
  - Do NOT add error logging beyond what exists (the `Notify` calls already inform the user)

  **Recommended Agent Profile**:
  - **Category**: `quick`
    - Reason: 3-line fix (return→continue) + small context key refactor
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `git-master`: Simple change, no complex git operations

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 1 (with Tasks 1, 2)
  - **Blocks**: Task 14 (integration)
  - **Blocked By**: None

  **References**:

  **Pattern References**:
  - `internal/app/run.go:110-151` — The outbound connection processor goroutine. Lines 128, 133, 143 have `return` that should be `continue`. The `for/select` loop processes transfer requests — `return` kills the entire goroutine, meaning no more outbound transfers work after the first failure.
  - `internal/app/run.go:146` — `context.WithValue(ctx, "filename", request.File)` — string key should be typed
  - `internal/transport/connection.go:58` — `ctx.Value("filename").(string)` — consumer of the context value
  - `internal/app/run.go:66-108` — Inbound processor for comparison — correctly uses `continue` on error paths (lines 79, 86, 93, 103)

  **Acceptance Criteria**:

  - [ ] `internal/app/run.go` lines 128, 133, 143: `return` replaced with `continue` (+ `conn.Close()` where applicable)
  - [ ] `context.WithValue` string key replaced with typed `contextKey`
  - [ ] `go vet ./...` → no errors
  - [ ] `go build ./cmd/drift` → exit 0

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Build succeeds after fix
    Tool: Bash
    Preconditions: None
    Steps:
      1. Run: go vet ./internal/app/
      2. Assert: exit code 0
      3. Run: go build ./cmd/drift
      4. Assert: exit code 0
      5. Assert: binary exists at ./drift
    Expected Result: Clean build with no vet warnings
    Evidence: Build output captured

  Scenario: Context key is typed (not string)
    Tool: Bash
    Preconditions: Fix applied
    Steps:
      1. Search for string "filename" used as context key in app/run.go
      2. Assert: no `context.WithValue(ctx, "filename"` pattern exists
      3. Assert: typed key constant is defined
    Expected Result: No string-typed context keys remain
    Evidence: grep output captured
  ```

  **Commit**: YES
  - Message: `fix(app): prevent outbound processor goroutine death on single transfer failure`
  - Files: `internal/app/run.go`, `internal/transport/connection.go`
  - Pre-commit: `go vet ./... && go build ./cmd/drift`

---

- [x] 4. Peer Change Observer

  **What to do**:
  - Extend `internal/zeroconf/zeroconf.go` `Peers` struct with observer callbacks:
    - Add `observers []func()` field to `Peers` struct
    - Add `OnChange(fn func())` method — registers a callback invoked on peer add/remove
    - Call all registered observers at the end of `add()` and `remove()` methods
    - Observers are called with the mutex NOT held (unlock before calling to prevent deadlocks)
  - Create `internal/zeroconf/zeroconf_test.go` (TDD):
    - RED: `TestPeersOnChangeCalledOnAdd` — register observer, add peer, assert callback fired
    - RED: `TestPeersOnChangeCalledOnRemove` — register observer, add then remove peer, assert callback fired for both
    - RED: `TestPeersMultipleObservers` — register 2 observers, add peer, both fire
    - RED: `TestPeersAllReturnsCopy` — verify `All()` returns snapshot (existing behavior, baseline test)
    - GREEN: Implement observer pattern
    - REFACTOR: Clean up

  **Must NOT do**:
  - Do NOT use generic event system or reactive streams — simple `func()` callback
  - Do NOT pass peer info to the callback (caller should call `All()` to get current state)
  - Do NOT change the `Peers` constructor signature — add `OnChange` as a post-construction call

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Small observer pattern addition to existing struct
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: No UI work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 5, 6)
  - **Blocks**: Task 9 (GUI peer list needs live updates)
  - **Blocked By**: Task 1 (test infrastructure patterns)

  **References**:

  **Pattern References**:
  - `internal/zeroconf/zeroconf.go:58-66` — `Peers` struct with `mu *sync.RWMutex` and `peers map[string]*PeerInfo` — add `observers []func()` here
  - `internal/zeroconf/zeroconf.go:108-112` — `add()` method — call observers after unlock
  - `internal/zeroconf/zeroconf.go:114-118` — `remove()` method — call observers after unlock
  - `internal/zeroconf/zeroconf.go:144-158` — `Browse()` callback in `Start()` — this is where `add`/`remove` are called from mDNS events

  **Acceptance Criteria**:

  **TDD:**
  - [ ] Test file: `internal/zeroconf/zeroconf_test.go`
  - [ ] `go test ./internal/zeroconf/ -v` → PASS
  - [ ] Tests cover: observer called on add, observer called on remove, multiple observers, snapshot from All()

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Zeroconf tests pass
    Tool: Bash
    Preconditions: None
    Steps:
      1. Run: go test ./internal/zeroconf/ -v -count=1
      2. Assert: exit code 0
      3. Assert: output contains "TestPeersOnChangeCalledOnAdd"
      4. Assert: output contains "PASS"
    Expected Result: All observer tests pass
    Evidence: Test output captured
  ```

  **Commit**: YES
  - Message: `feat(zeroconf): add peer change observer callbacks to Peers`
  - Files: `internal/zeroconf/zeroconf.go`, `internal/zeroconf/zeroconf_test.go`
  - Pre-commit: `go test ./internal/zeroconf/ -v`

---

- [x] 5. Progress Tracking Wrappers

  **What to do**:
  - Create `internal/transport/progress.go`:
    - `ProgressFunc` type: `type ProgressFunc func(bytesTransferred int64, totalBytes int64)`
    - `ProgressWriter` struct wrapping `io.Writer` — calls `ProgressFunc` after each `Write()`
    - `ProgressReader` struct wrapping `io.Reader` — calls `ProgressFunc` after each `Read()`
    - `NewProgressWriter(w io.Writer, total int64, fn ProgressFunc) *ProgressWriter`
    - `NewProgressReader(r io.Reader, total int64, fn ProgressFunc) *ProgressReader`
    - Both track cumulative bytes and report via callback
  - Integrate into `sendFile()` and `storeFile()`:
    - `sendFile`: Wrap `conn` (writer) with `ProgressWriter` before `f.WriteTo()`
    - `storeFile`: Wrap `io.LimitReader` result with `ProgressReader` before `f.ReadFrom()`
    - Pass `ProgressFunc` as parameter to both functions (nil = no progress tracking, backwards compatible)
  - Create `internal/transport/progress_test.go` (TDD):
    - RED: `TestProgressWriterReportsBytesCorrectly` — write known data, verify callback values
    - RED: `TestProgressReaderReportsBytesCorrectly` — read known data, verify callback values
    - RED: `TestProgressWriterNilCallback` — nil func doesn't panic
    - RED: `TestProgressMonotonicallyIncreasing` — callbacks report strictly increasing byte counts
    - GREEN: Implement
    - REFACTOR: Clean up

  **Must NOT do**:
  - Do NOT add rate limiting, throttling, or speed calculation to the wrappers — keep them simple (bytes transferred + total)
  - Do NOT change the encryption layer (`secret/secret.go`) — progress wraps above encryption
  - Do NOT change function signatures of `HandleConnection` or `SendFile` (the exported ones) — only internal `sendFile`/`storeFile`

  **Recommended Agent Profile**:
  - **Category**: `unspecified-low`
    - Reason: Clean io.Writer/Reader wrapper pattern, well-understood in Go
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: No UI work

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 2 (with Tasks 4, 6)
  - **Blocks**: Tasks 12 (progress UI), 14 (integration)
  - **Blocked By**: Task 1 (test patterns)

  **References**:

  **Pattern References**:
  - `internal/transport/connection.go:114-124` — `sendFile()` uses `f.WriteTo(writer)` — wrap `writer` with `ProgressWriter`
  - `internal/transport/connection.go:88-112` — `storeFile()` uses `io.LimitReader(reader, size)` then `f.ReadFrom(lr)` — wrap `lr` with `ProgressReader`
  - `internal/transport/connection.go:102-103` — `bytes, err := f.ReadFrom(lr)` — `bytes` is already captured but ignored (`_ = bytes`); progress wrapper replaces this

  **External References**:
  - io.Writer/io.Reader interface: standard Go pattern for wrapping with middleware

  **Acceptance Criteria**:

  **TDD:**
  - [ ] Test file: `internal/transport/progress_test.go`
  - [ ] `go test ./internal/transport/ -v -run TestProgress` → PASS
  - [ ] Tests cover: writer reports correctly, reader reports correctly, nil callback safety, monotonic increase
  - [ ] Existing tests still pass: `go test ./internal/transport/ -v` → PASS (no regressions)

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Progress wrapper tests pass
    Tool: Bash
    Preconditions: Task 1 tests exist
    Steps:
      1. Run: go test ./internal/transport/ -v -count=1
      2. Assert: exit code 0
      3. Assert: output contains "TestProgressWriterReportsBytesCorrectly"
      4. Assert: output contains "TestProgressReaderReportsBytesCorrectly"
      5. Assert: all existing tests still pass (no "FAIL" in output)
    Expected Result: All tests pass, zero regressions
    Evidence: Test output captured
  ```

  **Commit**: YES
  - Message: `feat(transport): add progress tracking io.Writer/Reader wrappers`
  - Files: `internal/transport/progress.go`, `internal/transport/progress_test.go`, `internal/transport/connection.go`
  - Pre-commit: `go test ./internal/transport/ -v`

---

- [x] 6. Protocol Extension — BATCH_OFFER Message Type

- [x] 7. Gateway Interface Evolution

- [x] 8. Build Tag Restructuring

- [x] 9. GTK4 Panel Window — Peer List, Drag-and-Drop, File Picker

- [x] 10. DBus StatusNotifierItem — System Tray

- [x] 11. Desktop Notifications

- [x] 12. Transfer Progress UI

- [ ] 13. Incoming Transfer Dialog with Countdown

  **What to do**:
  - Implement `AskBatch()` in the GUI gateway with a GTK4 dialog:
    - `GtkDialog` (or `GtkWindow` as modal) with:
      - Header: "Incoming files from {peerName}"
      - File list: `GtkListBox` with filename + formatted size per entry
      - Total: "N files, X.XX MiB total"
      - Countdown bar: `GtkProgressBar` that decreases from 1.0 to 0.0 over timeout duration
      - Countdown text: "Auto-declining in Xs"
      - Two buttons: "Accept" (suggested/default) and "Decline"
    - Timer: Use `glib.TimeoutAdd(1000, ...)` for 1-second countdown ticks
    - On timeout expiry: auto-close dialog, return "DECLINE"
    - On "Accept" click: close dialog, return "ACCEPT"
    - On "Decline" click: close dialog, return "DECLINE"
    - On window close (X button): same as "Decline"
  - Implement `Ask()` for single-file offers (non-batch fallback):
    - Same dialog but with single file info parsed from the question string
  - Thread safety:
    - `Ask()`/`AskBatch()` are called from transport goroutines
    - Must schedule dialog creation on GTK main thread via `glib.IdleAdd`
    - Block the calling goroutine on a channel until dialog returns result
    - Pattern: send dialog request to channel → GTK main thread picks up → shows dialog → sends result back
  - Fire desktop notification (from Task 11) when dialog opens: "Incoming files from {peerName}"
  - Support multiple simultaneous dialogs (unlike TUI which queues/auto-declines)
  - Read timeout duration from config (Task 2)

  **Must NOT do**:
  - Do NOT add selective file acceptance checkboxes — all-or-nothing
  - Do NOT add file preview or type inspection
  - Do NOT queue dialogs (show all simultaneously) — let GTK handle window stacking

  **Recommended Agent Profile**:
  - **Category**: `visual-engineering`
    - Reason: GTK4 dialog design, countdown timer UX, thread-safe dialog lifecycle
  - **Skills**: [`frontend-ui-ux`]
    - `frontend-ui-ux`: Countdown timer UX, dialog layout, button emphasis
  - **Skills Evaluated but Omitted**:
    - `playwright`: Not browser-based

  **Parallelization**:
  - **Can Run In Parallel**: YES
  - **Parallel Group**: Wave 5 (with Task 12)
  - **Blocks**: Task 14 (integration)
  - **Blocked By**: Tasks 7 (BatchGateway interface), 9 (panel window), 11 (notifications)

  **References**:

  **Pattern References**:
  - `internal/platform/gateway_linux_tui.go:156-172` — TUI `Ask()` — channel-based prompt/response pattern; GUI uses same pattern but with GTK dialog
  - `internal/platform/gateway_linux_tui.go:20-23` — `promptRequest` struct (question + response channel) — reuse this pattern
  - `internal/platform/gateway_linux_tui.go:48-67` — Goroutine that handles prompts + notifications — GUI needs similar goroutine bridge
  - `internal/platform/gateway.go:14-23` — `BatchGateway` interface with `AskBatch(peerName string, files []FileInfo) string`
  - `internal/transport/connection.go:74-86` — `waitForDecision()` — calls `gw.Ask()` in goroutine with context timeout
  - `internal/transport/message.go:124-144` — `formatSize()` — use for file size display in dialog

  **External References**:
  - GTK4 Dialog: `gtk.NewDialog()` or custom `gtk.NewWindow()` as modal
  - GLib TimeoutAdd: `glib.TimeoutAdd(1000, func() bool)` — 1-second interval for countdown
  - GTK4 ProgressBar: `SetFraction()` for countdown visualization

  **Acceptance Criteria**:

  - [ ] `AskBatch()` shows GTK4 dialog with file list, sizes, countdown, Accept/Decline buttons
  - [ ] `Ask()` shows same dialog for single-file offers
  - [ ] Timeout auto-declines and closes dialog
  - [ ] Dialog is thread-safe (scheduled on GTK main thread from transport goroutine)
  - [ ] `go build -tags gui ./cmd/drift` → exit 0

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Dialog code compiles
    Tool: Bash
    Preconditions: Tasks 7, 9, 11 complete
    Steps:
      1. Run: go build -tags gui -o drift-gui ./cmd/drift
      2. Assert: exit code 0
    Expected Result: Incoming transfer dialog compiles
    Evidence: Build output captured
  ```

  **Commit**: YES
  - Message: `feat(platform): add incoming transfer dialog with countdown timer and batch file display`
  - Files: `internal/platform/gateway_linux_gui.go`
  - Pre-commit: `go build -tags gui ./cmd/drift`

---

- [ ] 14. Integration — Wire GUI Gateway to App + Config

  **What to do**:
  - Update `internal/app/run.go`:
    - Load config at startup: `cfg, err := config.Load(config.DefaultPath())`
    - Pass config to gateway constructor: update `platform.NewGateway` to accept config (or pass via separate method)
    - Use `cfg.DownloadDir` in transport layer (replace hardcoded `xdg.UserDirs.Download + "/Drift"`)
    - Use `cfg.AcceptTimeout` in `waitForDecision` timeout (replace hardcoded 30s)
    - Use `cfg.Identity` as identity override if non-empty
    - Wire progress callbacks:
      - Outbound: Pass `ProgressFunc` through to `SendBatch`/`SendFile` → update GUI transfer list
      - Inbound: Pass `ProgressFunc` through `HandleConnection` → update GUI transfer list
  - Update `internal/transport/connection.go`:
    - `HandleConnection` signature: add `downloadDir string` and `progressFn ProgressFunc` parameters
    - `SendFile`/`SendBatch`: add `progressFn ProgressFunc` parameter
    - Replace hardcoded `xdg.UserDirs.Download + "/Drift"` with `downloadDir` parameter
  - Update `internal/transport/connection.go`:
    - `waitForDecision`: accept timeout as parameter instead of hardcoded 30s
  - Ensure all existing non-GUI code paths still work:
    - TUI gateway gets nil progress callbacks (no-op)
    - Config defaults match current hardcoded values
  - Run full test suite and verify everything integrates

  **Must NOT do**:
  - Do NOT change the TUI user experience — config provides defaults that match current behavior
  - Do NOT add config UI to the GUI yet — config is file-based only
  - Do NOT break Windows or macOS compilation

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Wiring multiple subsystems together, ensuring no regressions across 5+ modified files
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `frontend-ui-ux`: Integration work, not UI design

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 6 (sequential)
  - **Blocks**: Task 15
  - **Blocked By**: Tasks 2, 3, 5, 10, 12, 13

  **References**:

  **Pattern References**:
  - `internal/app/run.go:18-159` — Entire orchestration function — this is what gets modified
  - `internal/app/run.go:46-47` — Gateway instantiation — needs config parameter
  - `internal/app/run.go:110-151` — Outbound processor — needs progress callback wiring
  - `internal/app/run.go:66-108` — Inbound processor — needs progress callback + download dir
  - `internal/transport/connection.go:16-72` — `HandleConnection` — needs downloadDir + progressFn params
  - `internal/transport/connection.go:74-86` — `waitForDecision` — needs configurable timeout
  - `internal/transport/connection.go:49` — Hardcoded download path: `filepath.Join(xdg.UserDirs.Download, "Drift")`
  - `internal/config/config.go` — Config system from Task 2
  - `internal/transport/progress.go` — Progress wrappers from Task 5

  **Acceptance Criteria**:

  - [ ] `go build ./cmd/drift` → exit 0 (TUI with config support)
  - [ ] `go build -tags gui ./cmd/drift` → exit 0 (GUI fully integrated)
  - [ ] `go test ./... -count=1` → PASS (all tests, no regressions)
  - [ ] `go vet ./...` → no errors
  - [ ] `GOOS=windows go vet ./internal/...` → no errors (cross-platform check)
  - [ ] Config defaults match previous hardcoded values (download dir, 30s timeout)

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Full test suite passes after integration
    Tool: Bash
    Preconditions: All tasks 1-13 complete
    Steps:
      1. Run: go test ./... -v -count=1
      2. Assert: exit code 0
      3. Assert: no "FAIL" in output
      4. Count test cases: grep -c "--- PASS" in output
      5. Assert: count >= 20
    Expected Result: All tests pass, no regressions
    Evidence: Full test output captured

  Scenario: TUI builds and runs with config
    Tool: Bash
    Preconditions: Integration complete
    Steps:
      1. Run: go build -o drift-tui ./cmd/drift
      2. Assert: exit code 0
      3. Run: timeout 3 ./drift-tui 2>&1 || true
      4. Assert: output contains "Drift running on Linux" (TUI banner)
      5. Run: rm drift-tui
    Expected Result: TUI works with config integration
    Evidence: Output captured

  Scenario: GUI builds with all features
    Tool: Bash
    Preconditions: Integration complete
    Steps:
      1. Run: go build -tags gui -o drift-gui ./cmd/drift
      2. Assert: exit code 0
      3. Run: ldd drift-gui | grep "libgtk-4"
      4. Assert: GTK4 linked
      5. Run: rm drift-gui
    Expected Result: GUI binary includes all features
    Evidence: Build and ldd output captured

  Scenario: Cross-platform compilation check
    Tool: Bash
    Preconditions: Integration complete
    Steps:
      1. Run: GOOS=windows go vet ./internal/...
      2. Assert: exit code 0
      3. Run: GOOS=darwin go vet ./internal/...
      4. Assert: exit code 0 (may have darwinkit warnings, but no new errors)
    Expected Result: No cross-platform breakage from GUI changes
    Evidence: vet output captured
  ```

  **Commit**: YES
  - Message: `feat(app): integrate config, progress tracking, and batch transfers into app lifecycle`
  - Files: `internal/app/run.go`, `internal/transport/connection.go`
  - Pre-commit: `go test ./... -count=1 && go vet ./...`

---

- [ ] 15. Final Verification & Cross-Platform Build Check

  **What to do**:
  - Run complete verification suite:
    - `go test ./... -v -count=1` — all tests pass
    - `go vet ./...` — no issues
    - `go vet -tags gui ./...` — no issues with GUI code
    - `go build ./cmd/drift` — TUI builds
    - `go build -tags gui ./cmd/drift` — GUI builds
    - `GOOS=windows go build ./cmd/drift` — Windows still builds
  - Verify `go mod tidy` leaves no unused dependencies
  - Run `go mod vendor` to update vendor directory with new deps (gotk4, godbus, toml)
  - Verify file structure matches plan:
    - `internal/config/config.go` + `config_test.go` exist
    - `internal/transport/progress.go` + `progress_test.go` exist
    - `internal/transport/message_test.go` + `connection_test.go` exist
    - `internal/platform/gateway_linux_tui.go` exists (renamed from gateway_linux.go)
    - `internal/platform/gateway_linux_gui.go` exists (new)
    - `internal/platform/tray_linux.go` exists (new)
    - `internal/platform/notify_linux.go` exists (new)
    - `internal/platform/assets/drift-icon.svg` exists (new)
    - `internal/zeroconf/zeroconf_test.go` exists (new)
  - Smoke test: start GUI binary, verify it doesn't crash within 5 seconds

  **Must NOT do**:
  - Do NOT add new features in this task
  - Do NOT refactor or clean up — this is verification only
  - Do NOT push to remote (user decides when)

  **Recommended Agent Profile**:
  - **Category**: `deep`
    - Reason: Comprehensive verification across all deliverables
  - **Skills**: []
  - **Skills Evaluated but Omitted**:
    - `playwright`: No browser testing needed

  **Parallelization**:
  - **Can Run In Parallel**: NO
  - **Parallel Group**: Wave 6 (after Task 14)
  - **Blocks**: None (final task)
  - **Blocked By**: Task 14

  **References**:

  **Pattern References**:
  - All files created/modified in Tasks 1-14

  **Acceptance Criteria**:

  - [ ] `go test ./... -v -count=1` → PASS, 0 failures
  - [ ] `go vet ./...` → exit 0
  - [ ] `go vet -tags gui ./...` → exit 0
  - [ ] `go build ./cmd/drift` → exit 0
  - [ ] `go build -tags gui ./cmd/drift` → exit 0
  - [ ] `GOOS=windows go build ./cmd/drift` → exit 0
  - [ ] `go mod tidy` → no changes (clean deps)
  - [ ] All expected files exist (see list above)
  - [ ] Vendor directory updated with new deps

  **Agent-Executed QA Scenarios:**

  ```
  Scenario: Complete build and test verification
    Tool: Bash
    Preconditions: All tasks 1-14 complete
    Steps:
      1. Run: go test ./... -v -count=1
      2. Assert: exit code 0, no FAIL
      3. Run: go vet ./...
      4. Assert: exit code 0
      5. Run: go vet -tags gui ./...
      6. Assert: exit code 0
      7. Run: go build -o drift-tui ./cmd/drift
      8. Assert: exit code 0
      9. Run: go build -tags gui -o drift-gui ./cmd/drift
      10. Assert: exit code 0
      11. Run: GOOS=windows go build -o drift.exe ./cmd/drift
      12. Assert: exit code 0
      13. Clean up: rm drift-tui drift-gui drift.exe
    Expected Result: Everything builds and passes
    Evidence: All outputs captured

  Scenario: File structure verification
    Tool: Bash
    Preconditions: All tasks complete
    Steps:
      1. Verify: test -f internal/config/config.go
      2. Verify: test -f internal/config/config_test.go
      3. Verify: test -f internal/transport/progress.go
      4. Verify: test -f internal/transport/progress_test.go
      5. Verify: test -f internal/transport/message_test.go
      6. Verify: test -f internal/platform/gateway_linux_tui.go
      7. Verify: test -f internal/platform/gateway_linux_gui.go
      8. Verify: test -f internal/platform/tray_linux.go
      9. Verify: test -f internal/platform/notify_linux.go
      10. Verify: test -f internal/platform/assets/drift-icon.svg
      11. Verify: test -f internal/zeroconf/zeroconf_test.go
      12. Assert: all files exist (exit code 0 for each)
    Expected Result: All expected deliverables present
    Evidence: File existence checks captured

  Scenario: GUI smoke test
    Tool: interactive_bash (tmux)
    Preconditions: drift-gui binary built, desktop session available
    Steps:
      1. Start: ./drift-gui &
      2. Wait: 5 seconds
      3. Check: ps aux | grep drift-gui | grep -v grep
      4. Assert: process is still running (didn't crash)
      5. Kill: kill %1
    Expected Result: GUI starts and stays alive for 5+ seconds
    Evidence: Process list captured
  ```

  **Commit**: YES
  - Message: `chore: vendor new dependencies and verify complete build`
  - Files: `go.mod`, `go.sum`, `vendor/`
  - Pre-commit: `go test ./... -count=1 && go vet ./... && go build ./cmd/drift && go build -tags gui ./cmd/drift`

---

## Commit Strategy

| After Task | Message | Key Files | Verification |
|------------|---------|-----------|--------------|
| 1 | `test(transport): add baseline tests for message marshal/unmarshal and file storage` | `*_test.go` | `go test ./internal/transport/` |
| 2 | `feat(config): add XDG-compliant TOML config system` | `internal/config/*` | `go test ./internal/config/` |
| 3 | `fix(app): prevent outbound processor goroutine death on single transfer failure` | `run.go`, `connection.go` | `go vet ./...` |
| 4 | `feat(zeroconf): add peer change observer callbacks to Peers` | `zeroconf.go`, `*_test.go` | `go test ./internal/zeroconf/` |
| 5 | `feat(transport): add progress tracking io.Writer/Reader wrappers` | `progress.go`, `*_test.go` | `go test ./internal/transport/` |
| 6 | `feat(transport): add BATCH_OFFER protocol extension for multi-file transfers` | `message.go`, `connection.go` | `go test ./internal/transport/` |
| 7 | `feat(platform): add BatchGateway interface and multi-file Request support` | `gateway.go`, all gateway_*.go | `go vet ./...` |
| 8 | `refactor(platform): split Linux gateway into TUI and GUI build targets` | `gateway_linux_*.go` | `go build` both targets |
| 9 | `feat(platform): implement GTK4 panel window with peer list and drag-and-drop` | `gateway_linux_gui.go` | `go build -tags gui` |
| 10 | `feat(platform): implement DBus StatusNotifierItem system tray` | `tray_linux.go` | `go build -tags gui` |
| 11 | `feat(platform): add desktop notifications via freedesktop DBus` | `notify_linux.go` | `go build -tags gui` |
| 12 | `feat(platform): add transfer progress tracking UI` | `gateway_linux_gui.go` | `go build -tags gui` |
| 13 | `feat(platform): add incoming transfer dialog with countdown timer` | `gateway_linux_gui.go` | `go build -tags gui` |
| 14 | `feat(app): integrate config, progress, and batch transfers into app lifecycle` | `run.go`, `connection.go` | `go test ./...` |
| 15 | `chore: vendor new dependencies and verify complete build` | `vendor/`, `go.mod` | Full suite |

---

## Success Criteria

### Verification Commands
```bash
# All tests pass
go test ./... -v -count=1  # Expected: PASS, 0 failures

# Code quality
go vet ./...               # Expected: exit 0
go vet -tags gui ./...     # Expected: exit 0

# TUI builds (default, backwards compatible)
go build ./cmd/drift       # Expected: exit 0

# GUI builds (opt-in)
go build -tags gui ./cmd/drift  # Expected: exit 0

# Cross-platform still works
GOOS=windows go build ./cmd/drift  # Expected: exit 0

# Dependencies clean
go mod tidy                # Expected: no changes
```

### Final Checklist
- [ ] All "Must Have" items present and functional
- [ ] All "Must NOT Have" guardrails respected — no scope creep
- [ ] All tests pass (`go test ./...`)
- [ ] TUI unmodified in behavior (default build)
- [ ] GUI compiles and starts without crash
- [ ] BATCH_OFFER wire format documented and tested
- [ ] Config file optional — zero-config still works
- [ ] No gotk4 imports leak into non-GUI build paths
- [ ] No regressions on Windows/macOS compilation
