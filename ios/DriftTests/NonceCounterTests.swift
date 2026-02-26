import Testing
@testable import Drift

@Suite("NonceCounter Tests")
struct NonceCounterTests {
    @Test("Starts at all zeros")
    func initialState() {
        let nonce = NonceCounter()
        #expect(nonce.bytes == [UInt8](repeating: 0, count: 12))
    }

    @Test("Single increment")
    func singleIncrement() {
        var nonce = NonceCounter()
        nonce.increment()
        var expected = [UInt8](repeating: 0, count: 12)
        expected[11] = 1
        #expect(nonce.bytes == expected)
    }

    @Test("Multiple increments")
    func multipleIncrements() {
        var nonce = NonceCounter()
        for _ in 0..<255 {
            nonce.increment()
        }
        var expected = [UInt8](repeating: 0, count: 12)
        expected[11] = 255
        #expect(nonce.bytes == expected)
    }

    @Test("Carry on overflow")
    func carryOverflow() {
        var nonce = NonceCounter()
        // Increment 256 times to trigger carry
        for _ in 0..<256 {
            nonce.increment()
        }
        var expected = [UInt8](repeating: 0, count: 12)
        expected[10] = 1
        expected[11] = 0
        #expect(nonce.bytes == expected)
    }

    @Test("Double carry")
    func doubleCarry() {
        var nonce = NonceCounter()
        // Increment 65536 times: should carry twice
        for _ in 0..<65536 {
            nonce.increment()
        }
        var expected = [UInt8](repeating: 0, count: 12)
        expected[9] = 1
        expected[10] = 0
        expected[11] = 0
        #expect(nonce.bytes == expected)
    }

    @Test("Data representation is 12 bytes")
    func dataRepresentation() {
        let nonce = NonceCounter()
        #expect(nonce.data.count == 12)
    }
}
