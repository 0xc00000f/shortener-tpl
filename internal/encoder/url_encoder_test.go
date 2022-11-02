package encoder

import (
	"testing"

	"github.com/0xc00000f/shortener-tpl/internal/rand"
	"github.com/stretchr/testify/assert"
)

func TestURLEncoder_Encode(t *testing.T) {
	r := rand.New(false)
	tests := []struct {
		name    string
		letters int
	}{
		{
			name:    "6 letters url",
			letters: 6,
		},
		{
			name:    "72 letters url",
			letters: 72,
		},
		{
			name:    "0 letters url",
			letters: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ue := URLEncoder{length: tt.letters, rand: r}
			url := ue.encode()
			assert.Equal(t, len(url), tt.letters)
		})
	}
}

func TestURLEncoder_Short(t *testing.T) {

	//var store = s.NewMemoryStorage(nil)
	//
	//tests := []struct {
	//	name   string
	//	length int
	//	long   string
	//	err    error
	//}{
	//	{
	//		name:   "positive #1",
	//		length: PreferredLength,
	//		long:   "https://google.com",
	//		err:    nil,
	//	},
	//	{
	//		name:   "positive #2",
	//		length: PreferredLength,
	//		long:   "https://dzen.ru/",
	//		err:    nil,
	//	},
	//	{
	//		name:   "negative #1 - empty long url",
	//		length: PreferredLength,
	//		long:   "",
	//		err:    s.ErrEmptyValue,
	//	},
	//	{
	//		name:   "negative #2 - empty short url",
	//		length: 0,
	//		long:   "https://ya.ru/",
	//		err:    s.ErrEmptyKey,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		ue := New(
	//			SetLength(tt.length),
	//			SetStorage(store),
	//			SetRandom(rand.New(true)),
	//		)
	//		short, err := ue.Short(tt.long)
	//		assert.Equal(t, tt.err, err)
	//		if err != nil {
	//			assert.Equal(t, 0, len(short))
	//			return
	//		}
	//		assert.Equal(t, tt.length, len(short))
	//
	//		long, err := ue.Get(short)
	//		assert.Nil(t, err)
	//		assert.Equal(t, long, tt.long)
	//	})
	//}
}

func TestURLEncoder_Get(t *testing.T) {

	//var storage = s.NewMemoryStorage(nil)
	//storage.Store("ytAA2Z", "https://google.com")
	//storage.Store("hNaU8l", "https://dzen.ru/")
	//
	//tests := []struct {
	//	name  string
	//	short string
	//	long  string
	//	err   error
	//}{
	//	{
	//		name:  "positive #1",
	//		short: "ytAA2Z",
	//		long:  "https://google.com",
	//		err:   nil,
	//	},
	//	{
	//		name:  "positive #2",
	//		short: "hNaU8l",
	//		long:  "https://dzen.ru/",
	//		err:   nil,
	//	},
	//	{
	//		name:  "negative #1 - key is not exist",
	//		short: "not exist",
	//		long:  "https://dzen.ru/",
	//		err:   s.ErrNoKeyFound,
	//	},
	//	{
	//		name:  "negative #2 - empty short",
	//		short: "",
	//		long:  "https://dzen.ru/",
	//		err:   s.ErrEmptyKey,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		ue := New(
	//			SetLength(PreferredLength),
	//			SetStorage(storage))
	//		long, err := ue.Get(tt.short)
	//		require.Equal(t, tt.err, err)
	//		if err != nil {
	//			return
	//		}
	//		assert.Equal(t, tt.long, long)
	//
	//	})
	//}
}
