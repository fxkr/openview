package util

import (
	"crypto/rand"
	"math/big"
)

func NewPassword() (string, error) {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	const passwordLength = 32
	alphabetLength := big.NewInt(int64(len(alphabet)))
	result := make([]byte, 0, passwordLength)
	for i := 0; i < passwordLength; i++ {
		x, err := rand.Int(rand.Reader, alphabetLength)
		if err != nil {
			return "", err
		}
		result = append(result, byte(x.Uint64()))
	}
	return string(result), nil
}
