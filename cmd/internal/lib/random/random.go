package random

import (
	"math/rand"
	"sync"
	"time"
)

var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu sync.Mutex
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var chars = []rune(alphabet)

func NewRandomString(size int) string {
	b := make([]rune, size)
	rngMu.Lock()
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	rngMu.Unlock()
	return string(b)
}
