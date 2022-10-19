package api

import "os"

type ShortenerAPI struct {
	logic   Shortener
	BaseURL string
}

type Option func(sa *ShortenerAPI)

func NewShortenerAPI(options ...Option) *ShortenerAPI {
	sa := ShortenerAPI{}

	for _, fn := range options {
		fn(&sa)
	}

	return &sa
}

func (sa *ShortenerAPI) Logic() Shortener {
	return sa.logic
}

func SetLogic(logic Shortener) Option {
	return func(sa *ShortenerAPI) {
		sa.logic = logic
	}
}

func InitBaseURL(baseURL string) Option {
	return func(sa *ShortenerAPI) {
		if len(baseURL) > 0 {
			sa.BaseURL = baseURL
			return
		}

		sa.BaseURL = os.Getenv("BASE_URL")
	}
}
