package server

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
)

//set timeout of 15 sec in case the remote server is unresponsive
var httpClient = &http.Client{Timeout: 15 * time.Second}

type Repository struct {
	Name 		  string `json:"name"`
	Developer 	  string `json:"developer_login"`
	RepoURL       string `json:"repo_url"`
	Stars         int    `json:"stars"`
	Watchers      int    `json:"watchers"`
	Forks         int    `json:"forks"`
}

type Language struct {
	CountRepos 	 int           `json:"repos_count"`
	Repositories []*Repository `json:"repositories"`
}


//get 100 trending repos of all languages
func trendingURL(since time.Time) string {
	return fmt.Sprintf("https://api.github.com/search/repositories?q=" +
					"created:>=%s&sort=stars&order=desc&per_page=100", since.Format("2006-01-02"))
}

//get 100 trending repos by specific language
func trendingByLangURL(since time.Time, lang string) string {
	return fmt.Sprintf("https://api.github.com/search/repositories?q=" +
					"created:>=%s&language=%s&sort=stars&order=desc&per_page=100", since.Format("2006-01-02"), lang)
}

//get trending repositories by languages
func getLanguages(url string) (languages map[string]*Language, err error) {
	var repos []map[string]interface{}

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
func trendingRepos(url string) (repos []map[string]interface{}, err error) {
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
	repos = items.([]map[string]interface{})
	return
}

/* strore items given from repositories */
func reposByLang(items []map[string]interface{}) map[string]*Language {
	languages := make(map[string]*Language)
	for _, repo := range items {
		name := repo["language"].(string)
		languages[name].append(repo)
	}
	return languages
}

func (lang *Language) append(data map[string]interface{}) {
	lang.CountRepos++
	login, _ := unmarshalItem(data["owner"].(map[string]interface{}), "login")
	lang.Repositories = append(lang.Repositories, &Repository {
		Name: data["name"].(string),
		Developer: login.(string),
		RepoURL: data["html_url"].(string),
		Stars: data["stargazers_count"].(int),
		Watchers: data["watchers_count"].(int),
		Forks: data["forks_count"].(int),
	})
}