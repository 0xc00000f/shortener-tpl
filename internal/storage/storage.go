package storage

type URLStorage interface {
	Get(string) (string, bool)
	Set(string, string)
}

type dataStorage map[string]string

var Storage dataStorage

func init() {
	Storage = newStorage()
}

func newStorage() dataStorage {
	return make(dataStorage)
}

func (ds dataStorage) Get(key string) (value string, ok bool) {
	value, ok = ds[key]
	return
}

func (ds dataStorage) Set(key, value string) {
	ds[key] = value
}
