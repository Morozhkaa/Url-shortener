package usecases

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateKey() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 5
	var keyBuilder strings.Builder
	for i := 0; i < keyLength; i++ {
		char := alphabet[r.Intn(len(alphabet))]
		keyBuilder.WriteByte(char)
	}
	return keyBuilder.String()
}
