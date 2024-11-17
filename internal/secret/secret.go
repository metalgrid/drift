package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"net"

	"golang.org/x/crypto/curve25519"
)

type EncryptionKey *[32]byte

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

	// Generate a random nonce.
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Write the ephemeral public key and nonce to the underlying writer.
	if _, err := w.Write(publicKey[:]); err != nil {
		return nil, err
	}
	if _, err := w.Write(nonce); err != nil {
		return nil, err
	}

	return &EncryptWriter{
		writer: w,
		aead:   aead,
		nonce:  nonce,
		pubKey: publicKey,
	}, nil
}

// Write encrypts the data and writes it to the underlying writer.
func (ew *EncryptWriter) Write(data []byte) (int, error) {
	encrypted := ew.aead.Seal(nil, ew.nonce, data, nil)
	n, err := ew.writer.Write(encrypted)
	return n, err
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

	// Read the nonce.
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(r, nonce); err != nil {
		return nil, err
	}

	return &DecryptReader{
		reader: r,
		aead:   aead,
		nonce:  nonce,
	}, nil
}

// Read decrypts the data and reads from the underlying reader.
func (dr *DecryptReader) Read(buf []byte) (int, error) {
	n, err := dr.reader.Read(buf)
	if err != nil && err != io.EOF {
		return n, err
	}

	decrypted, err := dr.aead.Open(nil, dr.nonce, buf[:n], nil)
	if err != nil {
		return 0, err
	}

	copy(buf, decrypted)
	return len(decrypted), nil
}

// func main() {
// 	// Generate X25519 key pair for the recipient.
// 	privateKey, publicKey, _ := GenerateX25519KeyPair()

// 	// Data to encrypt.
// 	data := "This is a secret message!"

// 	// Encrypt the data.
// 	encBuf := &bytes.Buffer{}
// 	encWriter, _ := NewEncryptWriter(encBuf, publicKey)
// 	encWriter.Write([]byte(data))

// 	// Decrypt the data.
// 	decReader, _ := NewDecryptReader(encBuf, privateKey)
// 	decBuf := make([]byte, len(data))
// 	decReader.Read(decBuf)

// 	fmt.Println("Decrypted message:", string(decBuf))
// }

type WrappedConnection struct {
	net.Conn
	reader *DecryptReader
	writer *EncryptWriter
}

func (w *WrappedConnection) Read(p []byte) (n int, err error) {
	return w.reader.Read(p)
}

func (w *WrappedConnection) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
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
