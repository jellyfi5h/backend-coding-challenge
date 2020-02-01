package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//set timeout of 15 sec in case the remote server is unresponsive
var httpClient = &http.Client{Timeout: 15 * time.Second}

type Repository struct {
	Name      string `json:"name"`
	Developer string `json:"developer_login"`
	RepoURL   string `json:"repo_url"`
	Stars     int    `json:"stars"`
	Watchers  int    `json:"watchers"`
	Forks     int    `json:"forks"`
}

type Language struct {
	CountRepos   int           `json:"repos_count"`
	Repositories []*Repository `json:"repositories"`
}

//get 100 trending repos of all languages
func trendingURL(dtrange time.Time) string {
	return fmt.Sprintf("https://api.github.com/search/repositories?q="+
		"created:>=%s&sort=stars&order=desc&per_page=100", dtrange.Format("2006-01-02"))
}

//get 100 trending repos by specific language
func trendingByLangURL(dtrange time.Time, lang string) string {
	return fmt.Sprintf("https://api.github.com/search/repositories?q="+
		"created:>=%s&language=%s&sort=stars&order=desc&per_page=100", dtrange.Format("2006-01-02"), lang)
}

//get trending repositories by languages
func getLanguages(url string) (languages map[string]*Language, err error) {
	var repos []interface{}

	repos, err = trendingRepos(url)
	if err != nil {
		return
	}
	languages = reposByLang(repos)
	return
}

func unmarshalItem(data map[string]interface{}, field string) (item interface{}, err error) {
	var dump []byte

	dump, err = json.Marshal(data[field])
	if err != nil {
		return
	}
	err = json.Unmarshal(dump, &item)
	return
}

//get the trending repositories of the url(query string )
func trendingRepos(url string) (repos []interface{}, err error) {
	var data map[string]interface{}
	var resp *http.Response
	var items interface{}

	resp, err = httpClient.Get(url)
	if err != nil {
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}
	items, err = unmarshalItem(data, "items")
	if err != nil {
		return
	}
	repos = items.([]interface{})
	return
}

/* strore items given from repositories */
func reposByLang(items []interface{}) map[string]*Language {
	languages := make(map[string]*Language)
	for _, face := range items {
		repo := face.(map[string]interface{})
		if repo["language"] != nil {
			name := repo["language"].(string)
			languages[name] = appendLang(languages[name], repo)
		}
	}
	fmt.Println(languages)
	return languages
}

func appendLang(lang *Language, data map[string]interface{}) (newLang *Language) {
	if lang == nil {
		lang = &Language{}
	}
	login, _ := unmarshalItem(data["owner"].(map[string]interface{}), "login")
	newLang = &Language{
		CountRepos: lang.CountRepos + 1,
		Repositories: append(lang.Repositories, &Repository{
			Name:      data["name"].(string),
			Developer: login.(string),
			RepoURL:   data["html_url"].(string),
			Stars:     int(data["stargazers_count"].(float64)),
			Watchers:  int(data["watchers_count"].(float64)),
			Forks:     int(data["forks_count"].(float64)),
		}),
	}
	return
}
