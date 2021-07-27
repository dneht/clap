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
	"cana.io/clap/pkg/refer"
)

var spaceMap = make(map[uint64]*model.EnvironmentSpace)
var spaceInfoMap = make(map[uint64]*refer.SpaceInfo)

var appMap = make(map[uint64]*model.Project)
var appInfoMap = make(map[uint64]*refer.AppInfo)

var deployMap = make(map[uint64]*model.Deployment)
var appDeployMap = make(map[uint64]*refer.AppInfo)

var templateMap = make(map[uint64]*model.Template)

var permissionMap = make(map[uint64]*model.Permission)

var authInfoMap = make(map[string]*refer.AuthInfo)

func resetAllCache() {
	spaceMap = make(map[uint64]*model.EnvironmentSpace)
	spaceInfoMap = make(map[uint64]*refer.SpaceInfo)

	appMap = make(map[uint64]*model.Project)
	appInfoMap = make(map[uint64]*refer.AppInfo)

	deployMap = make(map[uint64]*model.Deployment)
	appDeployMap = make(map[uint64]*refer.AppInfo)

	templateMap = make(map[uint64]*model.Template)

	permissionMap = make(map[uint64]*model.Permission)

	authInfoMap = make(map[string]*refer.AuthInfo)
}
