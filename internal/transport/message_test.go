package transport

import (
	"bytes"
	"os"
	"testing"
)

// TestOfferMarshalMessage tests that Offer.MarshalMessage() produces correct wire format
func TestOfferMarshalMessage(t *testing.T) {
	offer := Offer{
		Message{"OFFER"},
		"test.txt",
		"application/octet-stream",
		1024,
	}
	expected := []byte("OFFER|test.txt|application/octet-stream|1024\n")
	result := offer.MarshalMessage()
	if !bytes.Equal(result, expected) {
		t.Errorf("MarshalMessage() = %q, want %q", result, expected)
	}
}

// TestOfferMarshalRoundTrip tests marshal → unmarshal → type assert → compare
func TestOfferMarshalRoundTrip(t *testing.T) {
	original := Offer{
		Message{"OFFER"},
		"document.pdf",
		"application/pdf",
		2048,
	}
	marshaled := original.MarshalMessage()
	unmarshaled := UnmarshalMessage(string(marshaled))

	offer, ok := unmarshaled.(Offer)
	if !ok {
		t.Fatalf("UnmarshalMessage() returned %T, want Offer", unmarshaled)
	}

	if offer.Type != original.Type {
		t.Errorf("Type = %q, want %q", offer.Type, original.Type)
	}
	if offer.Filename != original.Filename {
		t.Errorf("Filename = %q, want %q", offer.Filename, original.Filename)
	}
	if offer.Mimetype != original.Mimetype {
		t.Errorf("Mimetype = %q, want %q", offer.Mimetype, original.Mimetype)
	}
	if offer.Size != original.Size {
		t.Errorf("Size = %d, want %d", offer.Size, original.Size)
	}
}

// TestAnswerMarshalAccept tests Accept().MarshalMessage() produces correct format
func TestAnswerMarshalAccept(t *testing.T) {
	answer := Accept()
	expected := []byte("ANSWER|ACCEPT\n")
	result := answer.MarshalMessage()
	if !bytes.Equal(result, expected) {
		t.Errorf("Accept().MarshalMessage() = %q, want %q", result, expected)
	}
}

// TestAnswerMarshalDecline tests Decline().MarshalMessage() produces correct format
func TestAnswerMarshalDecline(t *testing.T) {
	answer := Decline()
	expected := []byte("ANSWER|DECLINE\n")
	result := answer.MarshalMessage()
	if !bytes.Equal(result, expected) {
		t.Errorf("Decline().MarshalMessage() = %q, want %q", result, expected)
	}
}

// TestAnswerAccepted tests Accepted() returns true for ACCEPT, false for DECLINE
func TestAnswerAccepted(t *testing.T) {
	acceptAnswer := Accept()
	if !acceptAnswer.Accepted() {
		t.Errorf("Accept().Accepted() = false, want true")
	}

	declineAnswer := Decline()
	if declineAnswer.Accepted() {
		t.Errorf("Decline().Accepted() = true, want false")
	}
}

// TestUnmarshalInvalidMessage tests that invalid message type returns nil
func TestUnmarshalInvalidMessage(t *testing.T) {
	result := UnmarshalMessage("GARBAGE|data\n")
	if result != nil {
		t.Errorf("UnmarshalMessage(\"GARBAGE|data\\n\") returned %v, want nil", result)
	}
}

// TestUnmarshalMalformedOffer tests that OFFER with wrong field count returns nil
func TestUnmarshalMalformedOffer(t *testing.T) {
	result := UnmarshalMessage("OFFER|onlytwopipes\n")
	if result != nil {
		t.Errorf("UnmarshalMessage(\"OFFER|onlytwopipes\\n\") returned %v, want nil", result)
	}
}

// TestUnmarshalMalformedOfferBadSize tests that OFFER with non-numeric size returns error
func TestUnmarshalMalformedOfferBadSize(t *testing.T) {
	result := UnmarshalMessage("OFFER|file|mime|notanumber\n")
	_, ok := result.(error)
	if !ok {
		t.Errorf("UnmarshalMessage(\"OFFER|file|mime|notanumber\\n\") returned %T, want error", result)
	}
}

// TestMakeOffer tests creating an offer from an actual file
func TestMakeOffer(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "test*.txt")
	if err != nil {
		t.Fatalf("CreateTemp() failed: %v", err)
	}
	defer tmpFile.Close()

	testData := []byte("test content")
	if _, err := tmpFile.Write(testData); err != nil {
		t.Fatalf("Write() failed: %v", err)
	}
	tmpFile.Close()

	offer, err := MakeOffer(tmpFile.Name())
	if err != nil {
		t.Fatalf("MakeOffer() failed: %v", err)
	}

	if offer.Type != "OFFER" {
		t.Errorf("Type = %q, want \"OFFER\"", offer.Type)
	}
	if offer.Filename != tmpFile.Name()[len(tmpDir)+1:] {
		t.Errorf("Filename = %q, want basename", offer.Filename)
	}
	if offer.Size != int64(len(testData)) {
		t.Errorf("Size = %d, want %d", offer.Size, len(testData))
	}
	if offer.Mimetype != "application/octet-stream" {
		t.Errorf("Mimetype = %q, want \"application/octet-stream\"", offer.Mimetype)
	}
}

// TestMakeOfferNonexistent tests that MakeOffer returns error for nonexistent file
func TestMakeOfferNonexistent(t *testing.T) {
	_, err := MakeOffer("/nonexistent/file/path/that/does/not/exist.txt")
	if err == nil {
		t.Errorf("MakeOffer() with nonexistent file returned nil, want error")
	}
}

// TestFormatSize tests size formatting for various byte counts
func TestFormatSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{500, "500 Bytes"},
		{1024, "1.00 KiB"},
		{1048576, "1.00 MiB"},
		{1073741824, "1.00 GiB"},
		{1099511627776, "1.00 TiB"},
		{0, "0 Bytes"},
		{512, "512 Bytes"},
		{2048, "2.00 KiB"},
		{5242880, "5.00 MiB"},
	}

	for _, tt := range tests {
		result := formatSize(tt.size)
		if result != tt.expected {
			t.Errorf("formatSize(%d) = %q, want %q", tt.size, result, tt.expected)
		}
	}
}
