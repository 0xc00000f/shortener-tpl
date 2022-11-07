package storage

import (
	"bufio"
	"encoding/json"
	"github.com/google/uuid"
	"os"

	"go.uber.org/zap"
)

type FileStorage struct {
	file   *os.File
	memory MemoryStorage

	l *zap.Logger
}

type url struct {
	UserID uuid.UUID `json:"userID,omitempty"`
	Short  string    `json:"short"`
	Long   string    `json:"long"`
}

func NewFileStorage(filename string, logger *zap.Logger) (FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return FileStorage{}, err
	}

	return FileStorage{
		file:   file,
		memory: NewMemoryStorage(nil),
		l:      logger,
	}, nil
}

func (fs FileStorage) Close() error {
	return fs.file.Close()
}

func (fs FileStorage) InitMemory() error {
	fi, err := fs.file.Stat()
	if err != nil {
		fs.l.Error("getting file info error", zap.Error(err))
		return err
	}
	if fi.Size() == 0 {
		return nil
	}

	scanner := bufio.NewScanner(fs.file)
	var url url

	for scanner.Scan() {
		data := scanner.Bytes()

		err = json.Unmarshal(data, &url)
		if err != nil {
			fs.l.Error("init memory unmarshal error", zap.Error(err))
			return err
		}

		fs.memory.storage[url.Short] = url.Long
		if url.UserID != uuid.Nil {
			fs.memory.history[url.UserID][url.Short] = url.Long
		}
	}
	return nil
}

func (fs FileStorage) Get(short string) (long string, err error) {
	return fs.memory.Get(short)
}

func (fs FileStorage) GetAll(userID uuid.UUID) (result map[string]string, err error) {
	return fs.memory.history[userID], nil
}

func (fs FileStorage) Store(userID uuid.UUID, short, long string) error {

	err := fs.memory.Store(userID, short, long)
	if err != nil {
		fs.l.Error("in-memory store error", zap.Error(err))
		return err
	}
	err = fs.writeURL(uuid.Nil, short, long)
	if err != nil {
		fs.l.Error("writing url in file error", zap.Error(err))
		return err
	}

	return nil
}

func (fs FileStorage) writeURL(userID uuid.UUID, short, long string) error {
	s := url{
		UserID: userID,
		Short:  short,
		Long:   long,
	}
	b, err := json.Marshal(s)
	if err != nil {
		fs.l.Error("writing url in file marshaling error", zap.Error(err))
		return err
	}
	b = append(b, '\n')

	_, err = fs.file.Write(b)
	if err != nil {
		fs.l.Error("writing in file error: %v", zap.Error(err))
		return err
	}
	return nil
}

func (fs FileStorage) IsKeyExist(short string) (bool, error) {
	return fs.memory.IsKeyExist(short)
}
