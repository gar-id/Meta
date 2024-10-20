package tools

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"#*?!&^%$@"

func StringWithCharset(length int, charset string) string {
	getTime := time.Now().UnixNano()
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(getTime))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}
