package config

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
)

const (
	FileStorageKey     = "FILE_STORAGE_PATH" // file storage path key -- environment variable
	SystemAddressKey   = "SERVER_ADDRESS"    // address key -- environment variable
	SystemBaseURLKey   = "BASE_URL"          // base url key -- environment variable
	SystemDBAddressKey = "DATABASE_DSN"      // database address key -- environment variable
	DefaultAddress     = "127.0.0.1:8080"
	DefaultAddressDB   = "127.0.0.1:8080"
)

type Cfg struct {
	Filepath        string // path to the file with shortened URLs
	Address         string // address of the HTTP server
	BaseURL         string // base URL of the resulting shortened URL
	DatabaseAddress string // address of the database

	L *zap.Logger // logger
}

func New(logger *zap.Logger) (Cfg, error) {
	cfg := Cfg{L: logger}

	flag.StringVar(&cfg.Filepath, "f", "", "responsible for the path to the file with shortened URLs")
	flag.StringVar(&cfg.Address, "a", "", "responsible for the start Address of the HTTP server")
	flag.StringVar(&cfg.DatabaseAddress, "d", "", "responsible for the database address")
	flag.StringVar(&cfg.BaseURL,
		"b",
		"",
		"responsible for the base Address of the resulting shortened URL")
	flag.Parse()

	cfg.chooseFilepath()
	cfg.chooseAddress()
	cfg.chooseBaseURL()

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

func (cfg *Cfg) chooseBaseURL() {
	// if is set by flags
	if cfg.BaseURL != "" {
		cfg.L.Info("choose base url from flag", zap.String("BaseURL", cfg.BaseURL))
		return
	}

	var ok bool
	// try to set value from system environment variable
	bu, ok := os.LookupEnv(SystemBaseURLKey)
	if !ok {
		// set default value
		cfg.L.Info("choose default base url")
		bu = fmt.Sprintf("http://%s", DefaultAddress)
	}

	cfg.BaseURL = bu
	cfg.L.Info("base url chosen", zap.String("BaseURL", cfg.BaseURL))
}

func (cfg *Cfg) chooseDatabaseAddress() {
	// if is set by flags
	if cfg.DatabaseAddress != "" {
		cfg.L.Info("choose database address from flag", zap.String("DatabaseAddress", cfg.DatabaseAddress))
		return
	}

	// try to set value from system environment variable
	cfg.DatabaseAddress = os.Getenv(SystemDBAddressKey)
	cfg.L.Info("database address chosen", zap.String("Database Address", cfg.DatabaseAddress))
}
