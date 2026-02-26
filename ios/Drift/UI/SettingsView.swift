import SwiftUI

struct SettingsView: View {
    let coordinator: AppCoordinator
    @State private var name: String = ""

    var body: some View {
        Form {
            Section("Identity") {
                TextField("Device Name", text: $name)
                    .onAppear {
                        name = coordinator.identityName
                    }
                    .onChange(of: name) { _, newValue in
                        coordinator.identityName = newValue
                    }
            }

            Section("Info") {
                LabeledContent("Public Key") {
                    Text(coordinator.keyPair.publicKeyHex.prefix(16) + "...")
                        .font(.caption.monospaced())
                        .foregroundStyle(.secondary)
                }
                LabeledContent("Protocol Version", value: "0.1")
                LabeledContent("Platform", value: "iOS")
            }
        }
        .navigationTitle("Settings")
        .navigationBarTitleDisplayMode(.inline)
    }
}
