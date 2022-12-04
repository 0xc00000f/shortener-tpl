package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	// path to the file with shortened URLs
	Filepath string `env:"FILE_STORAGE_PATH"`

	// address of the HTTP server
	Address string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1:8080"`

	// base URL of the resulting shortened URL
	BaseURL string `env:"BASE_URL" envDefault:"http://127.0.0.1:8080"`

	// address of the database
	DatabaseAddress string `env:"DATABASE_DSN"`
}

func New() (Cfg, error) {
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
	flag.Parse()

	return cfg, nil
}
