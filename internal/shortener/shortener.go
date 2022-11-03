package shortener

type Shortener interface {
	Short(long string) (short string, err error)
	Get(short string) (long string, err error)
	GetAll() (result map[string]string, err error)
}
