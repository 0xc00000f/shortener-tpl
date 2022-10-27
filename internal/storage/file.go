package storage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type FileStorage struct {
	file   *os.File
	memory MemoryStorage
}

type url struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

func NewFileStorage(filename string) (FileStorage, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return FileStorage{}, err
	}

	return FileStorage{
		file:   file,
		memory: NewMemoryStorage(),
	}, nil
}

func (fs FileStorage) Close() error {
	return fs.file.Close()
}

func (fs FileStorage) InitMemory() error {
	fi, err := fs.file.Stat()
	if err != nil {
		log.Printf("getting file info error: %v", err)
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
			log.Printf("init memory unmarshal error: %v", err)
			return err
		}

		fs.memory[url.Short] = url.Long
	}
	return nil
}

func (fs FileStorage) Get(short string) (long string, err error) {
	return fs.memory.Get(short)
}

func (fs FileStorage) Store(short, long string) error {

	err := fs.memory.Store(short, long)
	if err != nil {
		log.Printf("in-memory store error: %v", err)
		return err
	}
	err = fs.writeURL(short, long)
	if err != nil {
		log.Printf("writing url in file error: %v", err)
		return err
	}

	return nil
}

func (fs FileStorage) writeURL(short, long string) error {
	s := url{
		Short: short,
		Long:  long,
	}
	b, err := json.Marshal(s)
	if err != nil {
		log.Printf("writing url in file marshaling error: %v", err)
		return err
	}
	b = append(b, '\n')

	_, err = fs.file.Write(b)
	if err != nil {
		log.Printf("writing in file error: %v", err)
		return err
	}
	return nil
}

func (fs FileStorage) IsKeyExist(short string) (bool, error) {
	return fs.memory.IsKeyExist(short)
}
