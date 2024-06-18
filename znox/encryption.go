package znox

import (
	"bytes"
	cryptorand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	SaltSize   = 32         // in bytes
	NonceSize  = 24         // in bytes. taken from aead.NonceSize()
	KeySize    = uint32(32) // KeySize is 32 bytes (256 bits).
	KeyTime    = uint32(5)
	KeyMemory  = uint32(1024 * 64) // KeyMemory in KiB. here, 64 MiB.
	KeyThreads = uint8(4)
	chunkSize  = 1024 * 32 // chunkSize in bytes. here, 32 KiB.
)

func Encryption(unencryptedSrc, encryptedDst string, passw1, passw2 []byte) error {
	if !bytes.Equal(passw1, passw2) {
		return errors.New("passwords don't match")
	}

	salt := make([]byte, SaltSize)
	if n, err := cryptorand.Read(salt); err != nil || n != SaltSize {
		return fmt.Errorf("error reaind salt: %s", err.Error())
	}

	outfile, err := os.OpenFile(encryptedDst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error when opening/creating output file: %s", err)
	}
	defer outfile.Close()

	outfile.Write(salt)

	key := argon2.IDKey(passw1, salt, KeyTime, KeyMemory, KeyThreads, KeySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return fmt.Errorf("error creating cipher: %s", err.Error())
	}

	infile, err := os.Open(unencryptedSrc)
	if err != nil {
		return fmt.Errorf("error when opening input file:", err.Error())
	}
	defer infile.Close()

	buf := make([]byte, chunkSize)
	ad_counter := 0 // associated data is a counter

	for {
		n, err := infile.Read(buf)

		if n > 0 {
			// Select a random nonce, and leave capacity for the ciphertext.
			nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+n+aead.Overhead())
			if m, err := cryptorand.Read(nonce); err != nil || m != aead.NonceSize() {
				return fmt.Errorf("error when generating random nonce (size=%d): %s", m, err.Error())
			}

			msg := buf[:n]
			// Encrypt the message and append the ciphertext to the nonce.
			encryptedMsg := aead.Seal(nonce, nonce, msg, []byte(string(ad_counter)))
			outfile.Write(encryptedMsg)
			ad_counter += 1
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error when reading input file chunk: %s", err.Error())
		}
	}
	return nil
}

func Decryption(encryptedSrc, decryptedDst string, passw []byte) error {
	infile, err := os.Open(encryptedSrc)
	if err != nil {
		return fmt.Errorf("error when opening input file: %s", err.Error())
	}
	defer infile.Close()

	salt := make([]byte, SaltSize)
	n, err := infile.Read(salt)
	if n != SaltSize {
		return errors.New("invalid salt length")
	}
	if err == io.EOF {
		return errors.New("encountered EOF error")
	}
	if err != nil {
		return fmt.Errorf("error encountered: %s", err.Error())
	}

	key := argon2.IDKey(passw, salt, KeyTime, KeyMemory, KeyThreads, KeySize)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return fmt.Errorf("error initiating encryption: %s", err.Error())
	}

	decbufsize := aead.NonceSize() + chunkSize + aead.Overhead()

	outfile, err := os.OpenFile(decryptedDst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("error when opening output file: %s", err.Error())
	}
	defer outfile.Close()

	buf := make([]byte, decbufsize)
	ad_counter := 0 // associated data is a counter

	for {
		n, err := infile.Read(buf)
		if n > 0 {
			encryptedMsg := buf[:n]
			if len(encryptedMsg) < aead.NonceSize() {
				return errors.New("error: ciphertext is too short")
			}

			// Split nonce and ciphertext.
			nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]
			// Decrypt the message and check it wasn't tampered with.
			plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(string(ad_counter)))
			if err != nil {
				return errors.New("error when decrypting ciphertext. May be wrong password or file is damaged")
			}

			outfile.Write(plaintext)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error encountered. Read %d bytes: %v", n, err)
		}

		ad_counter += 1
	}
	return nil
}

func returnBaseName(path string) string {
	baseName := filepath.Base(path)
	if strings.Contains(baseName, ".") {
		idx := strings.Index(baseName, ".")
		baseName = baseName[:idx]
	}
	return baseName
}
