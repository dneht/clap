/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base

import (
	"strconv"
	"strings"
)

func parseString(value, def string) string {
	value = strings.TrimSpace(value)
	if "" == value {
		return def
	}
	return value
}

func parseBool(value string, def bool) bool {
	value = strings.TrimSpace(value)
	if "" == value {
		return def
	}
	return "true" == value || "yes" == value
}

func parseInt(value string, def int) int {
	value = strings.TrimSpace(value)
	if "" == value {
		return def
	}
	res, err := strconv.ParseInt(value, 10, 32)
	if nil != err {
		return def
	}
	return int(res)
}

func parseList(value string) []string {
	list := make([]string, 0, 4)
	value = strings.TrimSpace(value)
	if "" == value {
		return list
	}
	split := strings.Split(value, ",")
	for _, one := range split {
		add := strings.TrimSpace(one)
		if "" != add {
			list = append(list, add)
		}
	}
	return list
}

func parseListElse(value string, def []string) []string {
	list := parseList(value)
	if nil == def {
		return list
	}
	if len(list) == 0 {
		return def
	}
	return list
}

func parseDocument(kv map[string]string) map[string]DocumentProp {
	return parseMapPre("document", kv)
}

func parseMapPre(pre string, kv map[string]string) map[string]DocumentProp {
	props := make(map[string]DocumentProp, 4)
	for key, value := range kv {
		key = strings.TrimSpace(key)
		if strings.HasPrefix(key, pre) {
			split := strings.Split(key, ".")
			if len(split) >= 4 {
				fst := split[1] + "_" + split[2]
				get, ok := props[fst]
				if !ok {
					get = DocumentProp{}
				}
				switch split[3] {
				case "enable":
					if !get.Enable {
						get.Enable = parseBool(value, false)
					}
					break
				case "token":
					if "" == get.Token {
						get.Token = parseString(value, "")
					}
					break
				case "api_url":
					if "" == get.ApiUrl {
						get.ApiUrl = parseString(value, "")
					}
					break
				}
				props[fst] = get
			}
		}
	}
	return props
}
