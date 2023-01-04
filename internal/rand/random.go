package rand

import (
	"math/rand"
	"sync"
	"time"
)

type Random struct {
	rand *rand.Rand

	mu *sync.Mutex
}

func New(predictable bool) Random {
	var seed int64

	switch predictable {
	case true:
		seed = 1

	case false:
		seed = time.Now().UnixNano()
	}

	source := rand.NewSource(seed)
	random := rand.New(source) //nolint:gosec

	return Random{rand: random, mu: &sync.Mutex{}}
}

func (r *Random) String(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)

	r.mu.Lock()
	for i := range b {
		b[i] = letterRunes[r.rand.Intn(len(letterRunes))]
	}
	r.mu.Unlock()

	return string(b)
}
