import CryptoKit
import Foundation

struct DriftKeyPair {
    let privateKey: Curve25519.KeyAgreement.PrivateKey
    let publicKeyData: Data
    let publicKeyHex: String

    init() {
        let priv = Curve25519.KeyAgreement.PrivateKey()
        self.privateKey = priv
        self.publicKeyData = priv.publicKey.rawRepresentation
        self.publicKeyHex = priv.publicKey.rawRepresentation.hexString
    }
}
