package transport

import (
	"bytes"
	"testing"
)

// TestProgressWriterReportsBytesCorrectly verifies ProgressWriter calls callback with correct byte counts
func TestProgressWriterReportsBytesCorrectly(t *testing.T) {
	buf := &bytes.Buffer{}
	var reported []int64
	callback := func(bytesTransferred, totalBytes int64) {
		reported = append(reported, bytesTransferred)
	}

	pw := NewProgressWriter(buf, 10, callback)
	pw.Write([]byte("hello"))
	pw.Write([]byte("world"))

	if len(reported) != 2 {
		t.Errorf("expected 2 callbacks, got %d", len(reported))
	}
	if reported[0] != 5 {
		t.Errorf("expected first callback with 5 bytes, got %d", reported[0])
	}
	if reported[1] != 10 {
		t.Errorf("expected second callback with 10 bytes, got %d", reported[1])
	}
	if buf.String() != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", buf.String())
	}
}

// TestProgressReaderReportsBytesCorrectly verifies ProgressReader calls callback with correct byte counts
func TestProgressReaderReportsBytesCorrectly(t *testing.T) {
	data := []byte("hello world")
	reader := bytes.NewReader(data)
	var reported []int64
	callback := func(bytesTransferred, totalBytes int64) {
		reported = append(reported, bytesTransferred)
	}

	pr := NewProgressReader(reader, int64(len(data)), callback)
	buf := make([]byte, 5)
	pr.Read(buf)
	pr.Read(buf)

	if len(reported) != 2 {
		t.Errorf("expected 2 callbacks, got %d", len(reported))
	}
	if reported[0] != 5 {
		t.Errorf("expected first callback with 5 bytes, got %d", reported[0])
	}
	if reported[1] != 10 {
		t.Errorf("expected second callback with 10 bytes, got %d", reported[1])
	}
}

// TestProgressWriterNilCallback verifies ProgressWriter handles nil callback gracefully
func TestProgressWriterNilCallback(t *testing.T) {
	buf := &bytes.Buffer{}
	pw := NewProgressWriter(buf, 10, nil)
	n, err := pw.Write([]byte("hello"))

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes written, got %d", n)
	}
	if buf.String() != "hello" {
		t.Errorf("expected 'hello', got %q", buf.String())
	}
}

// TestProgressReaderNilCallback verifies ProgressReader handles nil callback gracefully
func TestProgressReaderNilCallback(t *testing.T) {
	data := []byte("hello")
	reader := bytes.NewReader(data)
	pr := NewProgressReader(reader, int64(len(data)), nil)
	buf := make([]byte, 5)
	n, err := pr.Read(buf)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Errorf("expected 5 bytes read, got %d", n)
	}
}

// TestProgressMonotonicallyIncreasing verifies bytes transferred never decreases
func TestProgressMonotonicallyIncreasing(t *testing.T) {
	buf := &bytes.Buffer{}
	var reported []int64
	callback := func(bytesTransferred, totalBytes int64) {
		reported = append(reported, bytesTransferred)
	}

	pw := NewProgressWriter(buf, 100, callback)
	for i := 0; i < 10; i++ {
		pw.Write([]byte("x"))
	}

	for i := 1; i < len(reported); i++ {
		if reported[i] < reported[i-1] {
			t.Errorf("bytes decreased at index %d: %d -> %d", i, reported[i-1], reported[i])
		}
	}
}

// TestProgressWriterTotalBytesParameter verifies totalBytes is passed to callback
func TestProgressWriterTotalBytesParameter(t *testing.T) {
	buf := &bytes.Buffer{}
	var totalBytes int64
	callback := func(bytesTransferred, total int64) {
		totalBytes = total
	}

	pw := NewProgressWriter(buf, 42, callback)
	pw.Write([]byte("test"))

	if totalBytes != 42 {
		t.Errorf("expected totalBytes=42, got %d", totalBytes)
	}
}

// TestProgressReaderTotalBytesParameter verifies totalBytes is passed to callback
func TestProgressReaderTotalBytesParameter(t *testing.T) {
	data := []byte("test")
	reader := bytes.NewReader(data)
	var totalBytes int64
	callback := func(bytesTransferred, total int64) {
		totalBytes = total
	}

	pr := NewProgressReader(reader, 42, callback)
	buf := make([]byte, 4)
	pr.Read(buf)

	if totalBytes != 42 {
		t.Errorf("expected totalBytes=42, got %d", totalBytes)
	}
}
