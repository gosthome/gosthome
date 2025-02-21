package util

import (
	"strings"

	"github.com/mozillazg/go-unidecode"
)

func CleanString(s string) string {
	return unidecode.Unidecode(s)
}

func SnakeCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
