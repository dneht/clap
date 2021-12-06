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
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

var roleMap = make(map[uint64]*model.RoleInfo)

func getRoleById(id uint64) (*model.RoleInfo, error) {
	value, ok := roleMap[id]
	if ok {
		return value, nil
	}

	var info model.RoleInfo
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		roleMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found role info")
	}
	return &info, err
}

func invalidRoleById(id uint64) {
	delete(roleMap, id)
}

func countRoleWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.RoleInfo
	sql := base.Engine.Cols()
	err := SelectAuth(c, model.RoleInfoTable, sql)
	if nil != err {
		return 0, err
	}
	if !ManagerAuth(c) {
		input.ApplyWithoutDisable(sql)
	}
	return sql.Count(&info)
}

func findRoleWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.RoleInfo, error) {
	var list []model.RoleInfo
	sql := base.Engine.Omit(model.CreatedAt, model.UpdatedAt)
	err := SelectAuth(c, model.RoleInfoTable, sql)
	if nil != err {
		return 0, nil, err
	}
	if !ManagerAuth(c) {
		input.ApplyWithoutDisable(sql)
	}
	err = sql.Find(&list)
	return len(list), &list, err
}

func findAllRoleSimple(c *fiber.Ctx) (int, *[]model.RoleInfo, error) {
	var list []model.RoleInfo
	sql := base.Engine.Cols(model.IdInRoleInfo, model.RoleNameInRoleInfo)
	err := sql.Where(model.IsDisableInRoleInfo + " = 0").
		Where(model.IsSuperInRoleInfo + " = 0").
		Where(model.RoleFromInRoleInfo + " = 0").Find(&list)
	return len(list), &list, err
}

func findRoleSimpleByIds(c *fiber.Ctx, ids []uint64) (*[]model.RoleInfo, error) {
	var list []model.RoleInfo
	sql := base.Engine.Cols(model.IdInRoleInfo, model.RoleNameInRoleInfo, model.IsSuperInRoleInfo, model.IsManageInRoleInfo)
	err := sql.Where(model.IsDisableInRoleInfo+" = 0").
		In(model.IdInRoleInfo, ids).Find(&list)
	return &list, err
}

func insertRole(c *fiber.Ctx, session *xorm.Session, info *model.RoleInfo) (int64, error) {
	if info.Id > 0 {
		return session.Omit(model.CreatedAt, model.UpdatedAt).
			InsertOne(info)
	} else {
		return session.Omit(model.IdInRoleInfo, model.CreatedAt, model.UpdatedAt).
			InsertOne(info)
	}
}

func updateRoleStatusById(c *fiber.Ctx, session *xorm.Session, id uint64, isDisable int) (int64, error) {
	res, err := session.Exec("update "+model.RoleInfoTable+" set "+model.IsDisableInRoleInfo+" = ? "+
		"where "+model.IdInRoleInfo+" = ?", isDisable, id)
	if nil != err {
		return 0, err
	}
	return res.RowsAffected()
}
