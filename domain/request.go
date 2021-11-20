package domain

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	Path    string
	Headers map[string]string
	Query   map[string]string
	Vars    map[string]string
	Body    map[string]interface{}
}

func requestFromHttp(r *http.Request) *Request {
	body := make(map[string]interface{})
	json.NewDecoder(r.Body).Decode(&body)

	return &Request{
		Path:    r.URL.Path,
		Headers: flatValues(r.Header),
		Query:   flatValues(r.URL.Query()),
		Vars:    mux.Vars(r),
		Body:    body,
	}
}

func flatValues(v map[string][]string) (r map[string]string) {
	r = make(map[string]string, len(v))
	for k, v := range v {
		r[k] = v[0]
	}
	return
}
