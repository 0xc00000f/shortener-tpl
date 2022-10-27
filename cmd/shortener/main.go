package main

import (
	"errors"
	"flag"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
)

var (
	fspF     *string
	addressF *string
	baseURL  *string

	systemAddressKey         = "SERVER_ADDRESS"
	defaultAddress           = ":8080"
	systemFileStoragePathKey = "FILE_STORAGE_PATH"
)

func init() {
	fspF = flag.String("f", "", "responsible for the path to the file with shortened URLs")
	addressF = flag.String("a", "", "responsible for the start address of the HTTP server")
	baseURL = flag.String(
		"b",
		"",
		"responsible for the base address of the resulting shortened URL")
}

func main() {
	flag.Parse()

	storage, err := chooseStorage()
	if err != nil {
		log.Printf("choosing storage error: %v", err)
		return
	}

	logic := logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	)

	sa := api.NewShortenerAPI(
		api.SetLogic(logic),
		api.InitBaseURL(*baseURL),
	)

	apiInstance := handlers.NewRouter(sa)
	address := chooseAddress()

	log.Printf("starting server - address: %s, handler: %v", address, apiInstance)
	log.Fatal(http.ListenAndServe(address, apiInstance))
}

func chooseStorage() (logic logic.URLStorager, err error) {
	if *fspF != "" {
		log.Printf("choose storage from flag: %s", *fspF)
		return creatingFileStorage(*fspF)
	}

	fsp, ok := os.LookupEnv(systemFileStoragePathKey)
	switch ok {
	case true:
		log.Printf("choose storage from environment variable: %s", fsp)
		return creatingFileStorage(fsp)
	case false:
		log.Printf("choose in-memory storage")
		storage := storage.NewMemoryStorage()
		return storage, nil
	default:
		return nil, errors.New("unknown storage")
	}
}

func creatingFileStorage(path string) (logic logic.URLStorager, err error) {
	storage, err := storage.NewFileStorage(path)
	if err != nil {
		log.Printf("creating file storage err: %s", err)
		return nil, err
	}

	err = storage.InitMemory()
	if err != nil {
		log.Printf("init file storage memory err: %s", err)
		return nil, err
	}
	return storage, nil
}

func chooseAddress() (address string) {
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
