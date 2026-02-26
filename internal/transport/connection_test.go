package transport

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/metalgrid/drift/internal/platform"
)

type mockGateway struct {
	mu            sync.Mutex
	notifications []string
}

func (m *mockGateway) Run(context.Context) error { return nil }
func (m *mockGateway) Shutdown()                 {}
func (m *mockGateway) NewRequest(string, string) error {
	return nil
}
func (m *mockGateway) Ask(string) string { return "" }
func (m *mockGateway) Notify(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.notifications = append(m.notifications, message)
}

func (m *mockGateway) hasNotification(value string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, notification := range m.notifications {
		if notification == value {
			return true
		}
	}
	return false
}

var _ platform.Gateway = (*mockGateway)(nil)

func newTCPConnPair(t *testing.T) (net.Conn, net.Conn) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed creating test listener: %v", err)
	}

	accepted := make(chan net.Conn, 1)
	acceptErr := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			acceptErr <- err
			return
		}
		accepted <- conn
	}()

	clientConn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		_ = listener.Close()
		t.Fatalf("failed dialing test listener: %v", err)
	}

	var serverConn net.Conn
	select {
	case serverConn = <-accepted:
	case err := <-acceptErr:
		_ = clientConn.Close()
		_ = listener.Close()
		t.Fatalf("failed accepting test connection: %v", err)
	case <-time.After(2 * time.Second):
		_ = clientConn.Close()
		_ = listener.Close()
		t.Fatal("timed out waiting for accepted connection")
	}

	if err := listener.Close(); err != nil {
		_ = serverConn.Close()
		_ = clientConn.Close()
		t.Fatalf("failed closing test listener: %v", err)
	}

	return serverConn, clientConn
}

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

func TestHandleConnectionAnswerAcceptWithoutFilenameState(t *testing.T) {
	gw := &mockGateway{}
	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, nil)
	}()

	if _, err := clientConn.Write(Accept().MarshalMessage()); err != nil {
		t.Fatalf("failed writing accept message: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if !gw.hasNotification(missingOutboundTransferStateToken) {
		t.Fatalf("expected notification %q", missingOutboundTransferStateToken)
	}

	if !gw.hasNotification(outboundAnswerNoSupportedPendingToken) {
		t.Fatalf("expected notification %q", outboundAnswerNoSupportedPendingToken)
	}
}

func TestHandleConnectionAnswerAcceptWithoutPendingOffer(t *testing.T) {
	gw := &mockGateway{}
	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	state := NewOutboundTransferState()
	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if _, err := clientConn.Write(Accept().MarshalMessage()); err != nil {
		t.Fatalf("failed writing accept message: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if !gw.hasNotification(unsolicitedAnswerIgnoredToken) {
		t.Fatalf("expected notification %q", unsolicitedAnswerIgnoredToken)
	}
}

func TestHandleConnectionAnswerAcceptSinglePendingFile(t *testing.T) {
	gw := &mockGateway{}
	state := NewOutboundTransferState()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "hello.txt")
	content := []byte("hello drift")
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}
	state.SetPendingFiles([]string{filePath})

	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if _, err := clientConn.Write(Accept().MarshalMessage()); err != nil {
		t.Fatalf("failed writing accept message: %v", err)
	}

	received := make([]byte, len(content))
	if _, err := io.ReadFull(clientConn, received); err != nil {
		t.Fatalf("failed reading sent file bytes: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if !bytes.Equal(received, content) {
		t.Fatalf("received bytes = %q, want %q", received, content)
	}

	if !gw.hasNotification("File sent: " + filePath) {
		t.Fatalf("expected file sent notification for %q", filePath)
	}
}

func TestHandleConnectionAnswerAcceptBatchPendingFiles(t *testing.T) {
	gw := &mockGateway{}
	state := NewOutboundTransferState()

	tmpDir := t.TempDir()
	firstPath := filepath.Join(tmpDir, "first.txt")
	secondPath := filepath.Join(tmpDir, "second.txt")
	firstContent := []byte("first")
	secondContent := []byte("second")

	if err := os.WriteFile(firstPath, firstContent, 0600); err != nil {
		t.Fatalf("failed creating first temp file: %v", err)
	}
	if err := os.WriteFile(secondPath, secondContent, 0600); err != nil {
		t.Fatalf("failed creating second temp file: %v", err)
	}

	state.SetPendingFiles([]string{firstPath, secondPath})

	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if _, err := clientConn.Write(Accept().MarshalMessage()); err != nil {
		t.Fatalf("failed writing accept message: %v", err)
	}

	combined := append(append([]byte(nil), firstContent...), secondContent...)
	received := make([]byte, len(combined))
	if _, err := io.ReadFull(clientConn, received); err != nil {
		t.Fatalf("failed reading sent batch bytes: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if !bytes.Equal(received, combined) {
		t.Fatalf("received bytes = %q, want %q", received, combined)
	}

	if !gw.hasNotification("Batch sent: 2 files") {
		t.Fatal("expected batch sent notification")
	}
}

func TestSendFileStateOrderingImmediateAnswer(t *testing.T) {
	gw := &mockGateway{}
	state := NewOutboundTransferState()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "race.txt")
	content := []byte("ordered-state")
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		t.Fatalf("failed creating temp file: %v", err)
	}

	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if err := SendFile(filePath, serverConn, state); err != nil {
		t.Fatalf("SendFile failed: %v", err)
	}

	reader := bufio.NewReader(clientConn)
	offer, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed reading offer from sender: %v", err)
	}
	if len(offer) == 0 {
		t.Fatal("expected non-empty offer")
	}

	if _, err := clientConn.Write(Accept().MarshalMessage()); err != nil {
		t.Fatalf("failed writing accept message: %v", err)
	}

	received := make([]byte, len(content))
	if _, err := io.ReadFull(reader, received); err != nil {
		t.Fatalf("failed reading sent file bytes: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if !bytes.Equal(received, content) {
		t.Fatalf("received bytes = %q, want %q", received, content)
	}

	if !gw.hasNotification("File sent: " + filePath) {
		t.Fatalf("expected file sent notification for %q", filePath)
	}
}

func TestHandleConnectionAnswerDeclineClearsPendingAndStops(t *testing.T) {
	gw := &mockGateway{}
	state := NewOutboundTransferState()
	state.SetPendingFiles([]string{"/tmp/declined.txt"})

	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if _, err := clientConn.Write(Decline().MarshalMessage()); err != nil {
		t.Fatalf("failed writing decline message: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}

	if files, ok := state.ConsumePendingFiles(); ok || len(files) != 0 {
		t.Fatal("expected pending files to be cleared after decline")
	}
}

func TestHandleConnectionMalformedAnswerMessage(t *testing.T) {
	gw := &mockGateway{}
	state := NewOutboundTransferState()

	serverConn, clientConn := newTCPConnPair(t)
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		HandleConnection(context.Background(), serverConn, gw, state)
	}()

	if _, err := clientConn.Write([]byte("ANSWER|ACCEPT|EXTRA\n")); err != nil {
		t.Fatalf("failed writing malformed answer message: %v", err)
	}
	_ = clientConn.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("HandleConnection did not return")
	}
}
