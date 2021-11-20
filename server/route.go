package server

import (
	"log"
	"net/http"
	"text/template"

	"github.com/themisir/databoard/domain"
)

type Route struct {
	path    string
	methods map[string]*Method
}

type Method struct {
	delegate MethodDelegate
	params   map[string]methodParam
}

type methodParam struct {
	tmpl       *template.Template
	validation *Validation
}

type MethodDelegate interface {
	Handle(w http.ResponseWriter, ctx *domain.Context, values map[string]string)
}

// Create a new route for given path
func NewRoute(path string) *Route {
	return &Route{
		path:    path,
		methods: make(map[string]*Method),
	}
}

// Add parameter to the method using given parameters. Returns error
// if parsing given valueTemplate fails.
func (m *Method) AddParam(name string, valueTemplate string, validation *Validation) error {
	tmpl, err := template.New("param." + name).Parse(valueTemplate)
	if err != nil {
		return err
	}
	m.params[name] = methodParam{tmpl: tmpl, validation: validation}
	return nil
}

// Returns route path
func (route *Route) Path() string {
	return route.path
}

// Adds handler for provided method name and returns pointer to created handler
func (route *Route) AddMethod(name string, delegate MethodDelegate) (method *Method) {
	method = &Method{
		delegate: delegate,
		params:   make(map[string]methodParam),
	}
	route.methods[name] = method
	return
}

// Handler that's used by net/http package to handle http requests
func (route *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if method, ok := route.methods[r.Method]; ok {
		context := domain.New(r)
		values := make(map[string]string, len(method.params))

		for k, v := range method.params {
			value, err := context.Transform(v.tmpl)
			if err != nil {
				// TODO: Log instead of panicking
				log.Fatalf("Failed to execute parameter '%s' template on route (%s) %s: %s", k, r.Method, r.URL.Path, err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if v.validation != nil && !v.validation.Validate(value) {
				// TODO: Log and respond validation errors
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			values[k] = value
		}

		method.delegate.Handle(w, context, values)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
