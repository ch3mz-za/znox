package znox

import (
	"crypto/rand"
	b64 "encoding/base64"
	"io/ioutil"
	"log"
	"os"
)

func GenerateAESkey() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Failed to generate key")
	}

	keyB64 := b64.StdEncoding.EncodeToString(key)

	err = os.WriteFile("key", []byte(keyB64), 0777)
	if err != nil {
		log.Fatal("Failed to write key")
	}
}

func ReadAESkey() []byte {
	keyB64, err := ioutil.ReadFile("key")
	if err != nil {
		log.Fatal(err)
	}
	key, err := b64.StdEncoding.DecodeString(string(keyB64))
	if err != nil {
		log.Fatal(err)
	}
	return key
}
