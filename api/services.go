package api

import (
	"net/http"

	"../server"
)

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	s := ServiceResponse{
		status: 200,
		fields: make(map[string]interface{}),
	}
	s.fields["/languages"] = "GET the languages used in the 100 trending repositories with an optional query string since={daily,weekly,montly(default)}"
	s.fields["/frameworks"] = "GET popular frameworks used in the trending github repositories with maximum of 60 repo (rate limit of github api is 60 request per hour)"
	s.exit(w)
}

/* GET /language?since={daily,weekly,monthly}
{
	languages: [
		{
			language: {string}
			repos_count: {int},
			repositories: [
				{
					name: {string},
					developer_login: {string},
					repo_url: {string},
					stars: {int}
					forks: {int}
				},
				...
			]
		},
		...
	]
}
*/
func LanguagesHandler() http.HandlerFunc {
	var err error
	var query map[string]string

	server.G_languages = nil
	s := ServiceResponse{
		status: 200,
		fields: make(map[string]interface{}),
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowedMethods(r.Method, "GET") {
			s.status = 405
			s.exit(w)
			return
		}
		query, err = getQueryValues(r.URL, "since")
		if err != nil {
			s.status = 400
			s.fields["error"] = err.Error()
			s.exit(w)
			return
		}
		server.G_languages, err = server.TrendingLanguages(query["since"])
		if err != nil {
			s.fields["error"] = err.Error()
			s.exit(w)
			return
		}
		s.fields["languages"] = server.G_languages
		server.G_frameworks = nil
		// get frameworks of trending repos
		//		async by running it in a goroutine(thread) to reduce time
		go server.Popularframeworks(server.G_languages)
		s.exit(w)
	}
}

/* GET /frameworks
{
	frameworks: [
		{
			name: {string},
			language: {string},
			repos_count: {int},
			repositories: [
				{
					repo info..
				},
				...
			]
		},
		...
	]
}
*/
func FrameworksHandler() http.HandlerFunc {
	var err error

	s := ServiceResponse{
		status: 200,
		fields: make(map[string]interface{}),
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowedMethods(r.Method, "GET") {
			s.status = 405
			s.exit(w)
			return
		}
		if server.G_frameworks == nil {
			//get trending repositories if framework list is empty
			if server.G_languages == nil {
				server.G_languages, err = server.TrendingLanguages("monthly")
				if err != nil {
					s.fields["error"] = err.Error()
					s.exit(w)
					return
				}
			}
			err = server.Popularframeworks(server.G_languages)
			if err != nil && server.G_frameworks == nil {
				s.fields["error"] = err.Error()
				s.exit(w)
				return
			}
		}
		s.fields["frameworks"] = server.G_frameworks
		s.exit(w)
	}
}
