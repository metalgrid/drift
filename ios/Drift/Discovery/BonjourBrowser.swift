import Foundation
import Network

/// Browses for _drift._tcp services on the local network.
/// Matches Go's zeroconf browse functionality.
final class BonjourBrowser {
    private var browser: NWBrowser?
    private let peerStore: PeerStore

    init(peerStore: PeerStore) {
        self.peerStore = peerStore
    }

    func start() {
        let params = NWParameters()
        params.includePeerToPeer = true

        let browser = NWBrowser(for: .bonjour(type: "_drift._tcp", domain: "local."), using: params)

        browser.stateUpdateHandler = { state in
            switch state {
            case .ready:
                print("Browser ready")
            case .failed(let error):
                print("Browser failed: \(error)")
            default:
                break
            }
        }

        browser.browseResultsChangedHandler = { [weak self] results, changes in
            self?.handleChanges(results: results, changes: changes)
        }

        browser.start(queue: .main)
        self.browser = browser
    }

    func stop() {
        browser?.cancel()
        browser = nil
    }

    private func handleChanges(results: Set<NWBrowser.Result>, changes: Set<NWBrowser.Result.Change>) {
        for change in changes {
            switch change {
            case .added(let result):
                resolveAndAdd(result)
            case .removed(let result):
                if case let .service(name, _, _, _) = result.endpoint {
                    peerStore.remove(instance: name)
                }
            case .changed(old: _, new: let result, flags: _):
                resolveAndAdd(result)
            @unknown default:
                break
            }
        }
    }

    private func resolveAndAdd(_ result: NWBrowser.Result) {
        guard case let .service(name, _, _, _) = result.endpoint else { return }

        // Extract TXT records from metadata
        var records: [String: String] = [:]
        if case let .bonjour(txtRecord) = result.metadata {
            // NWTXTRecord doesn't have a direct iteration API,
            // so we parse the raw dictionary representation
            let dict = txtRecord.dictionary
            records = dict
        }

        // Create connection to resolve addresses
        let params = NWParameters.tcp
        let connection = NWConnection(to: result.endpoint, using: params)

        connection.stateUpdateHandler = { [weak self] state in
            if case .ready = state {
                if let endpoint = connection.currentPath?.remoteEndpoint,
                   case let .hostPort(host, port) = endpoint {
                    let address: String
                    switch host {
                    case .ipv4(let addr):
                        address = "\(addr)"
                    case .ipv6(let addr):
                        address = "\(addr)"
                    case .name(let hostname, _):
                        address = hostname
                    @unknown default:
                        address = "unknown"
                    }

                    let resolvedPort = records["port"].flatMap(UInt16.init) ?? port.rawValue

                    if let peer = Peer.fromTXTRecords(records, instance: name, addresses: [address], port: resolvedPort) {
                        self?.peerStore.add(peer)
                    }
                }
                connection.cancel()
            } else if case .failed = state {
                connection.cancel()
            }
        }
        connection.start(queue: .main)
    }
}
