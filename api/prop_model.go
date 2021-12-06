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
	"xorm.io/xorm"
)

func getPropsById(id uint64) (*model.PropertyFile, bool, error) {
	var info model.PropertyFile
	result, err := base.Engine.ID(id).Get(&info)
	return &info, result, err
}

func getPropsByLink(id, res uint64) (*[]model.PropertyFile, error) {
	var list []model.PropertyFile
	err := base.Engine.Omit(model.CreatedAt, model.UpdatedAt).
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

func insertProp(session *xorm.Session, info *model.PropertyFile) (int64, error) {
	return session.Omit(model.IdInPropertyFile, model.CreatedAt, model.UpdatedAt).
		InsertOne(info)
}

func updatePropById(session *xorm.Session, id uint64, info *model.PropertyFile) (int64, error) {
	return session.Cols(model.FileReadmeInPropertyFile, model.FileContentInPropertyFile, model.FileHashInPropertyFile).
		Update(info, model.PropertyFile{Id: id})
}

func updatePropStatusById(session *xorm.Session, id uint64, disable int) (int64, error) {
	res, err := session.Cols(model.IsDisableInPropertyFile).
		Exec("update "+model.PropertyFileTable+" set "+model.IsDisableInPropertyFile+" = ? "+
			"where "+model.IdInPropertyFile+" = ?", disable, id)
	if nil != err {
		return 0, err
	}
	return res.RowsAffected()
}

func findPropByIds(ids []uint64) (*model.PropertySnap, error) {
	var info model.PropertySnap
	err := base.Engine.Omit(model.CreatedAt).
		In(model.IdInPropertySnap, ids).
		Find(&info)
	return &info, err
}

func getLatestPropSnap(res, prop uint64) (*model.PropertySnap, error) {
	var info model.PropertySnap
	err := base.Engine.Omit(model.CreatedAt).
		Where(model.ResIdInPropertySnap+"?", res).Where(model.PropIdInPropertySnap+"?", prop).
		Desc(model.IdInPropertySnap).Limit(1).
		Find(&info)
	return &info, err
}

func insertPropSnap(session *xorm.Session, info *model.PropertySnap) (int64, error) {
	return session.Omit(model.IdInPropertySnap, model.CreatedAt).
		InsertOne(info)
}
