import SwiftUI

@main
struct DriftApp: App {
    @State private var coordinator = AppCoordinator()

    var body: some Scene {
        WindowGroup {
            PeerListView(coordinator: coordinator)
                .onAppear {
                    NotificationService.shared.requestAuthorization()
                    NotificationService.shared.registerCategories()
                    coordinator.start()
                }
        }
    }
}
