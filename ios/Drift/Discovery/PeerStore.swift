import Foundation
import Observation

/// Observable store of discovered peers.
/// Matches Go's zeroconf.Peers.
@Observable
final class PeerStore {
    private(set) var peers: [Peer] = []

    func add(_ peer: Peer) {
        if let index = peers.firstIndex(where: { $0.instance == peer.instance }) {
            peers[index] = peer
        } else {
            peers.append(peer)
        }
    }

    func remove(instance: String) {
        peers.removeAll { $0.instance == instance }
    }

    func getByInstance(_ instance: String) -> Peer? {
        peers.first { $0.instance == instance }
    }

    func getByAddress(_ address: String) -> Peer? {
        peers.first { $0.addresses.contains(address) }
    }
}
