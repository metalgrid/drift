import Foundation

/// 12-byte sequential nonce matching Go's incrementNonce behavior.
/// Starts at all zeros and increments big-endian from the rightmost byte.
struct NonceCounter {
    private(set) var bytes: [UInt8]

    init() {
        bytes = [UInt8](repeating: 0, count: 12)
    }

    /// Increment the nonce using big-endian right-to-left carry.
    /// Matches Go: for i := len(nonce)-1; i >= 0; i-- { nonce[i]++; if nonce[i] != 0 { break } }
    mutating func increment() {
        for i in stride(from: bytes.count - 1, through: 0, by: -1) {
            bytes[i] &+= 1
            if bytes[i] != 0 { break }
        }
    }

    var data: Data {
        Data(bytes)
    }
}
