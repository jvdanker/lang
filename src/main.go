package main

import (
	"log"
	"net/http"

	"flag"
	"github.com/gorilla/handlers"
	"os"
	"time"
)

var app_id, app_key *string
var data_dir, web_dir *string
var index = Index{}

func main() {
	app_id = flag.String("app-id", "", "Application ID")
	app_key = flag.String("app-key", "", "Application Key")
	data_dir = flag.String("data", "", "Data directory")
	web_dir = flag.String("web", "", "Web directory")

	flag.Parse()
	if *app_id == "" || *app_key == "" || *web_dir == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	index.Reindex()

	log.Println("Started server and listening at port 8080...")

	srv := &http.Server{
		Handler:      handlers.CompressHandler(&WithCORS{NewRouter()}),
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
