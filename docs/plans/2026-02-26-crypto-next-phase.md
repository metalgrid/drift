# Crypto Next Phase Plan

## Goal

Replace the current opportunistic encrypted channel with an authenticated key exchange design that provides identity verification, forward secrecy, and explicit key separation.

## Problems In Current Design

- Peer key material is discovered via mDNS TXT and is not authenticated.
- Runtime-generated long-term keys prevent trust-on-first-use and stable identity pinning.
- KDF is a direct hash of ECDH output, without protocol context binding.

## Target Properties

- Authenticated handshake (active MITM resistance).
- Forward secrecy for captured traffic.
- Stable device identity keys (persisted, rotated deliberately).
- Session key derivation with role/context separation.
- Versioned handshake transcript and downgrade protection.

## Recommended Design

1. Persist a long-term static keypair per device.
2. Adopt a Noise-based handshake profile:
   - `XX` for first-contact workflows.
   - `IK` for known-peer optimized reconnects.
3. Derive traffic keys with HKDF labels:
   - `drift/v1/handshake`
   - `drift/v1/tx`
   - `drift/v1/rx`
4. Bind advertised peer identity to fingerprint verification UI (TOFU + pinning).
5. Reject protocol version downgrades at handshake time.

## Migration Strategy

1. Introduce persistent identity key storage in Go + iOS.
2. Add handshake version negotiation field (`v=0.2`) while keeping current transport for `v=0.1` peers.
3. Ship dual-stack handshake support (`legacy` + `noise-v1`).
4. Add peer trust store and explicit fingerprint confirmation for new peers.
5. Deprecate legacy mode after a compatibility window.

## Verification

- Unit tests: transcript validation, HKDF label separation, replay rejection.
- Integration tests: Go<->Go, iOS<->iOS, Go<->iOS compatibility matrix.
- Security tests: MITM simulation, downgrade attempt, stale replayed frames.
