package config

import (
	"flag"
	"os"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"go.uber.org/zap"
)

type cfg struct {
	filepath string              // path to the file with shortened URLs
	Address  string              // address of the HTTP server
	BaseURL  string              // base URL of the resulting shortened URL
	Storage  encoder.URLStorager // storage instance shortened URLs: in-memory / file

	l *zap.Logger // logger
}

func New(logger *zap.Logger) (cfg, error) {
	cfg := cfg{l: logger}

	flag.StringVar(&cfg.filepath, "f", "", "responsible for the path to the file with shortened URLs")
	flag.StringVar(&cfg.Address, "a", "", "responsible for the start Address of the HTTP server")
	flag.StringVar(&cfg.BaseURL,
		"b",
		"",
		"responsible for the base Address of the resulting shortened URL")
	flag.Parse()

	err := cfg.chooseStorage()
	if err != nil {
		cfg.l.Error("choose storage err", zap.Error(err))
		return cfg, err
	}
	cfg.chooseAddress()

	return cfg, nil
}

func (cfg *cfg) chooseStorage() (err error) {
	const fileStorageKey = "FILE_STORAGE_PATH" // file storage path key -- environment variable

	// if filepath is set by flags create file storage
	if cfg.filepath != "" {
		cfg.l.Info("choose storage from flag", zap.String("filepath", cfg.filepath))
		return cfg.creatingFileStorage(cfg.filepath)
	}

	// try to set filepath from system environment variable
	filepath, ok := os.LookupEnv(fileStorageKey)
	if !ok {
		// create in-memory storage
		cfg.l.Info("choose in-memory storage")
		cfg.Storage = storage.NewMemoryStorage(cfg.l)
		return nil
	}

	// filepath is set by system environment variable, create file storage
	cfg.filepath = filepath
	cfg.l.Info("choose storage from environment variable", zap.String("filepath", filepath))
	return cfg.creatingFileStorage(filepath)
}

func (cfg *cfg) creatingFileStorage(path string) (err error) {
	storage, err := storage.NewFileStorage(path, cfg.l)
	if err != nil {
		cfg.l.Error("creating file storage err", zap.Error(err))
		return err
	}

	err = storage.InitMemory()
	if err != nil {
		cfg.l.Error("init file storage memory err", zap.Error(err))
		return err
	}

	cfg.Storage = storage
	return nil
}

func (cfg *cfg) chooseAddress() {
	const systemAddressKey = "SERVER_ADDRESS" // address key -- environment variable
	const defaultAddress = ":8080"

	// if is set by flags
	if cfg.Address != "" {
		return
	}

	var ok bool
	// try to set value from system environment variable
	address, ok := os.LookupEnv(systemAddressKey)
	if !ok {
		// set default value
		address = defaultAddress
	}

	cfg.Address = address
}
