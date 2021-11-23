package domain

import (
	"bytes"
	"net/http"
	"strconv"
	"text/template"
)

type Context struct {
	r *Request
}

type Pagination struct {
	Offset int
	Limit  int
}

// Create new Context for given request
func New(r *http.Request) *Context {
	return &Context{r: requestFromHttp(r)}
}

// Execute template on this context
func (c *Context) Transform(tmpl *template.Template) (string, error) {
	buff := new(bytes.Buffer)
	if err := tmpl.Execute(buff, c); err != nil {
		return "", err
	}
	return buff.String(), nil
}

// Returns request details passed throught NewContext
func (c *Context) Req() *Request {
	return c.r
}

// Returns pagination details built using _offset and _limit parameters from
// http request query. This function returns nil if fails to parse those
// parameters.
func (c *Context) Pagination() *Pagination {
	offset, err1 := strconv.Atoi(c.r.Query["_offset"])
	limit, err2 := strconv.Atoi(c.r.Query["_limit"])
	if err1 == nil && err2 == nil {
		return &Pagination{Offset: offset, Limit: limit}
	}
	return nil
}
