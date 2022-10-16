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
		log.Print("FileStorage::InitMemory -- file.Stat error")
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
			log.Print("FileStorage::InitMemory -- json.Unmarshal error")
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
		log.Print("FileStorage::Store -- memory.Store error")
		return err
	}
	err = fs.writeURL(short, long)
	if err != nil {
		log.Print("FileStorage::Store -- writeURL error")
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
		log.Print("FileStorage::writeURL -- Marshal error")
		return err
	}
	b = append(b, '\n')

	_, err = fs.file.Write(b)
	if err != nil {
		log.Print("FileStorage::writeURL -- file.Write error")
		return err
	}
	return nil
}

func (fs FileStorage) IsKeyExist(short string) (bool, error) {

	fi, err := fs.file.Stat()
	if err != nil {
		log.Print("FileStorage::IsKeyExist -- file.Stat error")
		return false, err
	}
	if fi.Size() == 0 {
		return false, nil
	}

	scanner := bufio.NewScanner(fs.file)
	var url url

	for scanner.Scan() {

		data := scanner.Bytes()

		err = json.Unmarshal(data, &url)
		if err != nil {
			return false, err
		}

		if url.Short == short {
			return true, nil
		}
	}

	return false, nil
}
