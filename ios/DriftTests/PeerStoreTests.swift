import Foundation
import Testing
@testable import Drift

@Suite("PeerStore Tests")
struct PeerStoreTests {
    private func makePeer(instance: String, address: String = "192.168.1.1") -> Peer {
        Peer(
            id: instance,
            instance: instance,
            addresses: [address],
            port: 12345,
            publicKey: Data(repeating: 0xAB, count: 32),
            os: "linux",
            protocolVersion: "0.1"
        )
    }

    @Test("Add and retrieve peer")
    func addPeer() {
        let store = PeerStore()
        let peer = makePeer(instance: "alice's laptop")
        store.add(peer)
        #expect(store.peers.count == 1)
        #expect(store.peers[0].instance == "alice's laptop")
    }

    @Test("Update existing peer")
    func updatePeer() {
        let store = PeerStore()
        store.add(makePeer(instance: "test", address: "10.0.0.1"))
        store.add(makePeer(instance: "test", address: "10.0.0.2"))
        #expect(store.peers.count == 1)
        #expect(store.peers[0].addresses == ["10.0.0.2"])
    }

    @Test("Remove peer")
    func removePeer() {
        let store = PeerStore()
        store.add(makePeer(instance: "a"))
        store.add(makePeer(instance: "b"))
        store.remove(instance: "a")
        #expect(store.peers.count == 1)
        #expect(store.peers[0].instance == "b")
    }

    @Test("Get by instance")
    func getByInstance() {
        let store = PeerStore()
        store.add(makePeer(instance: "target"))
        store.add(makePeer(instance: "other"))
        let found = store.getByInstance("target")
        #expect(found?.instance == "target")
    }

    @Test("Get by address")
    func getByAddress() {
        let store = PeerStore()
        store.add(makePeer(instance: "peer1", address: "192.168.1.10"))
        store.add(makePeer(instance: "peer2", address: "192.168.1.20"))
        let found = store.getByAddress("192.168.1.20")
        #expect(found?.instance == "peer2")
    }

    @Test("Get by unknown instance returns nil")
    func unknownInstance() {
        let store = PeerStore()
        #expect(store.getByInstance("nope") == nil)
    }

    @Test("Get by unknown address returns nil")
    func unknownAddress() {
        let store = PeerStore()
        #expect(store.getByAddress("1.2.3.4") == nil)
    }
}
