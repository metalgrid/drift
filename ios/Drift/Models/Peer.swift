import Foundation
import Network

struct Peer: Identifiable, Equatable {
    let id: String // instance name as unique ID
    let instance: String
    var addresses: [String]
    var port: UInt16
    var publicKey: Data // 32 bytes
    var os: String
    var protocolVersion: String

    var publicKeyHex: String {
        publicKey.hexString
    }

    static func == (lhs: Peer, rhs: Peer) -> Bool {
        lhs.instance == rhs.instance
    }

    /// Parse a Peer from Bonjour TXT records and endpoint info.
    static func fromTXTRecords(_ records: [String: String], instance: String, addresses: [String], port: UInt16) -> Peer? {
        guard let pkHex = records["pk"],
              let pkData = Data(hexString: pkHex),
              pkData.count == 32 else {
            return nil
        }

        return Peer(
            id: instance,
            instance: instance,
            addresses: addresses,
            port: port,
            publicKey: pkData,
            os: records["os"] ?? "unknown",
            protocolVersion: records["v"] ?? "0.1"
        )
    }
}
