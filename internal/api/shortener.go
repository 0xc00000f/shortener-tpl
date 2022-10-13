package api

type ShortenerApi struct {
	logic Shortener
}

func NewShortenerApi(logic Shortener) *ShortenerApi {
	return &ShortenerApi{logic: logic}
}

func (sa *ShortenerApi) Logic() Shortener {
	return sa.logic
}
