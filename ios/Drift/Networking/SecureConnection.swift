import CryptoKit
import Foundation
import Network

/// Wraps an NWConnection with encrypted read/write via EncryptWriter + DecryptReader.
/// Matches Go's secret.SecureConnection().
final class SecureConnection {
    let connection: NWConnection
    private var encryptWriter: EncryptWriter?
    private var decryptReader: DecryptReader?

    init(connection: NWConnection) {
        self.connection = connection
    }

    /// Establish encryption on this connection.
    /// Creates EncryptWriter (using peer's pubkey) and DecryptReader (using own privkey).
    func establish(peerPublicKey: Data, localPrivateKey: Curve25519.KeyAgreement.PrivateKey) async throws {
        // Set up EncryptWriter - writes ephemeral pubkey then encrypts outgoing data
        self.encryptWriter = try await EncryptWriter(
            recipientPublicKey: peerPublicKey,
            writeFn: { [weak self] data in
                guard let self else { throw CryptoError.writeFailed }
                try await self.writeRaw(data)
            }
        )

        // Set up DecryptReader - reads ephemeral pubkey then decrypts incoming data
        self.decryptReader = try await DecryptReader(
            localPrivateKey: localPrivateKey,
            readFn: { [weak self] count in
                guard let self else { throw CryptoError.readFailed }
                return try await self.readExactly(count: count)
            }
        )
    }

    /// Write plaintext through the encrypted channel.
    func writeEncrypted(_ data: Data) async throws {
        guard let writer = encryptWriter else {
            throw CryptoError.writeFailed
        }
        try await writer.write(data)
    }

    /// Read one decrypted frame from the encrypted channel.
    func readDecrypted() async throws -> Data {
        guard let reader = decryptReader else {
            throw CryptoError.readFailed
        }
        return try await reader.read()
    }

    /// Read exactly `count` bytes from the NWConnection.
    private func readExactly(count: Int) async throws -> Data {
        var collected = Data()
        collected.reserveCapacity(count)

        while collected.count < count {
            let remaining = count - collected.count
            let chunk = try await withCheckedThrowingContinuation { continuation in
                connection.receive(minimumIncompleteLength: 1, maximumLength: remaining) { data, _, isComplete, error in
                    if let error {
                        continuation.resume(throwing: error)
                        return
                    }
                    guard let data, !data.isEmpty else {
                        if isComplete {
                            continuation.resume(throwing: CryptoError.readFailed)
                        } else {
                            continuation.resume(throwing: CryptoError.readFailed)
                        }
                        return
                    }
                    continuation.resume(returning: data)
                }
            }

            collected.append(chunk)
        }

        return collected
    }

    /// Write raw bytes to the NWConnection.
    private func writeRaw(_ data: Data) async throws {
        try await withCheckedThrowingContinuation { (continuation: CheckedContinuation<Void, Error>) in
            connection.send(content: data, completion: .contentProcessed { error in
                if let error {
                    continuation.resume(throwing: error)
                } else {
                    continuation.resume()
                }
            })
        }
    }

    func close() {
        connection.cancel()
    }
}
