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
	sql := base.Engine.Omit(model.CreatedAtInRoleInfo, model.UpdatedAtInRoleInfo)
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findRoleWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.RoleInfo, error) {
	var list []model.RoleInfo
	sql := base.Engine.Omit(model.CreatedAtInRoleInfo, model.UpdatedAtInRoleInfo)
	err := input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), &list, err
}
