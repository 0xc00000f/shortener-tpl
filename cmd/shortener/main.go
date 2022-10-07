package main

import (
	"github.com/0xc00000f/shortener-tpl/internal/app"
	"log"
)

func main() {

	log.Print("shortener-tpl: Enter main()")
	log.Fatal(app.Server().ListenAndServe())

}
