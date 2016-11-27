package main

import (
	"log"
	"net/http"

	"flag"
	"github.com/gorilla/mux"
	"os"
)

var app_id, app_key string

func main() {
	app_id := flag.String("app-id", "", "Application ID")
	app_key := flag.String("app-key", "", "Application Key")

	flag.Parse()
	if *app_id == "" || *app_key == "" {
		flag.PrintDefaults()
		os.Exit(-1)
	}

	reindex()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/v1/game/new", TodoIndex)
	router.HandleFunc("/v1/game/word", NextWord)
	router.HandleFunc("/v1/game/word/add", AddWord).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", &WithCORS{router}))
}
