import Foundation

extension Data {
    /// Initialize Data from a hex-encoded string.
    init?(hexString: String) {
        let len = hexString.count
        guard len % 2 == 0 else { return nil }

        var data = Data(capacity: len / 2)
        var index = hexString.startIndex
        for _ in 0..<len / 2 {
            let nextIndex = hexString.index(index, offsetBy: 2)
            guard let byte = UInt8(hexString[index..<nextIndex], radix: 16) else {
                return nil
            }
            data.append(byte)
            index = nextIndex
        }
        self = data
    }

    /// Lowercase hex string representation.
    var hexString: String {
        map { String(format: "%02x", $0) }.joined()
    }
}
