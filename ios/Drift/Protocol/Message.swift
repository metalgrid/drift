import Foundation

private let fieldSeparator = "|"
private let endOfMessage: Character = "\n"
private let mimeType = "application/octet-stream"

/// Wire protocol messages matching Go's transport/message.go.
enum DriftMessage {
    case offer(filename: String, mimetype: String, size: Int64)
    case batchOffer(files: [(filename: String, mimetype: String, size: Int64)])
    case answer(kind: String)

    // MARK: - Serialize

    /// Serialize to wire format: pipe-delimited, newline-terminated.
    /// Matches Go's MarshalMessage().
    func serialize() -> Data {
        let str: String
        switch self {
        case .offer(let filename, let mimetype, let size):
            str = ["OFFER", filename, mimetype, String(size)].joined(separator: fieldSeparator)
        case .batchOffer(let files):
            var parts = ["BATCH_OFFER", String(files.count)]
            for file in files {
                parts.append(contentsOf: [file.filename, file.mimetype, String(file.size)])
            }
            str = parts.joined(separator: fieldSeparator)
        case .answer(let kind):
            str = ["ANSWER", kind].joined(separator: fieldSeparator)
        }
        return Data((str + String(endOfMessage)).utf8)
    }

    // MARK: - Parse

    /// Parse a raw message string. Checks BATCH_OFFER before OFFER to avoid prefix collision.
    /// Matches Go's UnmarshalMessage().
    static func parse(_ raw: String) -> DriftMessage? {
        let msg = raw.hasSuffix(String(endOfMessage))
            ? String(raw.dropLast())
            : raw

        if msg.hasPrefix("BATCH_OFFER") {
            let parts = msg.split(separator: Character(fieldSeparator), omittingEmptySubsequences: false).map(String.init)
            guard parts.count >= 2 else { return nil }
            guard let count = Int(parts[1]), count > 0 else { return nil }
            let expectedParts = 2 + (count * 3)
            guard parts.count == expectedParts else { return nil }

            var files: [(filename: String, mimetype: String, size: Int64)] = []
            for i in 0..<count {
                let idx = 2 + (i * 3)
                guard let size = Int64(parts[idx + 2]) else { return nil }
                files.append((filename: parts[idx], mimetype: parts[idx + 1], size: size))
            }
            return .batchOffer(files: files)
        }

        if msg.hasPrefix("OFFER") {
            let parts = msg.split(separator: Character(fieldSeparator), omittingEmptySubsequences: false).map(String.init)
            guard parts.count == 4 else { return nil }
            guard let size = Int64(parts[3]) else { return nil }
            return .offer(filename: parts[1], mimetype: parts[2], size: size)
        }

        if msg.hasPrefix("ANSWER") {
            let parts = msg.split(separator: Character(fieldSeparator), omittingEmptySubsequences: false).map(String.init)
            guard parts.count == 2 else { return nil }
            return .answer(kind: parts[1])
        }

        return nil
    }

    // MARK: - Convenience constructors

    static func makeOffer(filename: String, size: Int64) -> DriftMessage {
        .offer(filename: filename, mimetype: mimeType, size: size)
    }

    static func makeBatchOffer(files: [(filename: String, size: Int64)]) -> DriftMessage {
        .batchOffer(files: files.map { ($0.filename, mimeType, $0.size) })
    }

    static func accept() -> DriftMessage {
        .answer(kind: "ACCEPT")
    }

    static func decline() -> DriftMessage {
        .answer(kind: "DECLINE")
    }

    var isAccepted: Bool {
        if case .answer(let kind) = self {
            return kind == "ACCEPT"
        }
        return false
    }
}
