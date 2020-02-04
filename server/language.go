package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//set timeout of 40 sec in case the remote server is unresponsive
var httpClient = &http.Client{Timeout: 40 * time.Second}

var G_languages []*Language

//Repository fetched fields
type Repository struct {
	Name      string `json:"name"`
	Developer string `json:"developer_login"`
	RepoURL   string `json:"repo_url"`
	Stars     int    `json:"stars"`
	Forks     int    `json:"forks"`
}

type Language struct {
	Name         string        `json:"language"`
	CountRepos   int           `json:"repos_count"`
	Repositories []*Repository `json:"repositories"`
}

/*
	get list of 100 trending repositories
		depending on number of stars && forks
*/
func trendingURL(created time.Time) string {
	return fmt.Sprintf("https://api.github.com/search/repositories?q="+
		"created:>=%s&sort=stars,forks&order=desc&per_page=100", created.Format("2006-01-02"))
}

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

//TrendingLanguages returns a list of the trending languages used in github since(daily|weekly|monthly)
func TrendingLanguages(since string) (languages []*Language, err error) {
	var repos []interface{}

	createdTime := createdSince(since)
	url := trendingURL(createdTime)
	repos, err = trendingRepos(url)
	if err != nil {
		return
	}
	languages = reposByLang(repos)
	return
}

//get the repositories from response of the request(url) passed
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

// filter the repositories taken and classify them by languages
func reposByLang(repositories []interface{}) []*Language {
	languages := []*Language{}
	for _, repo := range repositories {
		items := repo.(map[string]interface{})
		if items["language"] != nil {
			languages = appendLanguage(languages, items)
		}
	}
	return languages
}

func appendLanguage(languages []*Language, repo map[string]interface{}) []*Language {
	name := repo["language"].(string)
	index, _ := FindLanguage(languages, name)
	if index != -1 {
		languages[index].addRepo(repo)
		return languages
	}
	lang := &Language{Name: name}
	lang.addRepo(repo)
	languages = append(languages, lang)
	return languages
}

//counter of any repository added to the list
func (lang *Language) addRepo(repo map[string]interface{}) {
	lang.CountRepos++
	lang.Repositories = append(lang.Repositories, newRepository(repo))
}

//FindLanguage by its name
func FindLanguage(languages []*Language, name string) (int, *Language) {
	for index, lang := range languages {
		if name == lang.Name {
			return index, lang
		}
	}
	return -1, nil
}

//fetch repository fields
func newRepository(data map[string]interface{}) *Repository {
	login, _ := unmarshalItem(data["owner"].(map[string]interface{}), "login")
	return &Repository{
		Name:      data["name"].(string),
		Developer: login.(string),
		RepoURL:   data["html_url"].(string),
		Stars:     int(data["stargazers_count"].(float64)),
		Forks:     int(data["forks_count"].(float64)),
	}
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
