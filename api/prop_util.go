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

package api

import (
	"cana.io/clap/pkg/model"
	"cana.io/clap/util"
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

func splitPropBySpace(content string) []string {
	if "" == content {
		return []string{}
	}
	lines := strings.Split(content, "\n")
	if len(lines) == 1 {
		lines = strings.Split(content, "\r")
	}
	return lines
}

func formatPropBySpace(content string) string {
	if "" == content {
		return ""
	}
	lines := strings.Split(content, "\n")
	if len(lines) == 1 {
		lines = strings.Split(content, "\r")
	}
	return strings.Join(util.RemoveRepeatedElement(lines, "="), "\n")
}

func getPropHashByContent(content string) string {
	lines := splitPropBySpace(content)
	if len(lines) == 0 {
		return ""
	}

	sort.Strings(lines)
	hash := md5.New()
	for _, line := range lines {
		if "" == line {
			continue
		}
		hash.Write([]byte(strings.TrimSpace(line)))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func mergePropByName(list []model.PropertyFile) map[string]string {
	files := make(map[string][]string, 8)
	for _, one := range list {
		if "" == one.FileContent {
			continue
		}
		getList, ok := files[one.FileName]
		newList := splitPropBySpace(one.FileContent)
		if ok {
			files[one.FileName] = append(getList, newList...)
		} else {
			files[one.FileName] = newList
		}
	}

	table := make(map[string]string, len(list))
	for key, value := range files {
		table[key] = strings.Join(util.RemoveRepeatedElement(value, "="), "\n")
	}
	return table
}

func getPropReadme(userId uint64, created *time.Time, rollback bool) string {
	var readme string
	user, err := getUserById(userId)
	if nil == err {
		if rollback {
			readme = rollbackPropReadme(user.Nickname, created)
		} else {
			readme = updatePropReadme(user.Nickname, created)
		}
	} else {
		if rollback {
			readme = rollbackPropReadme("", created)
		} else {
			readme = updatePropReadme("", created)
		}
	}
	return readme
}

func updatePropReadme(name string, created *time.Time) string {
	if nil == created {
		return ""
	}

	if name == "" {
		return "配置在" + formatPropCreatedAt(created) + "被更新"
	} else {
		return "配置在" + formatPropCreatedAt(created) + "由[" + name + "]更新"
	}
}

func rollbackPropReadme(name string, created *time.Time) string {
	if nil == created {
		return ""
	}

	if name == "" {
		return "已回滚到" + formatPropCreatedAt(created) + "创建的配置"
	} else {
		return "已回滚到由[" + name + "]在" + formatPropCreatedAt(created) + "创建的配置"
	}
}

func formatPropCreatedAt(created *time.Time) string {
	return created.Format("2006年01月02日15点04分")
}
