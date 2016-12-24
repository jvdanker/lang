package main

import (
	"bytes"
	"encoding/json"

	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var url string = "https://od-api.oxforddictionaries.com/api/v1/entries/en/"

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	game := Game{Id: uuid.NewV4().String()}
	if err := json.NewEncoder(w).Encode(game); err != nil {
		panic(err)
	}
}

func NextWord(w http.ResponseWriter, r *http.Request) {
	i := index.GetNextRandom()
	idx := i[0]
	line := index.GetLine(i[0])

	// create three options
	options := make([]string, 3)
	options[0] = line[1]
	options[1] = index.GetLine(i[1])[1]
	options[2] = index.GetLine(i[2])[1]

	correct := index.Shuffle(options)

	word := Word{
		Index:   idx,
		Word:    line[0],
		Options: options,
		Correct: correct,
	}

	json.NewEncoder(w).Encode(word)
}

func AddWord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if index.HasWord(t.Word1) {
		http.Error(w, "Conflict", http.StatusConflict)
		return
	}

	// fetch definition
	body1, code := lookup(url + t.Word1)
	if code != 200 {
		log.Println(code)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// fetch synonym
	body2, code := lookup(url + t.Word1 + "/synonyms")
	if code != 200 {
		log.Println(code)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// write output
	s2 := "{\"entry\":" + string(body1) + ", \"synonyms\":" + string(body2) + "}"

	id := uuid.NewV4().String()
	f2, err := os.Create("data/" + id + ".json")
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	defer f2.Close()
	_, err = io.Copy(f2, bytes.NewReader([]byte(s2)))
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(s2))
	index.Reindex()
}

func lookup(url string) ([]byte, int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", *app_id)
	req.Header.Add("app_key", *app_key)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
		return nil, resp.StatusCode
	}

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		return body, resp.StatusCode
	}

	log.Println(resp.StatusCode)
	return nil, resp.StatusCode
}
