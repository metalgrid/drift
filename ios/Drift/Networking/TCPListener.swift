import Foundation
import Network

/// TCP listener on a random port with integrated Bonjour advertisement.
/// Matches Go's server.Start() + zeroconf Publish.
final class TCPListener {
    private var listener: NWListener?
    private(set) var port: UInt16 = 0
    var onConnection: ((NWConnection) -> Void)?
    var onReady: (() -> Void)?

    func start() throws {
        let listener = try NWListener(using: .tcp, on: .any)

        listener.stateUpdateHandler = { [weak self] state in
            switch state {
            case .ready:
                self?.port = listener.port?.rawValue ?? 0
                self?.onReady?()
            case .failed(let error):
                print("Listener failed: \(error)")
                listener.cancel()
            default:
                break
            }
        }

        listener.newConnectionHandler = { [weak self] connection in
            self?.onConnection?(connection)
        }

        listener.start(queue: .main)
        self.listener = listener
    }

    /// Set Bonjour service on the existing listener to advertise this device.
    /// Must be called after the listener is ready (port is assigned).
    /// TXT records match Go: v=0.1, pk=<hex>, os=ios, port=<port>
    func advertise(instanceName: String, publicKeyHex: String) {
        guard let listener else { return }

        let txtRecord = NWTXTRecord([
            "v": "0.1",
            "pk": publicKeyHex,
            "os": "ios",
            "port": String(port),
        ])

        listener.service = NWListener.Service(
            name: instanceName,
            type: "_drift._tcp",
            txtRecord: txtRecord
        )
    }

    func stop() {
        listener?.cancel()
        listener = nil
    }
}
