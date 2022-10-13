package main

import (
	"log"
	"net/http"

	"github.com/0xc00000f/shortener-tpl/internal/app"
)

func main() {
	log.Print("shortener-tpl: Enter main()")
	log.Fatal(http.ListenAndServe(":8080", app.NewRouter()))
}
