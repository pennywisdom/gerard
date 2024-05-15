package cmd

import (
	"crypto/rand"
	"encoding/hex"
)

func generateRandom10Char() (string, error) {
	b := make([]byte, 5) // Generate 5 random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	s := hex.EncodeToString(b) // Encode bytes to string
	return s, nil              // Outputs a 10 character long string
}
