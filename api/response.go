package api

import (
	"net/http"
	"net/url"
	"encoding/json"
)

type ServiceResponse struct {
	status int
	fields map[string]interface{}
}

func (resp ServiceResponse) exit(w http.ResponseWriter) {
	setHeader(w, resp.status)
	if resp.status != 200 {
		resp.fields["status"] = http.StatusText(resp.status)
	}
	data, _ := json.Marshal(resp.fields)
	w.Write(data)
}

/*AllowedMethods return
sucesss(200) in case req(request method) match one of the methodes given or
error(method not allowed 405)*/
func allowedMethods(req string, methods ...string) bool {
	for _, mt := range methods {
		if req == mt {
			return true
		}
	}
	return false
}

//GetQueryValues return the values
func getQueryValues(link *url.URL, keywords ...string) (values map[string]string, err error) {
	if err != nil {
		return nil, err
	}
	query := link.Query()
	values = make(map[string]string)
	for _, key := range keywords {
		values[key] = query.Get(key)
	}
	return
}

func setHeader(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
}

func statusDescription(status int) string {
	switch status {
	case 200:
		return "succeed"
	case 201:
		return "new element has been created"
	case 405:
		return "method not allowed"
	case 400:
		return "bad request : The server could not understand the request"
	default:
		return ""
	}
}
