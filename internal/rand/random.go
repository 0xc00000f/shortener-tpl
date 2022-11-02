package rand

import (
	"math/rand"
	"time"
)

type Random struct {
	rand *rand.Rand
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
	rand := rand.New(source)
	return Random{rand: rand}
}

func (r Random) String(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.rand.Intn(len(letterRunes))]
	}
	return string(b)
}
