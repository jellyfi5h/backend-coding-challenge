package server

import (
	"net/http"

	"../settings"
)

/* response example:
{
	lang_name: {
		repos_count: 40
		repositories: [
			{
				name: ...
				url: ...
				owner: ...
				starred: int
				forked: int
				watched: int
			},
			...
		]
	},
	...
}
*/

func entryPoints() *http.ServeMux {
	route := http.NewServeMux()

	route.HandleFunc("/languages", languagesHandler) // /languages?since={Daily | weekly | Monthly(default)}
	route.HandleFunc("/repos", reposByLangHandler)   // /repos/:language?since={..}&lang={name}
	return route
}

//WebServer : listening from host:port given in conf file */
func WebServer() error {
	route := entryPoints()
	values := settings.Restore()
	addr := values.Host + ":" + values.Port
	return http.ListenAndServe(addr, route)
}
