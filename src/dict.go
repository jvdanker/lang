package main

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

type Dict struct {
	data       map[string]interface{}
	word       string
	definition string
}

type list struct {
	data []interface{}
}

func NewDictFromReader(reader io.Reader) (*Dict, error) {
	r := new(Dict)

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&r.data); err != nil {
		return nil, err
	}

	r.word = r.Get("entry.results.word")[0].(string)
	r.definition = r.Get("entry.results.lexicalEntries.entries.senses.definitions")[0].(string)

	return r, nil
}

func (d *Dict) Get(path string) []interface{} {
	parts := strings.Split(path, ".")

	result := new(list)
	get(d.data, parts, result)

	return result.data
}

func get(data interface{}, path []string, result *list) {
	for i := 0; i < len(path); i++ {
		part := path[i]
		v := reflect.TypeOf(data)

		switch v.Kind() {
		case reflect.Map:
			data = data.(map[string]interface{})[part]
			if data == nil {
				return
			}

			if reflect.TypeOf(data).Kind() == reflect.Slice {
				for _, part2 := range data.([]interface{}) {
					get(part2, path[i+1:], result)
				}

				return
			}
			break
		case reflect.Slice:
			break
		}
	}

	if len(path) == 0 || reflect.TypeOf(data).Kind() == reflect.String {
		result.data = append(result.data, data)
	}
}
