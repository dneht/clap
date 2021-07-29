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
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

func getPropsById(id uint64) (*model.PropertyFile, bool, error) {
	var info model.PropertyFile
	result, err := base.Engine.ID(id).Get(&info)
	return &info, result, err
}

func getPropsByLink(id, res uint64) (*[]model.PropertyFile, error) {
	var list []model.PropertyFile
	err := base.Engine.Omit(model.CreatedAtInPropertyFile, model.UpdatedAtInPropertyFile).
		Where(model.ResIdInPropertyFile+" = ?", res).
		Where(model.LinkIdInPropertyFile+" = ?", id).
		Where(model.IsDisableInPropertyFile + " = 0").Find(&list)
	return &list, err
}

func getPropsByLinkWithName(id uint64, name string) *[]model.PropertyFile {
	res, ok := base.ResourceId(CommonPre + name)
	if ok {
		list, err := getPropsByLink(id, res)
		if nil == err {
			return list
		}
	}
	return nil
}

func insertProp(c *fiber.Ctx, session *xorm.Session, info *model.PropertyFile) (int64, error) {
	return session.Omit(model.IdInPropertyFile, model.CreatedAtInUserInfo, model.UpdatedAtInUserInfo).
		InsertOne(info)
}

func updatePropById(c *fiber.Ctx, session *xorm.Session, id uint64, info *model.PropertyFile) (int64, error) {
	return session.Cols(model.FileReadmeInPropertyFile, model.FileContentInPropertyFile, model.FileHashInPropertyFile).
		Update(info, model.PropertyFile{Id: id})
}

func updatePropStatusById(c *fiber.Ctx, session *xorm.Session, id uint64, disable int) (int64, error) {
	res, err := session.Cols(model.IsDisableInPropertyFile).
		Exec("update "+model.PropertyFileTable+" set "+model.IsDisableInPropertyFile+" = ? "+
			"where "+model.IdInPropertyFile+" = ?", disable, id)
	if nil != err {
		return 0, err
	}
	return res.RowsAffected()
}

func findPropByIds(c *fiber.Ctx, ids []uint64) (*model.PropertySnap, error) {
	var info model.PropertySnap
	err := base.Engine.Omit(model.CreatedAtInPropertySnap).
		In(model.IdInPropertySnap, ids).
		Find(&info)
	return &info, err
}

func getLatestPropSnap(c *fiber.Ctx, res, prop uint64) (*model.PropertySnap, error) {
	var info model.PropertySnap
	err := base.Engine.Omit(model.CreatedAtInPropertySnap).
		Where(model.ResIdInPropertySnap+"?", res).Where(model.PropIdInPropertySnap+"?", prop).
		Desc(model.IdInPropertySnap).Limit(1).
		Find(&info)
	return &info, err
}

func insertPropSnap(c *fiber.Ctx, session *xorm.Session, info *model.PropertySnap) (int64, error) {
	return session.Omit(model.IdInPropertySnap, model.CreatedAtInPropertySnap).
		InsertOne(info)
}

func generateNeedProps(appId, envId, spaceId, deployId uint64) *map[string]string {
	allList := make([]model.PropertyFile, 0, 8)
	deployList := getPropsByLinkWithName(deployId, model.DeploymentTable)
	if nil != deployList {
		allList = append(allList, *deployList...)
	}
	spaceList := getPropsByLinkWithName(spaceId, model.EnvironmentSpaceTable)
	if nil != spaceList {
		allList = append(allList, *spaceList...)
	}
	envList := getPropsByLinkWithName(envId, model.EnvironmentTable)
	if nil != envList {
		allList = append(allList, *envList...)
	}
	appList := getPropsByLinkWithName(appId, model.ProjectTable)
	if nil != appList {
		allList = append(allList, *appList...)
	}
	if len(allList) == 0 {
		return nil
	}

	allData := mergePropByName(&allList)
	if len(allData) == 0 {
		return nil
	}
	return &allData
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

	mountPath := refer.PropMountPath
	if nil != appInfo.Factor && "" != appInfo.Factor.ConfigMouthPath {
		mountPath = appInfo.Factor.ConfigMouthPath
	}

	allData := *generateNeedProps(appId, envId, spaceId, deployId)
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
