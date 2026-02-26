import Foundation

/// Callback types for UI interaction during connection handling.
typealias OfferPrompt = (TransferOffer) async -> Bool
typealias ProgressCallback = (Int64, Int64) -> Void

/// Handles the message loop on a secure connection.
/// Matches Go's transport.HandleConnection().
final class ConnectionHandler {
    private static let maxControlBufferBytes = 64 * 1024
    private let connection: SecureConnection
    private let onOffer: OfferPrompt
    private let onProgress: ProgressCallback?
    private let onComplete: (String) -> Void
    private let onError: (String) -> Void
    private var buffer = Data()

    init(
        connection: SecureConnection,
        onOffer: @escaping OfferPrompt,
        onProgress: ProgressCallback? = nil,
        onComplete: @escaping (String) -> Void,
        onError: @escaping (String) -> Void
    ) {
        self.connection = connection
        self.onOffer = onOffer
        self.onProgress = onProgress
        self.onComplete = onComplete
        self.onError = onError
    }

    /// Run the message loop: read decrypted chunks, buffer until newline, parse, dispatch.
    func run(peerName: String) async {
        do {
            while true {
                let chunk = try await connection.readDecrypted()
                buffer.append(chunk)

                if buffer.count > Self.maxControlBufferBytes {
                    throw CryptoError.readFailed
                }

                while let newlineIndex = buffer.firstIndex(of: UInt8(ascii: "\n")) {
                    let messageData = buffer[buffer.startIndex...newlineIndex]
                    buffer = Data(buffer[(newlineIndex + 1)...])

                    guard let messageStr = String(data: Data(messageData), encoding: .utf8),
                          let message = DriftMessage.parse(messageStr) else {
                        continue
                    }

                    try await handleMessage(message, peerName: peerName)
                }
            }
        } catch {
            // Connection closed or error - expected on disconnect
            print("Connection handler ended: \(error)")
        }
    }

    private func handleMessage(_ message: DriftMessage, peerName: String) async throws {
        switch message {
        case .offer(let filename, _, let size):
            let offer = TransferOffer(
                peerName: peerName,
                files: [FileEntry(filename: filename, mimetype: "application/octet-stream", size: size)],
                receivedAt: Date()
            )

            let accepted = await onOffer(offer)
            if accepted {
                try await connection.writeEncrypted(DriftMessage.accept().serialize())
                try await receiveFile(filename: filename, size: size)
                onComplete("File received: \(filename)")
            } else {
                try await connection.writeEncrypted(DriftMessage.decline().serialize())
            }

        case .batchOffer(let files):
            let entries = files.map { FileEntry(filename: $0.filename, mimetype: $0.mimetype, size: $0.size) }
            let offer = TransferOffer(peerName: peerName, files: entries, receivedAt: Date())

            let accepted = await onOffer(offer)
            if accepted {
                try await connection.writeEncrypted(DriftMessage.accept().serialize())
                for file in files {
                    try await receiveFile(filename: file.filename, size: file.size)
                }
                onComplete("Batch received: \(files.count) files")
            } else {
                try await connection.writeEncrypted(DriftMessage.decline().serialize())
            }

        case .answer(let kind):
            if kind != "ACCEPT" {
                print("Transfer declined by peer")
            }
        }
    }

    /// Receive a file through the encrypted stream using a size-limited read.
    /// Matches Go's storeFile with LimitReader.
    private func receiveFile(filename: String, size: Int64) async throws {
        let storage = FileStorage()
        let destination = try storage.destinationURLs(for: filename)
        let tempURL = destination.temp
        let finalURL = destination.final

        try storage.ensureDirectory()

        FileManager.default.createFile(atPath: tempURL.path, contents: nil)
        let handle = try FileHandle(forWritingTo: tempURL)
        defer { try? handle.close() }

        var bytesReceived: Int64 = 0
        while bytesReceived < size {
            let chunk = try await connection.readDecrypted()
            let remaining = size - bytesReceived
            let toWrite = Int64(chunk.count) <= remaining ? chunk : chunk.prefix(Int(remaining))

            handle.write(toWrite)
            bytesReceived += Int64(toWrite.count)
            onProgress?(bytesReceived, size)
        }

        try? handle.close()

        // Atomic rename like Go's os.Rename
        let fm = FileManager.default
        if fm.fileExists(atPath: finalURL.path) {
            try fm.removeItem(at: finalURL)
        }
        try fm.moveItem(at: tempURL, to: finalURL)
    }

    /// Send files through the encrypted stream.
    func sendFiles(_ fileURLs: [URL]) async throws {
        if fileURLs.count > 1 {
            var fileEntries: [(filename: String, size: Int64)] = []
            for url in fileURLs {
                let filename = url.lastPathComponent
                guard DriftMessage.isValidProtocolFilename(filename) else {
                    throw CryptoError.writeFailed
                }
                let attrs = try FileManager.default.attributesOfItem(atPath: url.path)
                let size = (attrs[.size] as? Int64) ?? 0
                fileEntries.append((filename: filename, size: size))
            }

            let offer = DriftMessage.makeBatchOffer(files: fileEntries)
            try await connection.writeEncrypted(offer.serialize())

            // Wait for answer
            let answerData = try await connection.readDecrypted()
            guard let answerStr = String(data: answerData, encoding: .utf8),
                  let answer = DriftMessage.parse(answerStr),
                  answer.isAccepted else {
                onError("Transfer declined")
                return
            }

            for url in fileURLs {
                try await sendFileData(url)
            }
            onComplete("Batch sent: \(fileURLs.count) files")
        } else if let url = fileURLs.first {
            guard DriftMessage.isValidProtocolFilename(url.lastPathComponent) else {
                throw CryptoError.writeFailed
            }
            let attrs = try FileManager.default.attributesOfItem(atPath: url.path)
            let size = (attrs[.size] as? Int64) ?? 0

            let offer = DriftMessage.makeOffer(filename: url.lastPathComponent, size: size)
            try await connection.writeEncrypted(offer.serialize())

            // Wait for answer
            let answerData = try await connection.readDecrypted()
            guard let answerStr = String(data: answerData, encoding: .utf8),
                  let answer = DriftMessage.parse(answerStr),
                  answer.isAccepted else {
                onError("Transfer declined")
                return
            }

            try await sendFileData(url)
            onComplete("File sent: \(url.lastPathComponent)")
        }
    }

    /// Stream file data in ~32KB chunks through the encrypted writer.
    private func sendFileData(_ url: URL) async throws {
        let handle = try FileHandle(forReadingFrom: url)
        defer { try? handle.close() }

        let attrs = try FileManager.default.attributesOfItem(atPath: url.path)
        let totalSize = (attrs[.size] as? Int64) ?? 0
        var bytesSent: Int64 = 0

        let chunkSize = 32 * 1024
        while true {
            let chunk = handle.readData(ofLength: chunkSize)
            if chunk.isEmpty { break }
            try await connection.writeEncrypted(chunk)
            bytesSent += Int64(chunk.count)
            onProgress?(bytesSent, totalSize)
        }
    }
}
