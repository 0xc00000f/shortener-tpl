package main

import (
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"go.uber.org/zap"
)

var (
	fspF     *string
	addressF *string
	baseURL  *string

	systemAddressKey         = "SERVER_ADDRESS"
	defaultAddress           = ":8080"
	systemFileStoragePathKey = "FILE_STORAGE_PATH"
)

type cfg struct {
	l *zap.Logger
}

func init() {
	fspF = flag.String("f", "", "responsible for the path to the file with shortened URLs")
	addressF = flag.String("a", "", "responsible for the start address of the HTTP server")
	baseURL = flag.String(
		"b",
		"",
		"responsible for the base address of the resulting shortened URL")
}

func main() {
	l, _ := zap.NewProduction()
	defer l.Sync()

	cfg := cfg{l}

	storage, err := cfg.chooseStorage()
	if err != nil {
		l.Error("choosing storage error: %v", zap.Error(err))
		return
	}

	address := cfg.chooseAddress()

	logic := logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
		logic.SetLogger(l),
	)

	sa := api.NewShortenerAPI(
		api.SetLogic(logic),
		api.InitBaseURL(*baseURL),
		api.SetLogger(l),
	)

	apiInstance := handlers.NewRouter(sa)

	l.Info("starting server", zap.String("address", address))
	l.Fatal("server fatal error", zap.Error(http.ListenAndServe(address, apiInstance)))
}

func (cfg cfg) chooseStorage() (logic logic.URLStorager, err error) {
	if *fspF != "" {
		cfg.l.Info("choose storage from flag", zap.String("fspF", *fspF))
		return cfg.creatingFileStorage(*fspF)
	}

	fsp, ok := os.LookupEnv(systemFileStoragePathKey)
	switch ok {
	case true:
		cfg.l.Info("choose storage from environment variable", zap.String("fsp", fsp))
		return cfg.creatingFileStorage(fsp)
	case false:
		cfg.l.Info("choose in-memory storage")
		storage := storage.NewMemoryStorage()
		return storage, nil
	default:
		return nil, errors.New("unknown storage")
	}
}

func (cfg cfg) creatingFileStorage(path string) (logic logic.URLStorager, err error) {
	storage, err := storage.NewFileStorage(path, cfg.l)
	if err != nil {
		cfg.l.Error("creating file storage err", zap.Error(err))
		return nil, err
	}

	err = storage.InitMemory()
	if err != nil {
		cfg.l.Error("init file storage memory err", zap.Error(err))
		return nil, err
	}
	return storage, nil
}

func (cfg cfg) chooseAddress() (address string) {
	if *addressF != "" {
		address = *addressF
	} else {
		var ok bool
		address, ok = os.LookupEnv(systemAddressKey)
		if !ok {
			address = defaultAddress
		}
	}

	return
}
