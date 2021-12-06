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
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
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
	if nil != err || !result {
		return nil, errors.New("can not found user info")
	}
	if nil == err {
		userMap[id] = &info
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
	if nil != err || !result {
		return nil, errors.New("can not found user info")
	}
	if nil == err {
		userMap[info.Id] = &info
		userTokenMap[token] = info.Id
	}
	return &info, err
}

func getUserByName(name string) (*model.UserInfo, error) {
	var info model.UserInfo
	result, err := base.Engine.Where(model.UserNameInUserInfo+" = ?", name).Get(&info)
	if nil != err || !result {
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
	sql := base.Engine.Cols()
	err := SelectAuth(c, model.UserInfoTable, sql)
	if nil != err {
		return 0, err
	}
	if !ManagerAuth(c) {
		input.ApplyWithoutDisable(sql)
	}
	return sql.Count(&info)
}

func findUserWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.UserInfo, error) {
	var list []model.UserInfo
	sql := base.Engine.Omit(model.CreatedAt, model.UpdatedAt)
	err := SelectAuth(c, model.UserInfoTable, sql)
	if nil != err {
		return 0, nil, err
	}
	if !ManagerAuth(c) {
		input.ApplyWithoutDisable(sql)
	}
	err = sql.Find(&list)
	return len(list), &list, err
}

func insertUser(c *fiber.Ctx, session *xorm.Session, info *model.UserInfo) (int64, error) {
	return session.Omit(model.IdInUserInfo, model.CreatedAt, model.UpdatedAt).
		InsertOne(info)
}

func updateUserTokenById(c *fiber.Ctx, session *xorm.Session, id uint64, newToken string) (int64, error) {
	res, err := session.Exec("update "+model.UserInfoTable+" set "+model.AccessTokenInUserInfo+" = ? "+
		"where "+model.IdInUserInfo+" = ?", newToken, id)
	if nil != err {
		return 0, err
	}
	return res.RowsAffected()
}

func updateUserStatusById(c *fiber.Ctx, session *xorm.Session, id uint64, newPasswd string, isDisable int) (int64, error) {
	res, err := session.Exec("update "+model.UserInfoTable+" set "+model.PasswordInUserInfo+" = ?, "+model.IsDisableInUserInfo+" = ? "+
		"where "+model.IdInUserInfo+" = ?", newPasswd, isDisable, id)
	if nil != err {
		return 0, err
	}
	return res.RowsAffected()
}

func updateUserRoleById(c *fiber.Ctx, session *xorm.Session, id uint64, base string, add []uint64) (int64, error) {
	list := make([]uint64, 0, 10)
	if "" != base && "[]" != base {
		err := json.Unmarshal([]byte(base), &list)
		if nil != err {
			list = add
		} else {
			for _, one := range add {
				list = append(list, one)
			}
		}
	} else {
		list = add
	}
	res, err := json.Marshal(list)
	if nil != err {
		res = []byte("[]")
	}
	exec, err := session.Exec("update "+model.UserInfoTable+" set "+model.RoleListInUserInfo+" = ? "+
		"where "+model.IdInRoleInfo+" = ?",
		string(res), id)
	if nil != err {
		return 0, err
	}
	return exec.RowsAffected()
}
