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
	log.Print("shortener-tpl: main -- entered")
	log.Print("shortener-tpl: main -- parse")
	flag.Parse()

	log.Print("shortener-tpl: main -- chooseStorage")
	storage, err := chooseStorage()
	if err != nil {
		log.Printf("main::chooseStorage -- error: %s", err)
		return
	}

	log.Print("shortener-tpl: main -- NewURLEncoder")
	logic := logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	)

	log.Print("shortener-tpl: main -- NewShortenerAPI")
	sa := api.NewShortenerAPI(
		api.SetLogic(logic),
		api.InitBaseURL(*baseURL),
	)

	log.Print("shortener-tpl: main -- creating NewRouter")
	apiInstance := handlers.NewRouter(sa)

	log.Print("shortener-tpl: main -- creating chooseAddress")
	address := chooseAddress()

	log.Printf("shortener-tpl: main -- starting server NewRouter, address: %s, handler: %v", address, apiInstance)
	log.Fatal(http.ListenAndServe(address, apiInstance))
}

func chooseStorage() (logic logic.URLStorager, err error) {
	log.Print("main::chooseStorage -- entered")
	defer log.Printf("main::chooseStorage -- finished, returned storager: %v, error: %s", logic, err)

	if *fspF != "" {
		log.Printf("main::chooseStorage -- choose storage from flag: %s", *fspF)
		return creatingFileStorage(*fspF)
	}

	log.Printf("main::chooseStorage -- trying choose storage from system variable: %s", systemFileStoragePathKey)
	fsp, ok := os.LookupEnv(systemFileStoragePathKey)
	switch ok {
	case true:
		log.Printf("main::chooseStorage -- chose storage from system variable. "+
			"System variable - %s, path - %s", systemFileStoragePathKey, fsp)
		return creatingFileStorage(fsp)
	case false:
		log.Printf("main::chooseStorage -- can't choose storage from system variable. "+
			"System variable - %s, path - %s", systemFileStoragePathKey, fsp)
		storage := storage.NewMemoryStorage()
		log.Print("main::chooseStorage -- created in-memory storage")
		return storage, nil
	default:
		return nil, errors.New("unknown storage")
	}

}

func creatingFileStorage(path string) (logic logic.URLStorager, err error) {
	log.Printf("main::creatingFileStorage -- entered, arguments - path: %s", path)
	defer log.Printf("main::creatingFileStorage -- finished, returned storager: %v, error: %v", logic, err)

	log.Printf("main::creatingFileStorage -- creating new file storage - path: %s", path)
	storage, err := storage.NewFileStorage(path)
	if err != nil {
		log.Printf("main::creatingFileStorage -- storage.NewFileStorage -- err != nil - err: %s", err)
		log.Fatal(err)
		return nil, err
	}
	log.Print("main::creatingFileStorage -- initing memory")
	err = storage.InitMemory()
	if err != nil {
		log.Printf("main::creatingFileStorage -- storage.InitMemory -- err != nil - err: %s", err)
		log.Fatal(err)
		return nil, err
	}
	return storage, nil
}

func chooseAddress() (address string) {
	log.Print("main::chooseAddress -- entered")
	if *addressF != "" {
		address = *addressF
		log.Printf("main::chooseAddress -- chose address from flag: %s", address)
	} else {
		var ok bool
		log.Printf("main::chooseAddress -- trying get address from system variable: %s", systemAddressKey)
		address, ok = os.LookupEnv(systemAddressKey)
		if ok {
			log.Printf("main::chooseAddress -- got address from system variable. "+
				"System variable - %s, address - %s", systemAddressKey, address)
		} else {
			log.Printf("main::chooseAddress -- can't get address from system variable. "+
				"System variable - %s, address - %s", systemAddressKey, address)
			address = defaultAddress
			log.Printf("main::chooseAddress -- taking default address - %s", address)
		}
	}

	log.Printf("main::chooseAddress -- finished, returned address: %s", address)
	return
}
