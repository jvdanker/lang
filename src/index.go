package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Game struct {
	Id string
}

type Word struct {
	Index   int `json:"index"`
	Word    string
	Options []string
	Correct int
}

type test_struct struct {
	Word1 string `json:"word1"`
	Word2 string
}

type Index struct {
	Test string
}

var indexSize int = 0
var files []string

func (i *Index) Reindex() {
	dirname := "data" + string(filepath.Separator)
	d, err := os.Open(dirname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	fi, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	indexSize = len(fi)
	files = nil
	for _, fi := range fi {
		if fi.Mode().IsRegular() {
			files = append(files, fi.Name())
		}
	}
}

func (i *Index) GetNextRandom() []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m := make(map[int]bool)
	for len(m) < 3 {
		rand := r.Intn(indexSize - 1)
		if _, ok := m[rand]; ok == false {
			m[rand] = true
		}

	}

	temp := make([]int, 3)
	x := 0
	for key := range m {
		temp[x] = key
		x++
	}

	return temp
}

func (i *Index) Shuffle(a []string) int {
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

func (i *Index) GetLine(index int) []string {
	f, _ := os.Open("data/" + files[index])

	d, err := NewDictFromReader(f)
	if err != nil {
		panic(err)
	}

	s := make([]string, 2)
	s[0] = d.word
	s[1] = d.definition

	return s
}

func (i *Index) HasWord(word string) bool {
	return false
}
