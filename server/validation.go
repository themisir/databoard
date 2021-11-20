package server

import "strconv"

type Validation struct {
	String *struct {
		MinLength *int `yaml:"minLength"`
		MaxLength *int `yaml:"maxLength"`
	} `yaml:"string"`
	Number *struct {
		Min *int `yaml:"min"`
		Max *int `yaml:"max"`
	}
}

func (v Validation) Validate(s string) bool {
	if v.Number != nil {
		num, err := strconv.Atoi(s)
		if err != nil {
			return false
		}
		if v.Number.Min != nil && num < *v.Number.Min {
			return false
		}
		if v.Number.Max != nil && num > *v.Number.Max {
			return false
		}
	}
	if v.String != nil {
		sl := len(s)
		if v.String.MinLength != nil && sl < *v.String.MinLength {
			return false
		}
		if v.String.MaxLength != nil && sl > *v.String.MaxLength {
			return false
		}
	}

	return true
}
