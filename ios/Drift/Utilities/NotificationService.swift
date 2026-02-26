import Foundation
import UserNotifications

/// Handles local notifications for incoming transfer offers.
final class NotificationService {
    static let shared = NotificationService()

    private init() {}

    func requestAuthorization() {
        UNUserNotificationCenter.current().requestAuthorization(options: [.alert, .sound, .badge]) { granted, error in
            if let error {
                print("Notification auth error: \(error)")
            }
        }
    }

    func postIncomingOffer(from peerName: String, filename: String, size: String) {
        let content = UNMutableNotificationContent()
        content.title = "Incoming File"
        content.body = "\(peerName) wants to send you \(filename) (\(size))"
        content.sound = .default
        content.categoryIdentifier = "TRANSFER_OFFER"

        let request = UNNotificationRequest(
            identifier: UUID().uuidString,
            content: content,
            trigger: nil
        )

        UNUserNotificationCenter.current().add(request)
    }

    func registerCategories() {
        let accept = UNNotificationAction(identifier: "ACCEPT", title: "Accept", options: .foreground)
        let decline = UNNotificationAction(identifier: "DECLINE", title: "Decline", options: .destructive)
        let category = UNNotificationCategory(
            identifier: "TRANSFER_OFFER",
            actions: [accept, decline],
            intentIdentifiers: []
        )
        UNUserNotificationCenter.current().setNotificationCategories([category])
    }
}
