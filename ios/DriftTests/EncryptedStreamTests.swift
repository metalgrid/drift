import CryptoKit
import Foundation
import Testing
@testable import Drift

@Suite("EncryptedStream Tests")
struct EncryptedStreamTests {
    @Test("Encrypt then decrypt round-trip")
    func roundTrip() async throws {
        var pipe = Data()
        var readOffset = 0

        let recipientKey = Curve25519.KeyAgreement.PrivateKey()
        let recipientPub = recipientKey.publicKey.rawRepresentation

        let writer = try await EncryptWriter(
            recipientPublicKey: Data(recipientPub),
            writeFn: { data in pipe.append(data) }
        )

        let plaintext1 = Data("Hello, Drift!".utf8)
        let plaintext2 = Data("Second message".utf8)
        try await writer.write(plaintext1)
        try await writer.write(plaintext2)

        let reader = try await DecryptReader(
            localPrivateKey: recipientKey,
            readFn: { count in
                guard readOffset + count <= pipe.count else {
                    throw CryptoError.readFailed
                }
                let data = pipe[readOffset..<(readOffset + count)]
                readOffset += count
                return Data(data)
            }
        )

        let decrypted1 = try await reader.read()
        let decrypted2 = try await reader.read()

        #expect(decrypted1 == plaintext1)
        #expect(decrypted2 == plaintext2)
    }

    @Test("Multiple messages maintain nonce sync")
    func nonceSync() async throws {
        var pipe = Data()
        var readOffset = 0

        let recipientKey = Curve25519.KeyAgreement.PrivateKey()

        let writer = try await EncryptWriter(
            recipientPublicKey: Data(recipientKey.publicKey.rawRepresentation),
            writeFn: { data in pipe.append(data) }
        )

        let messageCount = 100
        var originals: [Data] = []
        for i in 0..<messageCount {
            let msg = Data("Message \(i)".utf8)
            originals.append(msg)
            try await writer.write(msg)
        }

        let reader = try await DecryptReader(
            localPrivateKey: recipientKey,
            readFn: { count in
                guard readOffset + count <= pipe.count else {
                    throw CryptoError.readFailed
                }
                let data = pipe[readOffset..<(readOffset + count)]
                readOffset += count
                return Data(data)
            }
        )

        for i in 0..<messageCount {
            let decrypted = try await reader.read()
            #expect(decrypted == originals[i])
        }
    }

    @Test("Empty data round-trip")
    func emptyData() async throws {
        var pipe = Data()
        var readOffset = 0

        let recipientKey = Curve25519.KeyAgreement.PrivateKey()

        let writer = try await EncryptWriter(
            recipientPublicKey: Data(recipientKey.publicKey.rawRepresentation),
            writeFn: { data in pipe.append(data) }
        )

        let empty = Data()
        try await writer.write(empty)

        let reader = try await DecryptReader(
            localPrivateKey: recipientKey,
            readFn: { count in
                guard readOffset + count <= pipe.count else {
                    throw CryptoError.readFailed
                }
                let data = pipe[readOffset..<(readOffset + count)]
                readOffset += count
                return Data(data)
            }
        )

        let decrypted = try await reader.read()
        #expect(decrypted == empty)
    }

    @Test("Large data round-trip")
    func largeData() async throws {
        var pipe = Data()
        var readOffset = 0

        let recipientKey = Curve25519.KeyAgreement.PrivateKey()

        let writer = try await EncryptWriter(
            recipientPublicKey: Data(recipientKey.publicKey.rawRepresentation),
            writeFn: { data in pipe.append(data) }
        )

        // 64KB of data
        let large = Data((0..<65536).map { UInt8($0 & 0xFF) })
        try await writer.write(large)

        let reader = try await DecryptReader(
            localPrivateKey: recipientKey,
            readFn: { count in
                guard readOffset + count <= pipe.count else {
                    throw CryptoError.readFailed
                }
                let data = pipe[readOffset..<(readOffset + count)]
                readOffset += count
                return Data(data)
            }
        )

        let decrypted = try await reader.read()
        #expect(decrypted == large)
    }
}
