package server

import "net/http"

// /languages?since={daily|weekly|monthly(default)}
func languagesHandler(w http.ResponseWriter, r *http.Request) {

}

// /repositories/:language?since={...}&lang
func reposByLangHandler(w http.ResponseWriter, r *http.Request) {

}
