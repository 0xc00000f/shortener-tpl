package storage

import "errors"

type MemoryStorage map[string]string

func NewMemoryStorage() MemoryStorage {
	return make(MemoryStorage)
}

func (ms MemoryStorage) Get(short string) (long string, err error) {
	if len(short) == 0 {
		err = errors.New("empty string as a key isn't allowed")
		return "", err
	}
	long, ok := ms[short]
	if !ok {
		return "", errors.New("key is not exist")
	}
	return long, nil
}

func (ms MemoryStorage) Store(short, long string) error {
	if len(short) == 0 {
		return errors.New("empty string as a key isn't allowed")
	}
	if len(long) == 0 {
		return errors.New("empty string as a value isn't allowed")
	}
	ms[short] = long
	return nil
}

func (ms MemoryStorage) IsKeyExist(short string) (bool, error) {
	_, ok := ms[short]
	return ok, nil
}
