package znox

import (
	"crypto/rand"
	b64 "encoding/base64"
	"os"
)

func GenerateAESkey() error {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return err
	}

	keyB64 := b64.StdEncoding.EncodeToString(key)

	err = os.WriteFile("key", []byte(keyB64), 0777)
	if err != nil {
		return err
	}
	return nil
}

func ReadAESkey() ([]byte, error) {
	keyB64, err := os.ReadFile("key")
	if err != nil {
		return nil, err
	}
	key, err := b64.StdEncoding.DecodeString(string(keyB64))
	if err != nil {
		return nil, err
	}
	return key, nil
}
