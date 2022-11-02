package config

import (
	"flag"
	"os"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"go.uber.org/zap"
)

const (
	FileStorageKey   = "FILE_STORAGE_PATH" // file storage path key -- environment variable
	SystemAddressKey = "SERVER_ADDRESS"    // address key -- environment variable
	DefaultAddress   = ":8080"
)

type Cfg struct {
	Filepath string              // path to the file with shortened URLs
	Address  string              // address of the HTTP server
	BaseURL  string              // base URL of the resulting shortened URL
	Storage  encoder.URLStorager // storage instance shortened URLs: in-memory / file

	L *zap.Logger // logger
}

func New(logger *zap.Logger) (Cfg, error) {
	cfg := Cfg{L: logger}

	flag.StringVar(&cfg.Filepath, "f", "", "responsible for the path to the file with shortened URLs")
	flag.StringVar(&cfg.Address, "a", "", "responsible for the start Address of the HTTP server")
	flag.StringVar(&cfg.BaseURL,
		"b",
		"",
		"responsible for the base Address of the resulting shortened URL")
	flag.Parse()

	err := cfg.chooseStorage()
	if err != nil {
		cfg.L.Error("choose storage err", zap.Error(err))
		return cfg, err
	}
	cfg.chooseAddress()

	return cfg, nil
}

	// if filepath is set by flags create file storage
func (cfg *Cfg) chooseStorage() (err error) {
	if cfg.Filepath != "" {
		cfg.L.Info("choose storage from flag", zap.String("Filepath", cfg.Filepath))
		return cfg.creatingFileStorage(cfg.Filepath)
	}

	// try to set filepath from system environment variable
	filepath, ok := os.LookupEnv(FileStorageKey)
	if !ok {
		// create in-memory storage
		cfg.L.Info("choose in-memory storage")
		cfg.Storage = storage.NewMemoryStorage(cfg.L)
		return nil
	}

	// filepath is set by system environment variable, create file storage
	cfg.Filepath = filepath
	cfg.L.Info("choose storage from environment variable", zap.String("Filepath", filepath))
	return cfg.creatingFileStorage(filepath)
}

func (cfg *Cfg) creatingFileStorage(path string) (err error) {
	storage, err := storage.NewFileStorage(path, cfg.L)
	if err != nil {
		cfg.L.Error("creating file storage err", zap.Error(err))
		return err
	}

	err = storage.InitMemory()
	if err != nil {
		cfg.L.Error("init file storage memory err", zap.Error(err))
		return err
	}

	cfg.Storage = storage
	return nil
}

func (cfg *Cfg) chooseAddress() {
	// if is set by flags
	if cfg.Address != "" {
		return
	}

	var ok bool
	// try to set value from system environment variable
	address, ok := os.LookupEnv(SystemAddressKey)
	if !ok {
		// set default value
		address = DefaultAddress
	}

	cfg.Address = address
}
