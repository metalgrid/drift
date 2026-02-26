import Testing
@testable import Drift

@Suite("Message Tests")
struct MessageTests {
    @Test("Serialize single offer")
    func serializeOffer() {
        let msg = DriftMessage.makeOffer(filename: "test.txt", size: 1024)
        let data = msg.serialize()
        let str = String(data: data, encoding: .utf8)!
        #expect(str == "OFFER|test.txt|application/octet-stream|1024\n")
    }

    @Test("Parse single offer")
    func parseOffer() {
        let msg = DriftMessage.parse("OFFER|document.pdf|application/octet-stream|2048\n")
        guard case let .offer(filename, mimetype, size) = msg else {
            Issue.record("Expected offer")
            return
        }
        #expect(filename == "document.pdf")
        #expect(mimetype == "application/octet-stream")
        #expect(size == 2048)
    }

    @Test("Serialize batch offer")
    func serializeBatchOffer() {
        let msg = DriftMessage.makeBatchOffer(files: [
            (filename: "a.txt", size: 100),
            (filename: "b.pdf", size: 200),
        ])
        let data = msg.serialize()
        let str = String(data: data, encoding: .utf8)!
        #expect(str == "BATCH_OFFER|2|a.txt|application/octet-stream|100|b.pdf|application/octet-stream|200\n")
    }

    @Test("Parse batch offer")
    func parseBatchOffer() {
        let raw = "BATCH_OFFER|2|file1.txt|application/octet-stream|1024|file2.pdf|application/octet-stream|2048\n"
        let msg = DriftMessage.parse(raw)
        guard case let .batchOffer(files) = msg else {
            Issue.record("Expected batch offer")
            return
        }
        #expect(files.count == 2)
        #expect(files[0].filename == "file1.txt")
        #expect(files[0].size == 1024)
        #expect(files[1].filename == "file2.pdf")
        #expect(files[1].size == 2048)
    }

    @Test("Serialize accept answer")
    func serializeAccept() {
        let msg = DriftMessage.accept()
        let str = String(data: msg.serialize(), encoding: .utf8)!
        #expect(str == "ANSWER|ACCEPT\n")
    }

    @Test("Serialize decline answer")
    func serializeDecline() {
        let msg = DriftMessage.decline()
        let str = String(data: msg.serialize(), encoding: .utf8)!
        #expect(str == "ANSWER|DECLINE\n")
    }

    @Test("Parse accept answer")
    func parseAccept() {
        let msg = DriftMessage.parse("ANSWER|ACCEPT\n")
        #expect(msg?.isAccepted == true)
    }

    @Test("Parse decline answer")
    func parseDecline() {
        let msg = DriftMessage.parse("ANSWER|DECLINE\n")
        #expect(msg?.isAccepted == false)
    }

    @Test("Round-trip single offer")
    func roundTripOffer() {
        let original = DriftMessage.makeOffer(filename: "photo.jpg", size: 999999)
        let serialized = String(data: original.serialize(), encoding: .utf8)!
        let parsed = DriftMessage.parse(serialized)
        guard case let .offer(filename, _, size) = parsed else {
            Issue.record("Round-trip failed")
            return
        }
        #expect(filename == "photo.jpg")
        #expect(size == 999999)
    }

    @Test("Round-trip batch offer")
    func roundTripBatch() {
        let original = DriftMessage.makeBatchOffer(files: [
            (filename: "a.bin", size: 1),
            (filename: "b.bin", size: Int64.max),
        ])
        let serialized = String(data: original.serialize(), encoding: .utf8)!
        let parsed = DriftMessage.parse(serialized)
        guard case let .batchOffer(files) = parsed else {
            Issue.record("Round-trip failed")
            return
        }
        #expect(files.count == 2)
        #expect(files[1].size == Int64.max)
    }

    @Test("Parse invalid message returns nil")
    func parseInvalid() {
        #expect(DriftMessage.parse("GARBAGE|data\n") == nil)
        #expect(DriftMessage.parse("OFFER|too|few\n") == nil)
        #expect(DriftMessage.parse("BATCH_OFFER|1\n") == nil)
        #expect(DriftMessage.parse("ANSWER\n") == nil)
    }

    @Test("BATCH_OFFER checked before OFFER prefix")
    func batchBeforeOfferPrefix() {
        // "BATCH_OFFER" starts with "OFFER" would be ambiguous without prefix check
        let raw = "BATCH_OFFER|1|x.txt|application/octet-stream|42\n"
        let msg = DriftMessage.parse(raw)
        guard case .batchOffer = msg else {
            Issue.record("Should parse as batch, not single offer")
            return
        }
    }
}
