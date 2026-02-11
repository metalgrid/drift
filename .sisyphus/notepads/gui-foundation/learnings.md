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
