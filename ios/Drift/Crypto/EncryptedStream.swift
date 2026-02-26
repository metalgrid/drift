import CryptoKit
import Foundation

enum CryptoError: Error {
    case keyDerivationFailed
    case encryptionFailed
    case decryptionFailed
    case invalidEphemeralKey
    case readFailed
    case writeFailed
}

// MARK: - EncryptWriter

/// Encrypts data using AES-GCM and writes to an output function.
/// Matches Go's secret.EncryptWriter.
final class EncryptWriter {
    private let key: SymmetricKey
    private var nonce: NonceCounter
    private let writeFn: (Data) async throws -> Void

    /// Initialize by generating an ephemeral key, deriving shared secret, and writing the ephemeral
    /// public key to the stream.
    init(recipientPublicKey: Data, writeFn: @escaping (Data) async throws -> Void) async throws {
        // Generate ephemeral X25519 key pair
        let ephemeral = Curve25519.KeyAgreement.PrivateKey()

        // Derive shared secret
        guard let peerPub = try? Curve25519.KeyAgreement.PublicKey(rawRepresentation: recipientPublicKey) else {
            throw CryptoError.invalidEphemeralKey
        }
        let sharedSecret = try ephemeral.sharedSecretFromKeyAgreement(with: peerPub)

        // Derive AES key: SHA-256 of raw shared secret bytes
        let sharedSecretData = sharedSecret.withUnsafeBytes { Data($0) }
        let hash = SHA256.hash(data: sharedSecretData)
        self.key = SymmetricKey(data: Data(hash))

        self.nonce = NonceCounter()
        self.writeFn = writeFn

        // Write ephemeral public key (32 bytes) to stream
        try await writeFn(Data(ephemeral.publicKey.rawRepresentation))
    }

    /// Encrypt and write data. Writes [4-byte BE length][ciphertext || tag].
    /// Nonce is incremented after each seal (matching Go).
    func write(_ plaintext: Data) async throws {
        let nonceData = try AES.GCM.Nonce(data: nonce.data)
        let sealed = try AES.GCM.seal(plaintext, using: key, nonce: nonceData)
        nonce.increment()

        // Go's aead.Seal produces ciphertext || tag concatenated
        let encrypted = sealed.ciphertext + sealed.tag

        // Write 4-byte big-endian length
        let length = UInt32(encrypted.count)
        var lengthBuf = Data(count: 4)
        lengthBuf[0] = UInt8((length >> 24) & 0xFF)
        lengthBuf[1] = UInt8((length >> 16) & 0xFF)
        lengthBuf[2] = UInt8((length >> 8) & 0xFF)
        lengthBuf[3] = UInt8(length & 0xFF)

        try await writeFn(lengthBuf)
        try await writeFn(encrypted)
    }
}

// MARK: - DecryptReader

/// Decrypts data using AES-GCM from an input function.
/// Matches Go's secret.DecryptReader.
final class DecryptReader {
    private let key: SymmetricKey
    private var nonce: NonceCounter
    private let readFn: (Int) async throws -> Data

    /// Initialize by reading the ephemeral public key from the stream and deriving the shared secret.
    init(localPrivateKey: Curve25519.KeyAgreement.PrivateKey, readFn: @escaping (Int) async throws -> Data) async throws {
        // Read 32-byte ephemeral public key
        let ephemeralPubData = try await readFn(32)
        guard ephemeralPubData.count == 32 else {
            throw CryptoError.invalidEphemeralKey
        }

        let ephemeralPub = try Curve25519.KeyAgreement.PublicKey(rawRepresentation: ephemeralPubData)

        // Derive shared secret
        let sharedSecret = try localPrivateKey.sharedSecretFromKeyAgreement(with: ephemeralPub)

        // Derive AES key: SHA-256 of raw shared secret bytes
        let sharedSecretData = sharedSecret.withUnsafeBytes { Data($0) }
        let hash = SHA256.hash(data: sharedSecretData)
        self.key = SymmetricKey(data: Data(hash))

        self.nonce = NonceCounter()
        self.readFn = readFn
    }

    /// Read and decrypt one frame. Reads [4-byte BE length][ciphertext || tag].
    /// Nonce is incremented after each open (matching Go).
    func read() async throws -> Data {
        // Read 4-byte big-endian length
        let lengthBuf = try await readFn(4)
        guard lengthBuf.count == 4 else {
            throw CryptoError.readFailed
        }

        let length = (UInt32(lengthBuf[0]) << 24)
            | (UInt32(lengthBuf[1]) << 16)
            | (UInt32(lengthBuf[2]) << 8)
            | UInt32(lengthBuf[3])

        // Read encrypted data
        let encrypted = try await readFn(Int(length))
        guard encrypted.count == Int(length) else {
            throw CryptoError.readFailed
        }

        // Split: last 16 bytes are GCM tag, rest is ciphertext
        let tagSize = 16
        guard encrypted.count >= tagSize else {
            throw CryptoError.decryptionFailed
        }

        let ciphertext = encrypted.prefix(encrypted.count - tagSize)
        let tag = encrypted.suffix(tagSize)

        let nonceData = try AES.GCM.Nonce(data: nonce.data)
        let sealedBox = try AES.GCM.SealedBox(nonce: nonceData, ciphertext: ciphertext, tag: tag)
        let decrypted = try AES.GCM.open(sealedBox, using: key)

        nonce.increment()

        return decrypted
    }
}
