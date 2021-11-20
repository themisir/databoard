package config

import "github.com/themisir/databoard/server"

type Route struct {
	Path    string            `yaml:"path"`
	Methods map[string]Method `yaml:"methods"`
}

type Method struct {
	Parameters map[string]MethodParam `yaml:"parameters"`
	Query      *MethodQueryOptions    `yaml:"query"`
	Mutation   *MethodMutationOptions `yaml:"mutation"`
}

type MethodParam struct {
	Value      string             `yaml:"value"`
	Validation *server.Validation `yaml:"validation"`
}

type MethodQueryOptions struct {
	Name  string `yaml:"name"`
	First bool   `yaml:"first"`
}

type MethodMutationOptions struct {
	Name string `yaml:"name"`
}
