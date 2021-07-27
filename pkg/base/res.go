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
	"cana.io/clap/pkg/model"
	"encoding/json"
	"log"
)

func ResInit() {
	list := dangListFullRes()
	for _, one := range *list {
		val := one
		resIdMap[one.Id] = &val
		resNameMap[one.ResName] = one.Id
		resInfo := one.ResInfo
		if "" != resInfo {
			var info map[string]interface{}
			err := json.Unmarshal([]byte(resInfo), &info)
			if nil != err {
				log.Printf("get res error: %v", err)
			} else {
				resInfoMap[one.Id] = &info
			}
		}
	}
}

func Resources() (*map[uint64]*model.Resource, *map[uint64]*map[string]interface{}) {
	return &resIdMap, &resInfoMap
}

func Resource(id uint64) (*model.Resource, *map[string]interface{}) {
	return resIdMap[id], resInfoMap[id]
}

func dangListFullRes() *[]model.Resource {
	var list []model.Resource
	err := Engine.Omit(model.CreatedAtInResource, model.UpdatedAtInResource).
		OrderBy(model.ResOrderInResource).Find(&list)
	if nil != err {
		panic(err)
	}
	return &list
}
