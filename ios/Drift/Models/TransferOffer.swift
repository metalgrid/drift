import Foundation

struct FileEntry: Identifiable {
    let id = UUID()
    let filename: String
    let mimetype: String
    let size: Int64

    var formattedSize: String {
        Self.formatSize(size)
    }

    static func formatSize(_ size: Int64) -> String {
        let kib: Int64 = 1024
        let mib = kib * 1024
        let gib = mib * 1024
        let tib = gib * 1024

        switch size {
        case tib...:
            return String(format: "%.2f TiB", Double(size) / Double(tib))
        case gib...:
            return String(format: "%.2f GiB", Double(size) / Double(gib))
        case mib...:
            return String(format: "%.2f MiB", Double(size) / Double(mib))
        case kib...:
            return String(format: "%.2f KiB", Double(size) / Double(kib))
        default:
            return "\(size) Bytes"
        }
    }
}

struct TransferOffer: Identifiable {
    let id = UUID()
    let peerName: String
    let files: [FileEntry]
    let receivedAt: Date

    var totalSize: Int64 {
        files.reduce(0) { $0 + $1.size }
    }

    var formattedTotalSize: String {
        FileEntry.formatSize(totalSize)
    }

    var isBatch: Bool {
        files.count > 1
    }
}
