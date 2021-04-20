package simple

import (
	"math/rand"
	"time"
)

func RandomId(length int) (id string) {

	charset := "0123456789abcdef"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}
