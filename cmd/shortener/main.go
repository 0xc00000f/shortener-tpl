package main

import (
	"errors"
	"github.com/0xc00000f/shortener-tpl/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/handlers"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
)

func main() {
	log.Print("shortener-tpl: Enter main()")

	storage, err := chooseStorage("FILE_STORAGE_PATH")
	if err != nil {
		log.Printf("main::chooseStorage -- error: %s", err)
		return
	}

	sa := api.NewShortenerAPI(logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	))
	apiInstance := handlers.NewRouter(sa)

	address, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok {
		address = ":8080"
	}
	log.Fatal(http.ListenAndServe(address, apiInstance))
}

func chooseStorage(storagePath string) (logic.URLStorager, error) {
	fsp, ok := os.LookupEnv(storagePath)
	switch ok {
	case true:
		storage, err := storage.NewFileStorage(fsp)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		err = storage.InitMemory()
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		return storage, nil
	case false:
		storage := storage.NewMemoryStorage()
		return storage, nil
	default:
		return nil, errors.New("unknown storage")
	}
}