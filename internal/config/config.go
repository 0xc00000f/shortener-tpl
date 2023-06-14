package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Cfg struct {
	// path to the file with shortened URLs
	Filepath string `env:"FILE_STORAGE_PATH"`

	// address of the HTTP server
	Address string `env:"SERVER_ADDRESS"`

	// base URL of the resulting shortened URL
	BaseURL string `env:"BASE_URL"`

	// address of the database
	DatabaseAddress string `env:"DATABASE_DSN"`

	// is TLS enabled
	TLSEnabled bool `env:"ENABLE_HTTPS"`

	// tls certificate file
	TLSCertFile string `env:"TLS_CERT_FILE"`

	// tls key file
	TLSKeyFile string `env:"TLS_KEY_FILE"`

	// json config file
	JSONConfig string `env:"CONFIG"`

	// trusted subnet
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}

func New(logger *zap.Logger) (Cfg, error) {
	cfg := Cfg{}
	if err := env.Parse(&cfg); err != nil {
		panic("can't parse config")
	}

	flag.StringVar(
		&cfg.Filepath,
		"f",
		cfg.Filepath,
		"responsible for the path to the file with shortened URLs",
	)
	flag.StringVar(
		&cfg.Address,
		"a",
		cfg.Address,
		"responsible for the start Address of the HTTP server",
	)
	flag.StringVar(
		&cfg.DatabaseAddress,
		"d",
		cfg.DatabaseAddress,
		"responsible for the database address",
	)
	flag.StringVar(&cfg.BaseURL,
		"b",
		cfg.BaseURL,
		"responsible for the base Address of the resulting shortened URL")
	flag.BoolVar(&cfg.TLSEnabled,
		"s",
		cfg.TLSEnabled,
		"responsible for the TLS enabled")
	flag.StringVar(&cfg.JSONConfig,
		"c",
		cfg.JSONConfig,
		"responsible for the json config file")
	flag.StringVar(&cfg.JSONConfig,
		"config",
		cfg.JSONConfig,
		"responsible for the json config file",
	)
	flag.StringVar(&cfg.TrustedSubnet,
		"t",
		cfg.TrustedSubnet,
		"responsible for trusted subnet address")
	flag.Parse()

	if err := cfg.parseJSONConfig(cfg.JSONConfig); err != nil {
		logger.Error("can't parse json config", zap.Error(err))
	}

	return cfg, nil
}

func (c *Cfg) parseJSONConfig(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Cfg
	if err := json.Unmarshal(b, &cfg); err != nil {
		return err
	}

	if cfg.Filepath == "" {
		c.Filepath = cfg.Filepath
	}

	if cfg.Address == "" {
		c.Address = cfg.Address
	}

	if cfg.DatabaseAddress == "" {
		c.DatabaseAddress = cfg.DatabaseAddress
	}

	if cfg.BaseURL == "" {
		c.BaseURL = cfg.BaseURL
	}

	if !cfg.TLSEnabled {
		c.TLSEnabled = cfg.TLSEnabled
	}

	if cfg.TLSCertFile == "" {
		c.TLSCertFile = cfg.TLSCertFile
	}

	if cfg.TLSKeyFile == "" {
		c.TLSKeyFile = cfg.TLSKeyFile
	}

	if cfg.TrustedSubnet == "" {
		c.TrustedSubnet = cfg.TrustedSubnet
	}

	return nil
}
