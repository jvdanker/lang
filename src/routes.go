package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"TodoIndex",
		"GET",
		"/v1/game/new",
		TodoIndex,
	},
	Route{
		"NextWord",
		"GET",
		"/v1/game/word",
		NextWord,
	},
	Route{
		"AddWord",
		"POST",
		"/v1/game/word/add",
		AddWord,
	},
}
