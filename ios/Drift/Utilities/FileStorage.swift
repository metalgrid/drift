import Foundation

enum FileStorageError: Error {
    case invalidFilename
    case pathTraversalDetected
}

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

    func destinationURLs(for unsafeFilename: String) throws -> (temp: URL, final: URL) {
        let filename = try sanitizeFilename(unsafeFilename)
        let finalURL = basePath.appendingPathComponent(filename)
        let tempURL = basePath.appendingPathComponent("\(filename).\(UUID().uuidString).drift")

        let baseStandardized = basePath.standardizedFileURL
        let finalStandardized = finalURL.standardizedFileURL
        let tempStandardized = tempURL.standardizedFileURL

        let basePathString = baseStandardized.path
        guard finalStandardized.path.hasPrefix(basePathString + "/"),
              tempStandardized.path.hasPrefix(basePathString + "/") else {
            throw FileStorageError.pathTraversalDetected
        }

        return (temp: tempStandardized, final: finalStandardized)
    }

    private func sanitizeFilename(_ filename: String) throws -> String {
        let trimmed = filename.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !trimmed.isEmpty else {
            throw FileStorageError.invalidFilename
        }

        let base = URL(fileURLWithPath: trimmed).lastPathComponent
        guard base == trimmed,
              !base.contains("/"),
              !base.contains("\\"),
              base != ".",
              base != ".." else {
            throw FileStorageError.invalidFilename
        }

        return base
    }
}
