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

package util

import "strings"

func RemoveRepeatedElement(arr []string, sp string) []string {
	result := make([]string, 0, len(arr))
	repeat := make(map[string]bool, len(arr))
	for i := 0; i < len(arr); i++ {
		get := strings.TrimSpace(arr[i])
		if "" == get || strings.HasPrefix(get, "=") {
			continue
		}
		key := get
		if "" != sp {
			split := strings.Split(get, sp)
			if len(split) < 2 {
				continue
			}
			key = strings.TrimSpace(split[0])
			get = key + sp + strings.TrimSpace(split[1])
		}
		_, ok := repeat[key]
		if !ok {
			result = append(result, get)
			repeat[key] = true
		}
	}
	return result
}
