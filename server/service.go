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
	default:
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
			resp = []byte(err.Error())
			status = 400
		} else {
			createdTime := createdSince(query["since"])
			url := trendingURL(createdTime)
			languages, err := getLanguages(url)
			if err != nil {
				resp = []byte(err.Error())
			} else {
				title = "languages"
				resp, _ = json.Marshal(languages)
			}
		}
	}
	setHeader(w, status)
	data := web.JsonResponse(status, title, resp)
	w.Write(data)
}

// /repositories/:language?since={...}&lang
func reposByLangHandler(w http.ResponseWriter, r *http.Request) {

}
