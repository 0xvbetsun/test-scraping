package main

import (
	"log"
	"net/http"

	"github.com/vbetsun/scraping"
)

func main() {
	handler := scraping.NewHandler(0)
	http.Handle("/", handler)
	log.Print("Listening...")
	http.ListenAndServe(":3000", nil)
}
