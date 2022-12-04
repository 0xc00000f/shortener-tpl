package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"

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
		return FileStorage{}, fmt.Errorf("open file storage failed: %w", err)
	}

	return FileStorage{
		file:   file,
		memory: NewMemoryStorage(logger),
		l:      logger,
	}, nil
}

func (fs FileStorage) Close() error {
	if err := fs.file.Close(); err != nil {
		return fmt.Errorf("closing file storage failed: %w", err)
	}

	return nil
}

func (fs FileStorage) InitMemory() error {
	fi, err := fs.file.Stat()
	if err != nil {
		return fmt.Errorf("getting file info error: %w", err)
	}

	if fi.Size() == 0 {
		return nil
	}

	var url url

	scanner := bufio.NewScanner(fs.file)
	for scanner.Scan() {
		data := scanner.Bytes()

		err = json.Unmarshal(data, &url)
		if err != nil {
			return fmt.Errorf("init memory unmarshal error: %w", err)
		}

		fs.memory.storage[url.Short] = url.Long

		if url.UserID != uuid.Nil {
			if _, ok := fs.memory.history[url.UserID]; !ok {
				fs.memory.history[url.UserID] = map[string]string{}
			}

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

	err = fs.writeURL(userID, short, long)
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
		return fmt.Errorf("writing url in file marshaling error: %w", err)
	}

	b = append(b, '\n')

	_, err = fs.file.Write(b)
	if err != nil {
		return fmt.Errorf("writing result in file error: %w", err)
	}

	return nil
}

func (fs FileStorage) IsKeyExist(short string) (bool, error) {
	return fs.memory.IsKeyExist(short)
}
