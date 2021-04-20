package simple

import (
	"math/rand"
	"regexp"
	"strings"
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

func NamedId(name string, suffix int) (id string) {
	randId := RandomId(suffix)

	if name == "" {
		return randId
	}

	cleaner := regexp.MustCompile(`[^a-zA-B0-9-]`)

	id = name
	id = strings.ToLower(id)
	id = cleaner.ReplaceAllString(id, "_")
	id = strings.Trim(id, "_ -")

	id = id + "-" + randId

	return
}
