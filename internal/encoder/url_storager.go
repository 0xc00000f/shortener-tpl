package encoder

type URLStorager interface {
	Get(short string) (string, error)
	GetAll() (result map[string]string, err error)
	Store(short string, long string) error
	IsKeyExist(short string) (bool, error)
}
