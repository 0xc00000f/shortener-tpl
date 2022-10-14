package storage

import "errors"

type MemoryStorage map[string]string

func NewStorage() MemoryStorage {
	return make(MemoryStorage)
}

func (ds MemoryStorage) Get(short string) (value string, err error) {
	if len(short) == 0 {
		err = errors.New("empty string as a key isn't allowed")
		return "", err
	}
	value, ok := ds[short]
	if !ok {
		return "", errors.New("key is not exist")
	}
	return value, nil
}

func (ds MemoryStorage) Store(short, long string) error {
	if len(short) == 0 {
		return errors.New("empty string as a key isn't allowed")
	}
	if len(long) == 0 {
		return errors.New("empty string as a value isn't allowed")
	}
	ds[short] = long
	return nil
}

func (ds MemoryStorage) IsKeyExist(short string) bool {
	_, ok := ds[short]
	return ok
}
