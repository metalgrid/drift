# Security Hardening Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Eliminate high-risk protocol, file-write, and crypto safety issues in Go and iOS transfer paths while keeping wire compatibility for current peers.

**Architecture:** Add strict input validation and bounded parsing at protocol edges, then harden transport crypto/frame handling and socket lifecycle behavior. Keep message format stable for this phase; stage larger cryptographic redesign (authenticated key exchange and identity persistence) as a follow-up track.

**Tech Stack:** Go (net, io, filepath, crypto), Swift (CryptoKit, Network, Foundation), existing Drift protocol tests.

---

### Task 1: Go inbound filename and path traversal hardening

**Files:**
- Modify: `internal/transport/connection.go`
- Test: `internal/transport/connection_test.go`

**Step 1: Write failing tests**

Add tests for `storeFile` rejecting:
- absolute paths
- `..` traversal
- filename containing `/` or `\\`

```go
func TestStoreFileRejectsTraversal(t *testing.T) {
    err := storeFile(tmpDir, "../../evil.txt", 1, bytes.NewReader([]byte("x")), nil)
    require.Error(t, err)
}
```

**Step 2: Run tests to verify fail**

Run: `go test ./internal/transport -run StoreFile -v`
Expected: FAIL on traversal acceptance.

**Step 3: Minimal implementation**

Introduce filename sanitizer:
- require `name == filepath.Base(name)`
- reject empty or `.` or `..`
- reject any path separator
- resolve final target path and ensure it remains under incoming root

**Step 4: Re-run tests**

Run: `go test ./internal/transport -run StoreFile -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/transport/connection.go internal/transport/connection_test.go
git commit -m "fix(transport): prevent path traversal in inbound filenames"
```

### Task 2: Go buffered-stream correctness and control-frame limits

**Files:**
- Modify: `internal/transport/connection.go`
- Test: `internal/transport/connection_test.go`

**Step 1: Write failing tests**

Add tests for:
- control reader over-read does not corrupt file payload
- oversized control line is rejected

**Step 2: Run tests to verify fail**

Run: `go test ./internal/transport -run "(Buffered|Control)" -v`
Expected: FAIL.

**Step 3: Minimal implementation**

- Pass `reader` (bufio.Reader) to `storeFile` instead of raw `conn`.
- Replace unbounded `ReadString('\n')` with bounded read helper that errors at max control frame size.

**Step 4: Re-run tests**

Run: `go test ./internal/transport -run "(Buffered|Control)" -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/transport/connection.go internal/transport/connection_test.go
git commit -m "fix(transport): enforce bounded control frames and buffered payload reads"
```

### Task 3: Go decrypt reader safety (bounded frame + io.Reader contract)

**Files:**
- Modify: `internal/secret/secret.go`
- Test: `internal/secret/secret_test.go` (create)

**Step 1: Write failing tests**

Add tests for:
- short destination buffer returns partial bytes and subsequent reads continue
- oversized frame length is rejected before allocation

**Step 2: Run tests to verify fail**

Run: `go test ./internal/secret -v`
Expected: FAIL.

**Step 3: Minimal implementation**

- Add `maxEncryptedFrameSize` constant.
- Add `decBuf` draining behavior before reading next frame.
- Ensure `Read` returns `n <= len(buf)` and preserves remainder.

**Step 4: Re-run tests**

Run: `go test ./internal/secret -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/secret/secret.go internal/secret/secret_test.go
git commit -m "fix(secret): bound frame size and honor io.Reader semantics"
```

### Task 4: Go peer key and address validation in app loop

**Files:**
- Modify: `internal/app/run.go`
- Test: `internal/app/run_test.go` (create focused helpers tests)

**Step 1: Write failing tests**

Add tests for helper validation logic:
- peer public key must decode to exactly 32 bytes
- peer must provide at least one address

**Step 2: Run tests to verify fail**

Run: `go test ./internal/app -v`
Expected: FAIL.

**Step 3: Minimal implementation**

- Add helper validator(s) and use before `copy` and `peer.Addresses[0]`.
- On invalid key/address, notify and continue (no panic).

**Step 4: Re-run tests**

Run: `go test ./internal/app -v`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/app/run.go internal/app/run_test.go
git commit -m "fix(app): validate peer key length and address presence"
```

### Task 5: iOS filename sanitization and path containment

**Files:**
- Modify: `ios/Drift/Utilities/FileStorage.swift`
- Modify: `ios/Drift/Networking/ConnectionHandler.swift`
- Test: `ios/DriftTests/MessageTests.swift` (message-level) and new storage tests

**Step 1: Write failing tests**

Add tests to reject traversal or separator-containing filenames for receive path helper.

**Step 2: Run tests to verify fail**

Run: `xcodebuild test -project ios/Drift.xcodeproj -scheme DriftTests -destination 'platform=iOS Simulator,name=iPhone 15'`
Expected: FAIL.

**Step 3: Minimal implementation**

- Add canonical filename sanitizer in `FileStorage`.
- Resolve and compare standardized paths to ensure destination remains in Drift directory.

**Step 4: Re-run tests**

Expected: PASS.

**Step 5: Commit**

```bash
git add ios/Drift/Utilities/FileStorage.swift ios/Drift/Networking/ConnectionHandler.swift ios/DriftTests
git commit -m "fix(ios): block path traversal in received filenames"
```

### Task 6: iOS stream/framing robustness

**Files:**
- Modify: `ios/Drift/Networking/SecureConnection.swift`
- Modify: `ios/Drift/Crypto/EncryptedStream.swift`
- Test: `ios/DriftTests/EncryptedStreamTests.swift`

**Step 1: Write failing tests**

Add tests for:
- partial receive chunks still satisfy exact reads
- oversized encrypted frame length is rejected

**Step 2: Run tests to verify fail**

Run iOS tests.

**Step 3: Minimal implementation**

- Implement internal read accumulator loop in `readExactly`.
- Add max encrypted frame constant and reject oversized lengths.

**Step 4: Re-run tests**

Expected: PASS.

**Step 5: Commit**

```bash
git add ios/Drift/Networking/SecureConnection.swift ios/Drift/Crypto/EncryptedStream.swift ios/DriftTests/EncryptedStreamTests.swift
git commit -m "fix(ios): handle partial TCP reads and bound frame size"
```

### Task 7: iOS outbound stability and UI state safety

**Files:**
- Modify: `ios/Drift/AppCoordinator.swift`
- Modify: `ios/Drift/UI/PeerListView.swift`

**Step 1: Write failing tests (or deterministic harness checks)**

Cover:
- invalid port does not crash
- continuation resumed once
- sheet dismissal clears pending incoming offer

**Step 2: Run tests**

Run iOS tests/harness.

**Step 3: Minimal implementation**

- Guard `NWEndpoint.Port(rawValue:)`.
- One-shot continuation wrapper in connect wait path.
- Handle sheet setter dismiss case by declining/cleanup.

**Step 4: Re-run tests**

Expected: PASS.

**Step 5: Commit**

```bash
git add ios/Drift/AppCoordinator.swift ios/Drift/UI/PeerListView.swift
git commit -m "fix(ios): harden outbound connect state and incoming-offer dismissal"
```

### Task 8: Protocol hygiene, timeouts, and docs

**Files:**
- Modify: `internal/transport/message.go`
- Modify: `ios/Drift/Protocol/Message.swift`
- Modify: `internal/transport/connection.go`
- Add: `docs/plans/2026-02-26-crypto-next-phase.md`

**Step 1: Write failing tests**

Add tests rejecting protocol separator/newline in filename and enforcing read deadlines where applicable.

**Step 2: Run tests to verify fail**

Go test targeted packages.

**Step 3: Minimal implementation**

- Validate filename fields used in protocol serialization.
- Add connection deadlines/timeouts for handshake/control read phases.
- Document next-phase crypto migration (HKDF context binding, stable identity keys, authenticated key exchange).

**Step 4: Re-run tests**

Run: `go test ./...`
Expected: PASS.

**Step 5: Commit**

```bash
git add internal/transport/message.go ios/Drift/Protocol/Message.swift internal/transport/connection.go docs/plans/2026-02-26-crypto-next-phase.md
git commit -m "chore(security): enforce protocol field validation and document crypto migration"
```
