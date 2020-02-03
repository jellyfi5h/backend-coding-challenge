package server

import (
	"encoding/json"
	"net/http"
	"time"

	"../lib/web"
)

func createdSince(since string) time.Time {
	var created time.Time

	now := time.Now()
	switch since {
	case "daily":
		created = now.AddDate(0, 0, -1)
	case "weekly":
		created = now.AddDate(0, 0, -7)
	case "monthly":
		created = now.AddDate(0, -1, 0)
	default: //monthly
		created = now.AddDate(0, -1, 0)
	}
	return created
}

func setHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

// /languages?since={daily|weekly|monthly(default)}
func languagesHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte

	status := web.AllowedMethods(r.Method, "GET")
	title := "error"
	if status == 200 {
		query, err := web.GetQueryValues(r.URL, "since")
		if err != nil {
			resp = []byte(err.Error()) //syntax error in request
			status = 400
		} else {
			languages, err := TrendingLanguages(query["since"])
			if err != nil {
				resp = []byte(err.Error())
			} else {
				title = "languages"
				resp, _ = json.Marshal(languages)
				//clear frameworks
				go Popularframeworks(languages) // run in separate goroutine(thread) cause its take too much time due its complicated process
			}
		}
	}
	setHeader(w, status)
	data := web.JsonResponse(status, title, resp)
	w.Write(data)
}

// /frameworks
func frameworksHandler(w http.ResponseWriter, r *http.Request) {
	var resp []byte

	status := web.AllowedMethods(r.Method, "GET")
	title := "error"
	if status == 200 {
		if G_frameworks == nil { //framework list is empty will get the trending repos first
			languages, err := TrendingLanguages("monthly")
			Popularframeworks(languages)
			if err != nil {
				resp = []byte(err.Error())
			}
		} else {
			title = "frameworks"
			resp, _ = json.Marshal(G_frameworks)
		}
	}
	setHeader(w, status)
	data := web.JsonResponse(status, title, resp)
	w.Write(data)
}
