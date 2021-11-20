package config

type DataboardConfig struct {
	Database  DbConfig         `yaml:"db"`
	Queries   map[string]Query `yaml:"queries"`
	Mutations map[string]Query `yaml:"mutations"`
	Routes    []Route          `yaml:"routes"`
}

type DbConfig struct {
	Driver     string `yaml:"driver"`
	Connection string `yaml:"connection"`
}

type Query struct {
	Query      string     `yaml:"query"`
	Parameters []SqlParam `yaml:"parameters"`
}

type SqlParam struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Optional bool   `yaml:"optional"`
}
