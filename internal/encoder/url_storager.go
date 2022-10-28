package encoder

type URLStorager interface {
	Get(short string) (string, error)
	Store(short string, long string) error
	IsKeyExist(short string) (bool, error)
}
