package server

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

/*
	since github api does not provides frameworks used in repositories
	get framework process:
			- create a map and define for each language its package filename and
				the list of its few popular frameworks
			- get content of package file according to the language used in the repository
			- decode the content given from response data
			- search in the content for frameworks using simple regex expr
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

var packages = map[string]Package{
	"Python":     Package{fileName: "requirements.txt", frameworks: []string{"Django", "Flask", "Sanic"}},
	"JavaScript": Package{fileName: "package.json", frameworks: []string{"Vue", "React", "Express", "Angular"}},
	"TypeScript": Package{fileName: "package.json", frameworks: []string{"React", "Angular", "Vue", "Nest"}},
	"PHP":        Package{fileName: "composer.json", frameworks: []string{"Laravel", "Symfony", "Zend"}},
	"Ruby":       Package{fileName: "Gemfile", frameworks: []string{"Rails", "Sinatra"}},
	"Java":       Package{fileName: "pom.xml", frameworks: []string{"Spring", "Struts", "Hibernate"}},
}

// package file contents according to language and repo name
func getContentURL(login string, repo string, fileName string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", login, repo, fileName)
}

//get popular frameworks
func Popularframeworks(languages []*Language) error {
	for _, lang := range languages {
		for _, repo := range lang.Repositories {
			if _, exists := packages[lang.Name]; exists {
				arr, err := repoFrameworks(lang.Name, repo)
				if err != nil {
					return err
				}
				storeFrameworks(arr, lang.Name, repo)
			}
		}
	}
	return nil
}

//get frameworks from repository package file using regex
func repoFrameworks(lang string, repo *Repository) (frameworks []string, err error) {
	var content []byte

	content, err = packageContent(lang, repo.Developer, repo.Name)
	if err != nil {
		return
	}
	for _, name := range packages[lang].frameworks {
		pattern := `(?si)\b((` + name + `))\b`
		matched, _ := regexp.MatchString(pattern, string(content))
		if matched {
			frameworks = append(frameworks, name)
		}
	}
	return
}

//get the content of the package file and decode its value
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
	if resp.StatusCode == 403 && data["message"] != nil {
		//error message from response
		err = errors.New(data["message"].(string))
		return
	}
	if data["content"] != nil {
		//decode content by base64
		encoded := data["content"].(string)
		content, err = base64.StdEncoding.DecodeString(encoded) //decode content from base64 to string
	}
	return
}

func storeFrameworks(arr []string, langName string, repo *Repository) {
	for _, name := range arr {
		_, frameworkPtr := FindFramework(G_frameworks, name)
		if frameworkPtr != nil {
			frameworkPtr.CountRepos++
			frameworkPtr.Repositories = append(frameworkPtr.Repositories, repo)
			continue
		}
		G_frameworks = append(G_frameworks, &Framework{
			Name:         name,
			Lang:         langName,
			CountRepos:   1,
			Repositories: []*Repository{repo},
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
