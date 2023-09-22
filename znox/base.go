package znox

import (
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/term"
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

func Encryption(unencryptedSrc string, encryptedDst string) {
	fmt.Println("Encrypting.\nEnter a long and random password : ")
	bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Println("Error when reading password from terminal")
		panic(err)
	}

	fmt.Println("Enter the same password again : ")
	bytepw2, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Println("Error when reading password2 from terminal.")
		panic(err)
	}

	if !bytes.Equal(bytepw, bytepw2) {
		log.Println("Passwords don't match! Exiting.")
		os.Exit(1)
	}

	salt := make([]byte, SaltSize)
	if n, err := cryptorand.Read(salt); err != nil || n != SaltSize {
		log.Println("Error when generating radom salt.")
		panic(err)
	}

	outfile, err := os.OpenFile(encryptedDst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error when opening/creating output file.")
		panic(err)
	}
	defer outfile.Close()

	outfile.Write(salt)

	key := argon2.IDKey(bytepw, salt, KeyTime, KeyMemory, KeyThreads, KeySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		log.Println("Error when creating cipher.")
		panic(err)
	}

	infile, err := os.Open(unencryptedSrc)
	if err != nil {
		log.Println("Error when opening input file.")
		panic(err)
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
				log.Println("Error when generating random nonce :", err)
				log.Println("Generated nonce is of following size. m : ", m)
				panic(err)
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
			log.Println("Error when reading input file chunk :", err)
			panic(err)
		}
	}
}

func Decryption(encryptedSrc string, decryptedDst string) {
	fmt.Println("Decrypting.\nEnter the password : ")
	bytepw, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Println("Error when reading password from terminal.")
		panic(err)
	}

	infile, err := os.Open(encryptedSrc)
	if err != nil {
		log.Println("Error when opening input file.")
		panic(err)
	}
	defer infile.Close()

	salt := make([]byte, SaltSize)
	n, err := infile.Read(salt)
	if n != SaltSize {
		log.Printf("Error. Salt should be %d bytes long. salt n : %d", SaltSize, n)
		log.Print("Another Error :", err)
		panic("Generated salt is not of required length")
	}
	if err == io.EOF {
		log.Println("Encountered EOF error.")
		panic(err)
	}
	if err != nil {
		log.Println("Error encountered :", err)
		panic(err)
	}

	key := argon2.IDKey(bytepw, salt, KeyTime, KeyMemory, KeyThreads, KeySize)
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		log.Println("Error initiating encryption.")
		panic(err)
	}

	decbufsize := aead.NonceSize() + chunkSize + aead.Overhead()

	outfile, err := os.OpenFile(decryptedDst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error when opening output file.")
		panic(err)
	}
	defer outfile.Close()

	buf := make([]byte, decbufsize)
	ad_counter := 0 // associated data is a counter

	for {
		n, err := infile.Read(buf)
		if n > 0 {
			encryptedMsg := buf[:n]
			if len(encryptedMsg) < aead.NonceSize() {
				log.Println("Error. Ciphertext is too short.")
				panic("Ciphertext too short")
			}

			// Split nonce and ciphertext.
			nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]
			// Decrypt the message and check it wasn't tampered with.
			plaintext, err := aead.Open(nil, nonce, ciphertext, []byte(string(ad_counter)))
			if err != nil {
				log.Println("Error when decrypting ciphertext. May be wrong password or file is damaged.")
				panic(err)
			}

			outfile.Write(plaintext)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error encountered. Read %d bytes: %v", n, err)
			panic(err)
		}

		ad_counter += 1
	}
}

func returnBaseName(path string) string {
	baseName := filepath.Base(path)
	if strings.Contains(baseName, ".") {
		idx := strings.Index(baseName, ".")
		baseName = baseName[:idx]
	}
	return baseName
}
