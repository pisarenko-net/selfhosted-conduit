package router

import (
	"crypto/rand"
	"math/big"
)

const CODE_LENGTH int = 6

func generateRandomCode() string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result string
	for i := 0; i < CODE_LENGTH; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result += string(charset[randomIndex.Int64()])
	}
	return result
}
