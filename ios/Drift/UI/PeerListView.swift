import SwiftUI

struct PeerListView: View {
    let coordinator: AppCoordinator

    var body: some View {
        NavigationStack {
            Group {
                if coordinator.peerStore.peers.isEmpty {
                    ContentUnavailableView(
                        "No Peers Found",
                        systemImage: "antenna.radiowaves.left.and.right",
                        description: Text("Looking for nearby Drift devices...")
                    )
                } else {
                    List(coordinator.peerStore.peers) { peer in
                        PeerRowView(peer: peer, coordinator: coordinator)
                    }
                    .refreshable {
                        // Pull-to-refresh: peers update automatically via Bonjour
                        try? await Task.sleep(for: .seconds(1))
                    }
                }
            }
            .overlay {
                if coordinator.transferActive {
                    TransferProgressView(progress: coordinator.transferProgress)
                }
            }
            .navigationTitle("Drift")
            .toolbar {
                NavigationLink {
                    SettingsView(coordinator: coordinator)
                } label: {
                    Image(systemName: "gear")
                }
            }
            .sheet(item: Binding(
                get: { coordinator.incomingOffer },
                set: { _ in }
            )) { offer in
                IncomingTransferSheet(offer: offer) { accepted in
                    coordinator.respondToOffer(accepted: accepted)
                }
            }
        }
    }
}
