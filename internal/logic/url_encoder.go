package logic

import (
	"github.com/0xc00000f/shortener-tpl/internal/utils"
)

var preferredLength = 6

type URLEncoder struct {
	length  int
	storage URLStorager
}

type Option func(ue *URLEncoder)

func NewURLEncoder(options ...Option) *URLEncoder {
	ue := URLEncoder{length: preferredLength}

	for _, fn := range options {
		fn(&ue)
	}

	return &ue
}

func SetLength(length int) Option {
	return func(ue *URLEncoder) {
		ue.length = length
	}
}

func SetStorage(s URLStorager) Option {
	return func(ue *URLEncoder) {
		ue.storage = s
	}
}

func (ue *URLEncoder) encode() string {
	return utils.RandStringRunes(ue.length)
}

func (ue *URLEncoder) Short(long string) (short string, err error) {
	for {
		short = ue.encode()
		if ue.storage.IsKeyExist(short) {
			continue
		}
		break
	}

	err = ue.storage.Store(short, long)
	if err != nil {
		return "", err
	}
	return short, nil
}

func (ue *URLEncoder) Get(short string) (long string, err error) {
	return ue.storage.Get(short)
}