package main

import (
	"log"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/api"
	"github.com/0xc00000f/shortener-tpl/internal/logic"
	"github.com/0xc00000f/shortener-tpl/internal/storage"

	"github.com/0xc00000f/shortener-tpl/internal/handlers"
)

func main() {

	storage := storage.NewStorage()
	sa := api.NewShortenerAPI(logic.NewURLEncoder(
		logic.SetStorage(storage),
		logic.SetLength(7),
	))
	apiInstance := handlers.NewRouter(sa)

	log.Print("shortener-tpl: Enter main()")
	log.Fatal(http.ListenAndServe(":8080", apiInstance))
}
