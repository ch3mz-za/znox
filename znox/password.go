package znox

import (
	cryptorand "crypto/rand"
	"log"
	"math/big"
)

func GeneratePassword(pwLength int) string {
	smallAlpha := "abcdefghijklmnopqrstuvwxyz"
	bigAlpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	specialChars := "`~!@#$%^&*()_+-={}|[]\\;':\",./<>?"

	letters := smallAlpha + bigAlpha + digits + specialChars

	pw := ""
	for i := 0; i < pwLength; i++ {
		pw += string(letters[getRandNum(int64(len(letters)))])
	}

	return pw
}

func getRandNum(max int64) int64 {
	if i, err := cryptorand.Int(cryptorand.Reader, big.NewInt(max)); err != nil {
		log.Println("Error when generating random num : ", err)
		panic(err)
	} else {
		return i.Int64()
	}
}
