package query

import "strconv"

// Parameters used for mapping data from value map into SQL
// arguments when querying database.
type Parameter interface {
	Name() string
	IsOptional() bool
	Transform(s string) (interface{}, error)
}

// Create a new string parameter with given name
func String(name string, optional bool) Parameter {
	return stringParam{name, optional}
}

// Create a new int parameter with given name
func Int(name string, optional bool) Parameter {
	return intParam{name, optional}
}

type stringParam struct {
	name     string
	optional bool
}

type intParam struct {
	name     string
	optional bool
}

func (p stringParam) Name() string {
	return p.name
}

func (p stringParam) Transform(s string) (interface{}, error) {
	return s, nil
}

func (p stringParam) IsOptional() bool {
	return p.optional
}

func (p intParam) Name() string {
	return p.name
}

func (p intParam) Transform(s string) (interface{}, error) {
	return strconv.Atoi(s)
}

func (p intParam) IsOptional() bool {
	return p.optional
}
