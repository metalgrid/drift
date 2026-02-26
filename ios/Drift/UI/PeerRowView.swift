import SwiftUI
import UniformTypeIdentifiers

struct PeerRowView: View {
    let peer: Peer
    let coordinator: AppCoordinator
    @State private var showingFilePicker = false

    var body: some View {
        HStack {
            Image(systemName: osIcon)
                .font(.title2)
                .foregroundStyle(.secondary)
                .frame(width: 32)

            VStack(alignment: .leading, spacing: 2) {
                Text(peer.instance)
                    .font(.body)
                Text(peer.os)
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }

            Spacer()

            Button {
                showingFilePicker = true
            } label: {
                Image(systemName: "paperplane.fill")
                    .font(.title3)
            }
        }
        .fileImporter(
            isPresented: $showingFilePicker,
            allowedContentTypes: [.item],
            allowsMultipleSelection: true
        ) { result in
            switch result {
            case .success(let urls):
                coordinator.sendFiles(to: peer, fileURLs: urls)
            case .failure(let error):
                print("File picker error: \(error)")
            }
        }
    }

    private var osIcon: String {
        switch peer.os.lowercased() {
        case "darwin", "macos":
            return "desktopcomputer"
        case "linux":
            return "terminal"
        case "windows":
            return "pc"
        case "ios":
            return "iphone"
        default:
            return "questionmark.circle"
        }
    }
}
