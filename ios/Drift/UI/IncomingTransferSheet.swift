import SwiftUI

struct IncomingTransferSheet: View {
    let offer: TransferOffer
    let onRespond: (Bool) -> Void

    @State private var timeRemaining = 30
    @State private var timer: Timer?

    var body: some View {
        VStack(spacing: 20) {
            Image(systemName: "arrow.down.doc.fill")
                .font(.system(size: 48))
                .foregroundStyle(.blue)

            Text("Incoming Transfer")
                .font(.title2.bold())

            Text("From: \(offer.peerName)")
                .font(.subheadline)
                .foregroundStyle(.secondary)

            if offer.isBatch {
                VStack(alignment: .leading, spacing: 4) {
                    Text("\(offer.files.count) files")
                        .font(.headline)
                    ForEach(offer.files) { file in
                        HStack {
                            Text(file.filename)
                                .font(.caption)
                                .lineLimit(1)
                            Spacer()
                            Text(file.formattedSize)
                                .font(.caption)
                                .foregroundStyle(.secondary)
                        }
                    }
                }
                .padding()
                .background(.regularMaterial, in: RoundedRectangle(cornerRadius: 8))
            } else if let file = offer.files.first {
                VStack(spacing: 4) {
                    Text(file.filename)
                        .font(.headline)
                    Text(file.formattedSize)
                        .foregroundStyle(.secondary)
                }
            }

            Text("Total: \(offer.formattedTotalSize)")
                .font(.subheadline)

            Text("Auto-declining in \(timeRemaining)s")
                .font(.caption)
                .foregroundStyle(.secondary)

            HStack(spacing: 16) {
                Button(role: .destructive) {
                    stopTimer()
                    onRespond(false)
                } label: {
                    Text("Decline")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.bordered)

                Button {
                    stopTimer()
                    onRespond(true)
                } label: {
                    Text("Accept")
                        .frame(maxWidth: .infinity)
                }
                .buttonStyle(.borderedProminent)
            }
            .padding(.top)
        }
        .padding(24)
        .presentationDetents([.medium])
        .onAppear { startTimer() }
        .onDisappear { stopTimer() }
    }

    private func startTimer() {
        timeRemaining = 30
        timer = Timer.scheduledTimer(withTimeInterval: 1, repeats: true) { _ in
            timeRemaining -= 1
            if timeRemaining <= 0 {
                stopTimer()
                onRespond(false)
            }
        }
    }

    private func stopTimer() {
        timer?.invalidate()
        timer = nil
    }
}
