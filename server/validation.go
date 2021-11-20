package server

import "strconv"

type Validation struct {
	String *struct {
		MinLength *int `yaml:"minLength" json:"minLength"`
		MaxLength *int `yaml:"maxLength" json:"maxLength"`
	} `yaml:"string" json:"string"`
	Number *struct {
		Min *int `yaml:"min" json:"min"`
		Max *int `yaml:"max" json:"max"`
	} `yaml:"number" json:"number"`
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
