package server

import (
	"net/http"

	"../settings"
)

//WebServer : listening from host:port given in conf file */
func WebServer(route *http.ServeMux) error {
	var values settings.Values

	values = settings.Restore()
	//route.Handle("/", http.FileServer(s.Root))
	addr := values.Host + ":" + values.Port
	return http.ListenAndServe(addr, nil)
}
