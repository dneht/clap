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
	"cana.io/clap/pkg/auth"
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"xorm.io/xorm"
)

const LoginApi = "/login"
const UserApiPre = "/api/user"
const RoleApiPre = "/api/role"

func generateToken() string {
	return auth.MakeSMD5Hash([]byte(strconv.FormatUint(util.UniqueId(), 10)))
}

func generatePasswd(passwd string) string {
	return auth.MakeSMD5Hash([]byte(passwd))
}

func LoginUser(c *fiber.Ctx) error {
	login := new(refer.LoginUserInput)
	err := c.BodyParser(login)
	if nil != err || "" == login.UserName || "" == login.Password {
		return util.ErrorInput(c, "input error")
	}

	info, err := getUserByName(login.UserName)
	if nil != err || info.Id <= 0 {
		return util.ErrorInput(c, "user not found")
	}
	if info.IsDisable != 0 {
		return util.ErrorInput(c, "user is disabled")
	}
	if !auth.CheckSMD5Hash(info.Password, login.Password) {
		return util.ErrorInput(c, "password incorrect")
	}
	cleanUser(info)
	newToken := generateToken()
	_, err = updateUserTokenById(c, base.Engine.NewSession(), info.Id, newToken)
	return util.ResultParam(c, err, true, newToken)
}

func GetUser(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "user id must be set")
	}
	info, err := getUserById(id)
	return util.ResultParam(c, err, true, info)
}

func ListUser(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputLog(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countUserWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findUserWithPage(c, input)
		})
}

func checkUserInput(c *fiber.Ctx) (*model.UserInfo, error) {
	info := new(model.UserInfo)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputLog(c, err, "input is error")
	}
	return info, nil
}

func CreateUser(c *fiber.Ctx) error {
	info, err := checkUserInput(c)
	if err != nil {
		return util.ErrorInput(c, "user body check error")
	}
	if "" == info.UserName || "" == info.Nickname || "" == info.Password {
		return util.ErrorInput(c, "input is error")
	}

	info.Password = generatePasswd(info.Password)
	info.AccessToken = generateToken()
	if "" == info.RoleList {
		info.RoleList = "[]"
	} else {
		roleList := make([]uint64, 0)
		err = json.Unmarshal([]byte(info.RoleList), &roleList)
		if nil != err {
			info.RoleList = "[]"
		}
	}
	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		result, err := insertUser(c, session, info)
		if nil != err {
			return nil, util.ErrorInternal(c, err)
		}

		role := &model.RoleInfo{
			Id:         info.Id,
			RoleName:   "user:" + info.UserName,
			RoleFrom:   1,
			RoleRemark: "用户同名角色",
		}
		result, err = insertRole(c, session, role)
		if nil != err {
			return nil, util.ErrorInternal(c, err)
		}

		result, err = updateUserRoleById(c, session, info.Id, info.RoleList, []uint64{info.Id})
		return nil, util.ResultParamWithMessage(c, err, result > 0, "create user error", info.Id)
	})
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return c.SendString("ok")
}

func ChangeUser(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "user id must be set")
	}
	info, err := checkUserInput(c)
	if err != nil {
		return util.ErrorInput(c, "user body check error")
	}
	cleanUser(info)
	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		if info.IsDisable != 0 || info.Password != "" {
			isDisable := 0
			if info.IsDisable != 0 {
				isDisable = 1
			}
			_, err = updateUserStatusById(c, session, id, generatePasswd(info.Password), isDisable)
			if nil != err {
				return 0, err
			}
			_, err = updateRoleStatusById(c, session, id, isDisable)
			if nil != err {
				return 0, err
			}
		}
		if info.RoleList != "" && info.RoleList != "[]" {
			roleList := make([]uint64, 0)
			err = json.Unmarshal([]byte(info.RoleList), &roleList)
			if nil != err {
				return 0, err
			}
			return updateUserRoleById(c, session, id, "", roleList)
		}
		return 1, nil
	})
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return c.SendString("ok")
}

func cleanUser(user *model.UserInfo) {
	if nil != user {
		invalidUserById(user.Id)
		delete(authInfoMap, user.AccessToken)
	}
}

func GetRole(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "app id must be set")
	}
	info, err := getRoleById(id)
	return util.ResultParam(c, err, true, info)
}

func ListRole(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputLog(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countRoleWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findRoleWithPage(c, input)
		})
}

func SimpleRole(c *fiber.Ctx) error {
	size, list, err := findAllRoleSimple(c)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	return util.ResultList(c, err, size, list)
}

func checkRoleInput(c *fiber.Ctx) (*model.RoleInfo, error) {
	info := new(model.RoleInfo)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputLog(c, err, "input is error")
	}
	if "" == info.RoleName {
		return nil, util.ErrorInput(c, "input is error")
	}
	return info, nil
}

func CreateRole(c *fiber.Ctx) error {
	info, err := checkRoleInput(c)
	if err != nil {
		return err
	}

	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		result, err := insertRole(c, session, info)
		return nil, util.ResultParamWithMessage(c, err, result > 0, "create role error", info.Id)
	})
	return err
}
