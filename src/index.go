package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing/iotest"
	"time"
)

var indexSize int = 0

type Game struct {
	Id string
}

type Word struct {
	Index   int
	Word    string
	Options []string
	Correct int
}

type test_struct struct {
	Word1 string
	Word2 string
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

func reindex() {
	createIndex()
	indexSize, _ = readIndex()
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
