package util

import (
	"regexp"
	"strings"
)

var matchAllCamel = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchAllCamel.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

var matchAllSnake = regexp.MustCompile("(_|-)([a-zA-Z]+)")

func ToCamelCase(str string) string {
	camel := matchAllSnake.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)

	return camel
}

func ConvertBoolToYesOrNo(input bool) string {
	if input {
		return "yes"
	} else {
		return "no"
	}
}

func IfNotExistTailAdd(input string, suffix string) string {
	if strings.HasSuffix(input, suffix) {
		return input
	} else {
		return input + suffix
	}
}
