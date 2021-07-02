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

var userMap = make(map[uint64]*model.UserInfo)
var userTokenMap = make(map[string]uint64)

func getUserById(id uint64) (*model.UserInfo, error) {
	value, ok := userMap[id]
	if ok {
		return value, nil
	}

	var info model.UserInfo
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		userMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found user info")
	}
	return &info, err
}

func getUserByToken(token string) (*model.UserInfo, error) {
	id, ok := userTokenMap[token]
	if ok {
		return getUserById(id)
	}

	var info model.UserInfo
	result, err := base.Engine.Where(model.AccessTokenInUserInfo+" = ?", token).Get(&info)
	if nil == err {
		userMap[id] = &info
		userTokenMap[token] = info.Id
	}
	if !result {
		return nil, errors.New("can not found user info")
	}
	return &info, err
}

func invalidUserById(id uint64) {
	get, ok := userMap[id]
	if ok {
		delete(userMap, id)
		delete(userTokenMap, get.AccessToken)
	}
}

func countUserWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.UserInfo
	sql := base.Engine.Omit(model.CreatedAtInUserInfo, model.UpdatedAtInUserInfo)
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findUserWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.UserInfo, error) {
	var list []model.UserInfo
	sql := base.Engine.Omit(model.CreatedAtInUserInfo, model.UpdatedAtInUserInfo)
	err := input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), &list, err
}
