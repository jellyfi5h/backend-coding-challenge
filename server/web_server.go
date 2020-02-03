package server

import (
	"net/http"

	"../settings"
)

func entryPoints() *http.ServeMux {
	route := http.NewServeMux()

	route.HandleFunc("/languages", languagesHandler) // /languages?since={Daily | weekly | Monthly(default)}
	route.HandleFunc("/frameworks", frameworksHandler)
	return route
}

//WebServer : listening from host:port given in conf file */
func WebServer() error {
	route := entryPoints()
	values := settings.Restore()
	addr := values.Host + ":" + values.Port
	return http.ListenAndServe(addr, route)
}
