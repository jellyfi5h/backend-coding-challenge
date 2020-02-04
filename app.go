package main

import (
	"fmt"
	"net/http"

	"./api"
)

func main() {
	route := entryPoints()
	err := http.ListenAndServe(":8100", route)
	if err != nil {
		fmt.Println(err)
	}
}

func entryPoints() *http.ServeMux {
	route := http.NewServeMux()

	// help manual
	route.HandleFunc("/", api.HelpHandler)
	// GET /languages?since={Daily | weekly | Monthly(default)}
	route.HandleFunc("/languages", api.LanguagesHandler())
	// GET /frameworks
	route.HandleFunc("/frameworks", api.FrameworksHandler())
	return route
}
