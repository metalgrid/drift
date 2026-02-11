package transport

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestStoreFileSuccess tests storing a file with content
func TestStoreFileSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	testData := []byte("test file content")
	reader := bytes.NewReader(testData)

	err := storeFile(tmpDir, "test.txt", int64(len(testData)), reader, nil)
	if err != nil {
		t.Fatalf("storeFile() failed: %v", err)
	}

	// Verify file exists with correct content
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir() failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file in directory, got %d", len(files))
	}

	filePath := filepath.Join(tmpDir, "test.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	if !bytes.Equal(content, testData) {
		t.Errorf("File content = %q, want %q", content, testData)
	}
}

// TestStoreFileZeroBytes tests storing a file with zero bytes
func TestStoreFileZeroBytes(t *testing.T) {
	tmpDir := t.TempDir()
	reader := bytes.NewReader([]byte{})

	err := storeFile(tmpDir, "empty.txt", 0, reader, nil)
	if err != nil {
		t.Fatalf("storeFile() with zero bytes failed: %v", err)
	}

	// Verify file exists with 0 bytes
	filePath := filepath.Join(tmpDir, "empty.txt")
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Stat() failed: %v", err)
	}

	if fileInfo.Size() != 0 {
		t.Errorf("File size = %d, want 0", fileInfo.Size())
	}
}

// TestStoreFileDirCreation tests that storeFile creates non-existent directories
func TestStoreFileDirCreation(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir", "nested")
	testData := []byte("nested file content")
	reader := bytes.NewReader(testData)

	err := storeFile(subDir, "nested.txt", int64(len(testData)), reader, nil)
	if err != nil {
		t.Fatalf("storeFile() with nested directory failed: %v", err)
	}

	// Verify directory was created and file exists
	filePath := filepath.Join(subDir, "nested.txt")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() failed: %v", err)
	}

	if !bytes.Equal(content, testData) {
		t.Errorf("File content = %q, want %q", content, testData)
	}
}
