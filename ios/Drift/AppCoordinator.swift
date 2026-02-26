import CryptoKit
import Foundation
import Network
import Observation

/// Wires together discovery, server, crypto, and UI.
/// Mirrors Go's app.Run().
@Observable
final class AppCoordinator {
    let keyPair: DriftKeyPair
    let peerStore = PeerStore()

    private let tcpListener = TCPListener()
    private var browser: BonjourBrowser?

    var incomingOffer: TransferOffer?
    var incomingOfferContinuation: CheckedContinuation<Bool, Never>?

    var transferProgress: Double = 0
    var transferActive = false
    var statusMessage: String?

    var identityName: String {
        get { UserDefaults.standard.string(forKey: "drift_identity") ?? UIDevice.current.name }
        set { UserDefaults.standard.set(newValue, forKey: "drift_identity") }
    }

    init() {
        keyPair = DriftKeyPair()
    }

    /// Start all services: TCP listener with Bonjour, browser for peer discovery.
    func start() {
        do {
            tcpListener.onConnection = { [weak self] connection in
                self?.handleInboundConnection(connection)
            }

            // When listener is ready and port is assigned, advertise + browse
            tcpListener.onReady = { [weak self] in
                self?.startNetworkServices()
            }

            try tcpListener.start()
        } catch {
            statusMessage = "Failed to start: \(error.localizedDescription)"
        }
    }

    private func startNetworkServices() {
        // Advertise Bonjour service on the same listener (no port conflict)
        tcpListener.advertise(instanceName: identityName, publicKeyHex: keyPair.publicKeyHex)

        // Start Bonjour browser to discover peers
        let brow = BonjourBrowser(peerStore: peerStore)
        brow.start()
        browser = brow
    }

    func stop() {
        tcpListener.stop()
        browser?.stop()
    }

    // MARK: - Inbound connections

    private func handleInboundConnection(_ nwConnection: NWConnection) {
        nwConnection.start(queue: .main)

        nwConnection.stateUpdateHandler = { [weak self] state in
            guard let self else { return }
            if case .ready = state {
                Task { await self.processInbound(nwConnection) }
            }
        }
    }

    private func processInbound(_ nwConnection: NWConnection) async {
        guard let path = nwConnection.currentPath,
              let remoteEndpoint = path.remoteEndpoint,
              case let .hostPort(host, _) = remoteEndpoint else {
            nwConnection.cancel()
            return
        }

        let remoteAddress: String
        switch host {
        case .ipv4(let addr): remoteAddress = "\(addr)"
        case .ipv6(let addr): remoteAddress = "\(addr)"
        case .name(let name, _): remoteAddress = name
        @unknown default: nwConnection.cancel(); return
        }

        guard let peer = peerStore.getByAddress(remoteAddress) else {
            print("Unknown peer from \(remoteAddress)")
            nwConnection.cancel()
            return
        }

        do {
            let secureConn = SecureConnection(connection: nwConnection)
            try await secureConn.establish(
                peerPublicKey: peer.publicKey,
                localPrivateKey: keyPair.privateKey
            )

            let handler = ConnectionHandler(
                connection: secureConn,
                onOffer: { [weak self] offer in
                    await self?.promptUserForOffer(offer) ?? false
                },
                onProgress: { [weak self] current, total in
                    DispatchQueue.main.async {
                        self?.transferProgress = Double(current) / Double(total)
                        self?.transferActive = total > 0
                    }
                },
                onComplete: { [weak self] message in
                    DispatchQueue.main.async {
                        self?.transferActive = false
                        self?.transferProgress = 0
                        self?.statusMessage = message
                    }
                },
                onError: { [weak self] message in
                    DispatchQueue.main.async {
                        self?.transferActive = false
                        self?.statusMessage = message
                    }
                }
            )
            await handler.run(peerName: peer.instance)
        } catch {
            print("Failed to secure inbound connection: \(error)")
            nwConnection.cancel()
        }
    }

    // MARK: - Outbound transfers

    func sendFiles(to peer: Peer, fileURLs: [URL]) {
        Task {
            do {
                transferActive = true
                transferProgress = 0

                let nwConnection = NWConnection(
                    host: NWEndpoint.Host(peer.addresses.first ?? ""),
                    port: NWEndpoint.Port(rawValue: peer.port)!,
                    using: .tcp
                )

                try await withCheckedThrowingContinuation { (continuation: CheckedContinuation<Void, Error>) in
                    nwConnection.stateUpdateHandler = { state in
                        switch state {
                        case .ready:
                            continuation.resume()
                        case .failed(let error):
                            continuation.resume(throwing: error)
                        default:
                            break
                        }
                    }
                    nwConnection.start(queue: .main)
                }

                let secureConn = SecureConnection(connection: nwConnection)
                try await secureConn.establish(
                    peerPublicKey: peer.publicKey,
                    localPrivateKey: keyPair.privateKey
                )

                let handler = ConnectionHandler(
                    connection: secureConn,
                    onOffer: { _ in false },
                    onProgress: { [weak self] current, total in
                        DispatchQueue.main.async {
                            self?.transferProgress = Double(current) / Double(total)
                        }
                    },
                    onComplete: { [weak self] message in
                        DispatchQueue.main.async {
                            self?.transferActive = false
                            self?.transferProgress = 0
                            self?.statusMessage = message
                        }
                    },
                    onError: { [weak self] message in
                        DispatchQueue.main.async {
                            self?.transferActive = false
                            self?.statusMessage = message
                        }
                    }
                )

                try await handler.sendFiles(fileURLs)
            } catch {
                transferActive = false
                statusMessage = "Send failed: \(error.localizedDescription)"
            }
        }
    }

    // MARK: - User prompt

    @MainActor
    private func promptUserForOffer(_ offer: TransferOffer) async -> Bool {
        incomingOffer = offer
        return await withCheckedContinuation { continuation in
            incomingOfferContinuation = continuation
        }
    }

    @MainActor
    func respondToOffer(accepted: Bool) {
        incomingOfferContinuation?.resume(returning: accepted)
        incomingOfferContinuation = nil
        incomingOffer = nil
    }
}
