package api

type Shortener interface {
	Short(long string) (short string, err error)
	Get(short string) (long string, err error)
}
