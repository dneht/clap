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
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"xorm.io/xorm"
)

func GetProp(c *fiber.Ctx) error {
	get := c.Params("type")
	if "" == get {
		return util.ErrorInput(c, "type input error")
	}
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "id input error")
	}
	res, ok := base.ResourceId(CommonPre + get)
	if !ok {
		return util.ErrorInput(c, "type get error")
	}
	err = ResourceAuth(c, get, id, AllowThisPropertyView)
	if nil != err {
		return err
	}

	list, err := getPropsByLink(id, res)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	files := make([]refer.PropBaseOutput, 0, len(list))
	for _, one := range list {
		files = append(files, refer.PropBaseOutput{
			Id:      one.Id,
			Name:    one.FileName,
			Readme:  one.FileReadme,
			Content: one.FileContent,
		})
	}
	return util.ResultList(c, err, len(files), files)
}

func CreateProp(c *fiber.Ctx) error {
	if !ManagerAuth(c) {
		return fiber.ErrForbidden
	}
	info, err := checkPropertyAdd(c)
	if err != nil {
		return err
	}
	info.FileHash = getPropHashByContent(info.FileContent)

	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err = insertProp(session, info)
		if nil != err {
			return 0, util.ErrorInternal(c, err)
		}
		_, err = insertConfigSnap(session, &model.PropertySnap{
			UserId:      contextUserId(c),
			ResId:       info.ResId,
			LinkId:      info.LinkId,
			PropId:      info.Id,
			FileName:    info.FileName,
			FileContent: info.FileContent,
		})
		if nil != err {
			return 0, util.ErrorInternal(c, err)
		}
		return 1, nil
	})
	return util.ResultParamWithMessage(c, err, true, "create prop error", info.Id)
}

func UpdateProp(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "id input error")
	}
	info, err := checkPropertyBase(c)
	if nil != err {
		return err
	}
	get, result, err := getPropsById(id)
	if nil != err || !result {
		return util.ErrorInternal(c, err)
	}
	needKube, resName, envBase, spaceBase, appBase, deployBase, err := checkNeedUpdateKube(get)
	if nil != err {
		return util.ErrorInputLog(c, err, "res get error")
	}
	if needKube {
		if nil == envBase || nil == spaceBase || nil == appBase || nil == deployBase {
			return util.ErrorInputLog(c, err, "res not found")
		}
	}
	err = ResourceAuth(c, resName, get.LinkId, AllowThisPropertyUpdate)
	if nil != err {
		return err
	}
	var k8s *kubernetes.Clientset
	if needKube {
		k8s, _, err = base.K8S(envBase.Id)
		if nil != err {
			return util.ErrorInput(c, "can not get cluster info")
		}
	}

	fileHash := getPropHashByContent(info.FileContent)
	if info.FileReadme == get.FileReadme && fileHash == get.FileHash {
		return c.SendString("ok")
	}
	info.FileHash = fileHash

	userId := contextUserId(c)
	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err = updatePropById(session, id, info)
		if nil != err {
			return 0, util.ErrorInternal(c, err)
		}
		if "" != get.FileContent && fileHash != get.FileHash {
			_, err = insertConfigSnap(session, &model.PropertySnap{
				UserId:      userId,
				ResId:       get.ResId,
				LinkId:      get.LinkId,
				PropId:      get.Id,
				FileName:    get.FileName,
				FileContent: info.FileContent,
			})
			if nil != err {
				return 0, util.ErrorInternal(c, err)
			}
		}
		return 1, nil
	})

	if needKube {
		propMap := generateNeedProps(appBase.Id, envBase.Id, spaceBase.Id, deployBase.Id)
		if nil != propMap && len(propMap) > 0 {
			configName := refer.GetConfigName(appBase, spaceBase, refer.PropGenerateName)
			err = onlyUpdateConfig(k8s, spaceBase.SpaceKeep, configName, &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: configName,
				},
				Data: propMap,
			})
			if nil != err {
				return util.ErrorInternal(c, err)
			}
			invalidAppInfoById(appBase.Id)
			invalidDeployInfoById(deployBase.Id)
		}
	}
	return util.ResultParamWithMessage(c, err, true, "update prop error", info.Id)
}

func checkPropertyBase(c *fiber.Ctx) (*model.PropertyFile, error) {
	info := new(model.PropertyFile)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputLog(c, err, "input is error")
	}
	if "" == info.FileContent {
		return nil, util.ErrorInput(c, "prop content set error")
	}
	info.FileContent = formatPropBySpace(info.FileContent)
	return info, nil
}

func checkPropertyAdd(c *fiber.Ctx) (*model.PropertyFile, error) {
	info, err := checkPropertyBase(c)
	if nil != err {
		return nil, err
	}

	get := c.Params("type")
	if "" == get {
		return nil, util.ErrorInput(c, "type input error")
	}
	res, ok := base.ResourceId(CommonPre + get)
	if ok {
		info.ResId = res
	} else {
		return nil, util.ErrorInput(c, "type not found")
	}
	if info.LinkId <= 0 {
		return nil, util.ErrorInput(c, "link id set error")
	}
	if "" == info.FileName {
		return nil, util.ErrorInput(c, "file name set error")
	}
	return info, nil
}

func checkNeedUpdateKube(info *model.PropertyFile) (bool, string, *model.Environment, *model.EnvironmentSpace, *model.Project, *model.Deployment, error) {
	resBase, _ := base.Resource(info.ResId)
	if nil == resBase {
		return false, "", nil, nil, nil, nil, nil
	}
	if resBase.ResName == CommonPre+model.DeploymentTable {
		envBase, spaceBase, appBase, deployBase, err := getMoreModels(info.LinkId)
		if nil != err {
			return false, "", nil, nil, nil, nil, nil
		}
		return true, strings.Replace(resBase.ResName, CommonPre, "", 1),
			envBase, spaceBase, appBase, deployBase, nil
	}
	return false, "", nil, nil, nil, nil, nil
}

func generateNeedProps(appId, envId, spaceId, deployId uint64) map[string]string {
	allList := make([]model.PropertyFile, 0, 8)
	existNames := make(map[string]bool)
	deployList := getPropsByLinkWithName(deployId, model.DeploymentTable)
	if nil != deployList && len(deployList) > 0 {
		allList = append(allList, deployList...)
		for _, deployOne := range deployList {
			name := strings.TrimSpace(deployOne.FileName)
			if "" != name {
				existNames[name] = true
			}
		}
	}

	spaceList := getPropsByLinkWithName(spaceId, model.EnvironmentSpaceTable)
	if nil != spaceList && len(spaceList) > 0 {
		allList = append(allList, appendNeedProps(existNames, spaceList)...)
	}
	envList := getPropsByLinkWithName(envId, model.EnvironmentTable)
	if nil != envList && len(envList) > 0 {
		allList = append(allList, appendNeedProps(existNames, envList)...)
	}
	appList := getPropsByLinkWithName(appId, model.ProjectTable)
	if nil != appList && len(appList) > 0 {
		allList = append(allList, appendNeedProps(existNames, appList)...)
	}
	if len(allList) == 0 {
		return nil
	}

	allData := mergePropByName(allList)
	if nil == allData || len(allData) == 0 {
		allData = make(map[string]string, 0)
		return allData
	}
	return allData
}

func appendNeedProps(names map[string]bool, props []model.PropertyFile) []model.PropertyFile {
	list := make([]model.PropertyFile, 0, len(props))
	for _, prop := range props {
		name := strings.TrimSpace(prop.FileName)
		_, ok := names[name]
		if ok {
			list = append(list, prop)
		}
	}
	return list
}

func generateRenderProps(appId, envId, spaceId, deployId uint64, appInfo *refer.AppInfo) error {
	if nil == appInfo || nil == appInfo.Param {
		return errors.New("get app info and param error")
	}
	var volumeGet []interface{}
	volumeName := refer.PropGenerateName
	volumeMounts, ok := appInfo.Param[refer.KeyVolumeMounts]
	if ok {
		volumeGet, ok = volumeMounts.([]interface{})
		if ok {
			for _, volumeOne := range volumeGet {
				volumeThis, cok := volumeOne.(refer.VolumeMountInfo)
				if cok && volumeThis.Name == volumeName {
					return nil
				}
			}
		} else {
			return errors.New("convert exist volumes error")
		}
	}

	allData := generateNeedProps(appId, envId, spaceId, deployId)
	if nil == allData || len(allData) == 0 {
		return nil
	}

	mountPath := refer.PropMountPath
	if nil != appInfo.Factor && "" != appInfo.Factor.ConfigMouthPath {
		mountPath = appInfo.Factor.ConfigMouthPath
	}
	volumeList := []interface{}{
		refer.VolumeMountInfo{
			Name:      volumeName,
			Type:      "Config",
			Data:      allData,
			MountPath: mountPath,
			ReadOnly:  true,
		},
	}

	if nil != volumeGet {
		volumeList = append(volumeList, volumeGet...)
	}
	appInfo.Param[refer.KeyVolumeMounts] = volumeList
	return nil
}
