package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/google/uuid"

	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/models"
)

type FileStorage struct {
	file   *os.File
	memory MemoryStorage

	mu sync.RWMutex
	l  *zap.Logger
}

func NewFileStorage(filename string, logger *zap.Logger) (*FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &FileStorage{
		file:   file,
		memory: *NewMemoryStorage(logger),
		l:      logger,
	}, nil
}

func (fs *FileStorage) Close() error {
	return fs.file.Close()
}

func (fs *FileStorage) InitMemory() error {
	fi, err := fs.file.Stat()
	if err != nil {
		fs.l.Error("getting file info error", zap.Error(err))
		return err
	}

	if fi.Size() == 0 {
		return nil
	}

	var url models.URL

	scanner := bufio.NewScanner(fs.file)
	for scanner.Scan() {
		data := scanner.Bytes()

		err = json.Unmarshal(data, &url)
		if err != nil {
			fs.l.Error("init memory unmarshal error", zap.Error(err))
			return err
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

func (fs *FileStorage) Get(ctx context.Context, short string) (long string, err error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.memory.Get(ctx, short)
}

//revive:disable-next-line
func (fs *FileStorage) GetAll(ctx context.Context, userID uuid.UUID) (result map[string]string, err error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.memory.history[userID], nil
}

func (fs *FileStorage) Store(ctx context.Context, userID uuid.UUID, short, long string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	err := fs.memory.Store(ctx, userID, short, long)
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

func (fs *FileStorage) writeURL(userID uuid.UUID, short, long string) error {
	s := models.URL{
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

func (fs *FileStorage) IsKeyExist(ctx context.Context, short string) (bool, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return fs.memory.IsKeyExist(ctx, short)
}
