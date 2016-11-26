package main

import (
	"fmt"
	"log"
	"net/http"

	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing/iotest"
	"time"
	"flag"
)

type Game struct {
	Id string
}

type Word struct {
	Index   int
	Word    string
	Options []string
	Correct int
}

func createIndex() {
	f, _ := os.Open("words.txt")
	defer f.Close()
	reader := bufio.NewReader(f)

	f2, _ := os.Create("words.idx")
	defer f2.Close()

	var count int16 = 0
	var index int64 = 0
	for {
		s1, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil && err != iotest.ErrTimeout {
			panic("GetLines: " + err.Error())
		}

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(index))
		f2.Write(b)

		count++
		index += int64(len(s1))
	}

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(count))
	f2.Write(b)
}

func readIndex() (int, error) {
	f, _ := os.Open("words.idx")
	defer f.Close()

	f.Seek(-2, 2)

	b := make([]byte, 2)
	f.Read(b)

	return int(binary.LittleEndian.Uint16(b)), nil
}

func getOffset(line int) int64 {
	f, _ := os.Open("words.idx")
	defer f.Close()

	f.Seek(int64(line)*int64(8), 0)

	b := make([]byte, 8)
	f.Read(b)

	return int64(binary.LittleEndian.Uint64(b))
}

func readLine(offset int64) string {
	f, _ := os.Open("words.txt")
	defer f.Close()

	f.Seek(offset, 0)

	reader := bufio.NewReader(f)
	s1, _ := reader.ReadString('\n')

	return s1
}

func hasWord(word string) bool {
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		parts := strings.Split(s, "\t")

		if strings.Compare(parts[0], word) == 0 {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return false
}

var indexSize int = 0

func reindex() {
	createIndex()
	indexSize, _ = readIndex()
}

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

type WithCORS struct {
	r *mux.Router
}

// Simple wrapper to Allow CORS.
// See: http://stackoverflow.com/a/24818638/1058612.
func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(res, req)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	game := Game{Id: uuid.NewV4().String()}
	json.NewEncoder(w).Encode(game)
}

func Shuffle(a []string) int {
	var correct int = 0

	for i := range a {
		j := rand.Intn(i + 1)
		if i == correct {
			correct = j
		}
		if j == correct {
			correct = i
		}
		a[i], a[j] = a[j], a[i]
	}

	return correct
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

func getNextRandom() []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m := make(map[int]bool)
	for len(m) < 3 {
		rand := r.Intn(indexSize)
		if _, ok := m[rand]; ok == false {
			m[rand] = true
		}

	}

	temp := make([]int, 3)
	i := 0
	for key := range m {
		temp[i] = key
		i++
	}

	return temp
}

func getLine(index int) []string {
	s := readLine(getOffset(index))
	parts := strings.Split(s, "\t")
	return parts
}

type test_struct struct {
	Word1 string
	Word2 string
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
	req.Header.Add("app_id", "1aca1eed")
	req.Header.Add("app_key", "831583fd2dedf2b02cbe6c3f3b75549d")

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

	return nil, resp.StatusCode
}
