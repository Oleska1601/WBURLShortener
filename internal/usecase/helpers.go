package usecase

import (
	"math/rand"
	"time"
)

const (
	symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 6
)

func generateShortURL() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = symbols[rnd.Intn(len(symbols))]
	}

	return string(b)
}
