package web

import (
	"net/http"
	"net/url"
)

/*AllowedMethods return
sucesss(200) in case req(request method) match one of the methodes given or
error(method not allowed 405)*/
func AllowedMethods(req string, methods ...string) int {
	for _, mt := range methods {
		if req == mt {
			return 200
		}
	}
	return 405
}

//GetQueryValues return the values
func GetQueryValues(link *url.URL, keywords ...string) (values map[string]string, err error) {
	//u, err := url.Parse(link)
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
		return "error"
	}
}

func JsonResponse(status int, title string, data []byte) []byte {
	var resp string

	if status != http.StatusOK {
		resp = `"status":"` + statusDescription(status) + `",`
	}
	if len(data) > 0 {
		if len(title) > 0 {
			resp += `"` + title + `":` + string(data)
		} else if len(resp) == 0 {
			return data
		}
	}
	return []byte(`{` + resp + `}`)
}
