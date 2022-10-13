package main

import (
	"github.com/0xc00000f/shortener-tpl/internal/app"
	"log"
	"net/http"
)

func main() {

	log.Print("shortener-tpl: Enter main()")
	log.Fatal(http.ListenAndServe(":8080", app.NewRouter()))

}
