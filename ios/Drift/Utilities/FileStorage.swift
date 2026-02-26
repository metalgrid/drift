import Foundation

/// Manages file storage for received transfers.
/// Files are saved to Documents/Drift, matching Go's ~/Downloads/Drift.
struct FileStorage {
    private let basePath: URL

    init() {
        let docs = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask).first!
        basePath = docs.appendingPathComponent("Drift", isDirectory: true)
    }

    func ensureDirectory() throws {
        try FileManager.default.createDirectory(at: basePath, withIntermediateDirectories: true)
    }

    func fileURL(for filename: String) -> URL {
        basePath.appendingPathComponent(filename)
    }

    func tempFileURL(for filename: String) -> URL {
        basePath.appendingPathComponent(filename + ".\(UUID().uuidString).drift")
    }
}
