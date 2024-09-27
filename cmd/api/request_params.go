package main

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/WanderingAura/quotable/internal/validator"
)

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) readBool(qs url.Values, key string, defaultValue bool, v *validator.Validator) bool {
	s := qs.Get(key)

	switch s {
	case "":
		return defaultValue
	case "true":
		return true
	case "false":
		return false
	default:
		v.AddError(key, "must be either true or false")
		return defaultValue
	}
}
