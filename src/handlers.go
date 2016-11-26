package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"os"
	"io"
	"bytes"
	"log"
	"io/ioutil"
	"github.com/satori/go.uuid"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	game := Game{Id: uuid.NewV4().String()}
	json.NewEncoder(w).Encode(game)
}

func NextWord(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	i := getNextRandom()
	index := i[0]
	line := getLine(i[0])

	// create three options
	options := make([]string, 3)
	options[0] = line[1]
	options[1] = getLine(i[1])[1]
	options[2] = getLine(i[2])[1]

	correct := Shuffle(options)

	word := Word{
		Index:   index,
		Word:    line[0],
		Options: options,
		Correct: correct,
	}

	json.NewEncoder(w).Encode(word)
}

var url string = "https://od-api.oxforddictionaries.com:443/api/v1/entries/en/"

func AddWord(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	if !hasWord(t.Word1) {
		f, err := os.OpenFile("words.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(t.Word1 + "\t" + t.Word2 + "\n"); err != nil {
			panic(err)
		}

		reindex()
	} else {
		http.Error(w, "Dulicate", 409)
	}

	id := uuid.NewV4().String()
	fmt.Println(id)

	// fetch definition
	body1, code := lookup(url + t.Word1)
	if code != 200 {
		panic(code)
	}

	// fetch synonym
	body2, code := lookup(url + t.Word1 + "/synonyms")
	body3 := ""
	if code == 200 {
		body3 = ", \"synonyms\":" + string(body2)
	}

	// write output
	s2 := "{\"entry\":" + string(body1) + body3 + "}"

	f2, _ := os.Create("data/" + id + ".json")
	defer f2.Close()
	_, err = io.Copy(f2, bytes.NewReader([]byte(s2)))
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(s2))
}

func lookup(url string) ([]byte, int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", app_id)
	req.Header.Add("app_key", app_key)

	fmt.Println("id = ", app_id)
	return nil, 404

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if (resp.StatusCode == 200) {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		return body, resp.StatusCode
	}

	log.Fatal(resp.StatusCode)
	return nil, resp.StatusCode
}
