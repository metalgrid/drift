package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/curve25519"
)

type EncryptionKey *[32]byte

const maxEncryptedFrameSize = 8 * 1024 * 1024

// GenerateX25519KeyPair generates an X25519 key pair.
func GenerateX25519KeyPair() (EncryptionKey, EncryptionKey, error) {
	var privateKey, publicKey [32]byte

	if _, err := rand.Read(privateKey[:]); err != nil {
		return nil, nil, err
	}
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return &privateKey, &publicKey, nil
}

// DeriveSharedSecret derives a shared secret using the X25519 private key and the peer's public key.
func DeriveSharedSecret(privateKey, peerPublicKey EncryptionKey) (EncryptionKey, error) {
	var sharedSecret [32]byte
	curve25519.ScalarMult(&sharedSecret, privateKey, peerPublicKey)
	if sharedSecret == ([32]byte{}) {
		return nil, errors.New("invalid shared secret")
	}

	return &sharedSecret, nil
}

// KeyDerivation derives a 32-byte AES key from the shared secret using SHA-256.
func KeyDerivation(sharedSecret EncryptionKey) []byte {
	hash := sha256.Sum256(sharedSecret[:])
	return hash[:]
}

// EncryptWriter encrypts data using AES-GCM and writes to the underlying writer.
type EncryptWriter struct {
	writer io.Writer
	aead   cipher.AEAD
	nonce  []byte
	pubKey EncryptionKey
	encBuf []byte
}

// NewEncryptWriter initializes an EncryptWriter with the recipient's public X25519 key.
func NewEncryptWriter(w io.Writer, recipientPublicKey EncryptionKey) (*EncryptWriter, error) {
	// Generate an ephemeral X25519 key pair.
	privateKey, publicKey, err := GenerateX25519KeyPair()
	if err != nil {
		return nil, err
	}

	// Derive a shared secret using the ephemeral private key and recipient's public key.
	sharedSecret, err := DeriveSharedSecret(privateKey, recipientPublicKey)
	if err != nil {
		return nil, err
	}

	// Derive an AES key from the shared secret.
	aesKey := KeyDerivation(sharedSecret)

	// Create AES-GCM cipher.
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Initialize sequential nonce.
	nonce := make([]byte, aead.NonceSize())

	// Write the ephemeral public key to the underlying writer.
	if _, err := w.Write(publicKey[:]); err != nil {
		return nil, err
	}

	return &EncryptWriter{
		writer: w,
		aead:   aead,
		nonce:  nonce,
		pubKey: publicKey,
	}, nil
}

func (ew *EncryptWriter) Write(data []byte) (int, error) {
	// Encrypt the data.
	encrypted := ew.aead.Seal(nil, ew.nonce, data, nil)
	incrementNonce(ew.nonce)

	// Write the length of the encrypted data (4 bytes, big-endian).
	length := uint32(len(encrypted))
	lengthBuf := make([]byte, 4)
	lengthBuf[0] = byte(length >> 24)
	lengthBuf[1] = byte(length >> 16)
	lengthBuf[2] = byte(length >> 8)
	lengthBuf[3] = byte(length)

	if _, err := ew.writer.Write(lengthBuf); err != nil {
		return 0, err
	}

	// Write the encrypted data.
	_, err := ew.writer.Write(encrypted)

	// Return the original length of the data not to confuse the caller.
	return len(data), err
}

// DecryptReader decrypts data using AES-GCM and reads from the underlying reader.
type DecryptReader struct {
	reader io.Reader
	aead   cipher.AEAD
	nonce  []byte
	decBuf []byte
}

// NewDecryptReader initializes a DecryptReader with the recipient's private X25519 key.
func NewDecryptReader(r io.Reader, recipientPrivateKey EncryptionKey) (*DecryptReader, error) {
	// Read the ephemeral public key.
	var publicKey [32]byte
	if _, err := io.ReadFull(r, publicKey[:]); err != nil {
		return nil, err
	}

	// Derive a shared secret using the recipient's private key and ephemeral public key.
	sharedSecret, err := DeriveSharedSecret(recipientPrivateKey, &publicKey)
	if err != nil {
		return nil, err
	}

	// Derive an AES key from the shared secret.
	aesKey := KeyDerivation(sharedSecret)

	// Create AES-GCM cipher.
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Initialize sequential nonce.
	nonce := make([]byte, aead.NonceSize())

	return &DecryptReader{
		reader: r,
		aead:   aead,
		nonce:  nonce,
	}, nil
}

func (dr *DecryptReader) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}

	if len(dr.decBuf) > 0 {
		n := copy(buf, dr.decBuf)
		dr.decBuf = dr.decBuf[n:]
		return n, nil
	}

	// Read the length of the encrypted data (4 bytes, big-endian).
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(dr.reader, lengthBuf); err != nil {
		return 0, err
	}
	length := (uint32(lengthBuf[0]) << 24) |
		(uint32(lengthBuf[1]) << 16) |
		(uint32(lengthBuf[2]) << 8) |
		uint32(lengthBuf[3])
	if length == 0 || length > maxEncryptedFrameSize {
		return 0, fmt.Errorf("encrypted frame length %d out of bounds", length)
	}

	// Allocate a buffer to hold the encrypted data.
	encrypted := make([]byte, length)
	if _, err := io.ReadFull(dr.reader, encrypted); err != nil {
		return 0, err
	}

	// Decrypt the data.
	decrypted, decryptErr := dr.aead.Open(nil, dr.nonce, encrypted, nil)
	if decryptErr != nil {
		return 0, fmt.Errorf("decryption error: %w", decryptErr)
	}

	// Increment the nonce after each decryption.
	incrementNonce(dr.nonce)

	// Copy the decrypted data into the provided buffer.
	n := copy(buf, decrypted)
	if n < len(decrypted) {
		dr.decBuf = append(dr.decBuf[:0], decrypted[n:]...)
	}

	return n, nil
}

type WrappedConnection struct {
	net.Conn
	reader *DecryptReader
	writer *EncryptWriter
}

func (w *WrappedConnection) Read(p []byte) (n int, err error) {
	n, err = w.reader.Read(p)
	if err != nil {
		return n, fmt.Errorf("read error: %w", err)
	}
	return n, err
}

func (w *WrappedConnection) Write(p []byte) (n int, err error) {
	n, err = w.writer.Write(p)
	if err != nil {
		return n, fmt.Errorf("write error: %w", err)
	}
	return n, err
}

func (w *WrappedConnection) Close() error {
	fmt.Println("Closing connection")
	return w.Conn.Close()
}

// SecureConnection establishes a secure connection with the peer by using the *remote* public key and the *local* private key
func SecureConnection(conn net.Conn, peerPublicKey, localPrivateKey EncryptionKey) (net.Conn, error) {
	writer, err := NewEncryptWriter(conn, peerPublicKey)
	if err != nil {
		return nil, err
	}

	reader, err := NewDecryptReader(conn, localPrivateKey)
	if err != nil {
		return nil, err
	}

	return &WrappedConnection{
		conn,
		reader,
		writer,
	}, nil

}

func incrementNonce(nonce []byte) {
	for i := len(nonce) - 1; i >= 0; i-- {
		nonce[i]++
		if nonce[i] != 0 {
			break
		}
	}
}
