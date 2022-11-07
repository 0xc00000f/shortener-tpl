package config

import (
	"flag"
	"os"

	"go.uber.org/zap"
)

const (
	FileStorageKey   = "FILE_STORAGE_PATH" // file storage path key -- environment variable
	SystemAddressKey = "SERVER_ADDRESS"    // address key -- environment variable
	DefaultAddress   = ":8080"
)

type Cfg struct {
	Filepath string // path to the file with shortened URLs
	Address  string // address of the HTTP server
	BaseURL  string // base URL of the resulting shortened URL

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

	cfg.chooseFilepath()
	cfg.chooseAddress()

	return cfg, nil
}

func (cfg *Cfg) chooseFilepath() {
	// if filepath is set by flags
	if cfg.Filepath != "" {
		cfg.L.Info("choose filepath from flag", zap.String("Filepath", cfg.Filepath))
		return
	}

	// try to set filepath from system environment variable
	filepath, ok := os.LookupEnv(FileStorageKey)
	if !ok {
		cfg.L.Info("Filepath is empty", zap.String("Filepath", cfg.Filepath))
		return
	}

	// filepath is set by system environment variable, create file storage
	cfg.Filepath = filepath
	cfg.L.Info("Filepath found in environment variable",
		zap.String("environment variable", FileStorageKey), zap.String("Filepath", cfg.Filepath))
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
