package api

type ShortenerAPI struct {
	logic Shortener
}

func NewShortenerAPI(logic Shortener) *ShortenerAPI {
	return &ShortenerAPI{logic: logic}
}

func (sa *ShortenerAPI) Logic() Shortener {
	return sa.logic
}
