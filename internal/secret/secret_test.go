package secret

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func TestDecryptReaderRespectsSmallDestinationBuffers(t *testing.T) {
	recipientPriv, recipientPub, err := GenerateX25519KeyPair()
	if err != nil {
		t.Fatalf("GenerateX25519KeyPair() failed: %v", err)
	}

	var wire bytes.Buffer
	ew, err := NewEncryptWriter(&wire, recipientPub)
	if err != nil {
		t.Fatalf("NewEncryptWriter() failed: %v", err)
	}

	plaintext := []byte("hello-world-plaintext")
	if _, err := ew.Write(plaintext); err != nil {
		t.Fatalf("EncryptWriter.Write() failed: %v", err)
	}

	dr, err := NewDecryptReader(bytes.NewReader(wire.Bytes()), recipientPriv)
	if err != nil {
		t.Fatalf("NewDecryptReader() failed: %v", err)
	}

	var out bytes.Buffer
	buf := make([]byte, 4)
	for out.Len() < len(plaintext) {
		n, err := dr.Read(buf)
		if err != nil && err != io.EOF {
			t.Fatalf("DecryptReader.Read() failed: %v", err)
		}
		if n == 0 {
			break
		}
		if n > len(buf) {
			t.Fatalf("DecryptReader.Read() returned n=%d > len(buf)=%d", n, len(buf))
		}
		out.Write(buf[:n])
	}

	if !bytes.Equal(out.Bytes(), plaintext) {
		t.Fatalf("decrypted plaintext = %q, want %q", out.Bytes(), plaintext)
	}
}

func TestDecryptReaderRejectsOversizedFrameBeforeAllocation(t *testing.T) {
	var frame bytes.Buffer
	_ = binary.Write(&frame, binary.BigEndian, uint32(maxEncryptedFrameSize+1))

	dr := &DecryptReader{reader: bytes.NewReader(frame.Bytes())}
	buf := make([]byte, 32)

	if _, err := dr.Read(buf); err == nil {
		t.Fatal("expected oversized frame error, got nil")
	}
}
