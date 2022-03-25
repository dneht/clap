package util

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
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

func SpliceNumbersToMD5(nums ...uint64) string {
	hash := md5.New()
	for idx, num := range nums {
		hash.Write([]byte(strconv.FormatInt(int64(idx), 10)))
		hash.Write([]byte(strconv.FormatUint(num, 10)))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func GenerateMD5(list ...string) string {
	hash := md5.New()
	for idx, one := range list {
		hash.Write([]byte(strconv.FormatInt(int64(idx), 10)))
		hash.Write([]byte(one))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func StringClone(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}
