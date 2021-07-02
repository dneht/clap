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

var resourceMap = make(map[uint64]*model.Resource)
var permissionMap = make(map[uint64]*model.Permission)

func getResourceById(id uint64) (*model.Resource, error) {
	value, ok := resourceMap[id]
	if ok {
		return value, nil
	}

	var info model.Resource
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		resourceMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found resource")
	}
	return &info, err
}

func countResourceWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.Resource
	sql := base.Engine.Omit(model.CreatedAtInResource, model.UpdatedAtInResource)
	return input.Apply(sql).Count(&info)
}

func findResourceWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.Resource, error) {
	var list []model.Resource
	sql := base.Engine.Omit(model.CreatedAtInResource, model.UpdatedAtInResource)
	err := input.Apply(sql).Find(&list)
	return len(list), &list, err
}

func getPermissionById(id uint64) (*model.Permission, error) {
	value, ok := permissionMap[id]
	if ok {
		return value, nil
	}

	var info model.Permission
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		permissionMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found permission")
	}
	return &info, err
}

func countPermissionWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.Permission
	sql := base.Engine.Omit(model.CreatedAtInPermission, model.UpdatedAtInPermission)
	return input.Apply(sql).Count(&info)
}

func findPermissionWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.Permission, error) {
	var list []model.Permission
	sql := base.Engine.Omit(model.CreatedAtInPermission, model.UpdatedAtInPermission)
	err := input.Apply(sql).Find(&list)
	return len(list), &list, err
}

func findPermissionByRole(c *fiber.Ctx, input *util.MainInput, role []uint64) (int, *[]model.Permission, error) {
	var list []model.Permission
	sql := base.Engine.Omit(model.CreatedAtInPermission, model.UpdatedAtInPermission).
		In(model.RoleIdInPermission, role)
	err := input.Apply(sql).Find(&list)
	return len(list), &list, err
}
