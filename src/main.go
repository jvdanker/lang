package main

import (
	"log"
	"net/http"

	"flag"
	"os"
)

var app_id, app_key string
var index = Index{}

func main() {
	app_id := flag.String("app-id", "", "Application ID")
	app_key := flag.String("app-key", "", "Application Key")

	flag.Parse()
	if *app_id == "" || *app_key == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	index.Reindex()

	log.Fatal(http.ListenAndServe(":8080", &WithCORS{NewRouter()}))
}
