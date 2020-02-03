package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
)

/*
	since github api does not provides frameworks used in repositories
	i will search in every repository in their package file (accord the language) if exist and only in root (not in subfolders)
	for specific frameworks (according to language) using a simple regex expression.
	and get's its popularity by the number of repos using it
*/

var G_frameworks []*Framework

type Package struct {
	fileName   string
	frameworks []string
}

type Framework struct {
	Name         string        `json:"name"`
	Lang         string        `json:"language"`
	CountRepos   int           `json:"repos_count"`
	Repositories []*Repository `json:"repositories"`
}

//list of frameworks in each language and the package file where to find them
var packages = map[string]Package{
	"Python":     Package{fileName: "requirements.txt", frameworks: []string{"Django", "Flask", "Sanic", "Tornado"}},
	"JavaScript": Package{fileName: "package.json", frameworks: []string{"Vue", "React", "Express", "Angular", "Ember"}},
	"TypeScript": Package{fileName: "package.json", frameworks: []string{"React", "Angular", "Vue", "Nest", "Loopback", "Stix"}},
	"PHP":        Package{fileName: "composer.json", frameworks: []string{"Laravel", "Symfony", "Zend", "Phalcon"}},
	"Ruby":       Package{fileName: "Gemfile", frameworks: []string{"Rails", "Sinatra"}},
}

// package file contents according to language and repo name
func getContentURL(login string, repo string, fileName string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", login, repo, fileName)
}

//get popular frameworks from the list of the trending 100 repositories
func Popularframeworks(languages []*Language) {
	for _, lang := range languages {
		for _, repo := range lang.Repositories {
			if _, exists := packages[lang.Name]; exists { // if this language is
				arr := repoFrameworks(lang.Name, repo)
				storeFrameworks(arr, lang.Name)
			}
		}
	}
}

//search in the content of the package file of the given repository for frameworks using regex exp
func repoFrameworks(lang string, repo *Repository) []string {
	frameworks := []string{}
	content, _ := packageContent(lang, repo.Developer, repo.Name)
	for _, name := range packages[lang].frameworks {
		pattern := `(?si)\b((` + name + `))\b`
		matched, _ := regexp.MatchString(pattern, string(content))
		if matched {
			frameworks = append(frameworks, name)
		}
	}
	return frameworks
}

//get the content of the package file by the language given in a specific repository
func packageContent(lang string, login string, repo string) (content []byte, err error) {
	var data map[string]interface{}

	url := getContentURL(login, repo, packages[lang].fileName)
	resp, err := httpClient.Get(url)
	if err != nil {
		return // timeout request
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}
	if data["content"] != nil { // file found in root dir
		encoded := data["content"].(string)
		content, err = base64.StdEncoding.DecodeString(encoded) //decode content from base64 to string
	}
	return
}

func storeFrameworks(arr []string, langName string) {
	for _, name := range arr {
		id, _ := FindFramework(G_frameworks, name)
		if id != -1 {
			G_frameworks[id].CountRepos++
			continue
		}
		G_frameworks = append(G_frameworks, &Framework{
			Name:       name,
			Lang:       langName,
			CountRepos: 1,
		})
	}
}

//FindFramework by its name
func FindFramework(frameworks []*Framework, name string) (int, *Framework) {
	for id, farmew := range frameworks {
		if farmew.Name == name {
			return id, farmew
		}
	}
	return -1, nil
}
