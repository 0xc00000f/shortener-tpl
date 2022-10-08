package storage

type URLStorage interface {
	Get(string) (string, bool)
	Set(string, string)
}

type DataStorage map[string]string

func NewStorage() DataStorage {
	return make(DataStorage)
}

func (ds DataStorage) Get(key string) (value string, ok bool) {
	value, ok = ds[key]
	return
}

func (ds DataStorage) Set(key, value string) {
	ds[key] = value
}
