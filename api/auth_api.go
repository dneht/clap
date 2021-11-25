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
	"cana.io/clap/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

const StaticApi = "/api/static"
const ResApiPre = "/api/res"
const PowApiPre = "/api/pow"

func GetRes(c *fiber.Ctx) error {
	if !ManagerAuth(c) {
		return fiber.ErrForbidden
	}
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "user id must be set")
	}
	info, err := getResourceById(id)
	return util.ResultParam(c, err, true, info)
}

func ListRes(c *fiber.Ctx) error {
	if !ManagerAuth(c) {
		return fiber.ErrForbidden
	}
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputLog(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countResourceWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findResourceWithPage(c, input)
		})
}

func SimpleRes(c *fiber.Ctx) error {
	list, _ := base.Resources()
	result := make([]map[string]interface{}, 0, len(*list))
	for _, one := range *list {
		result = append(result, map[string]interface{}{
			"id":   one.Id,
			"name": one.ResName,
			"pow":  one.ResInfo,
		})
	}
	return util.ResultList(c, nil, len(result), result)
}

func StaticRes(c *fiber.Ctx) error {
	auth, err := getAuthFromHeader(c)
	if nil != err {
		return err
	}
	return util.ResultParam(c, err, true, auth.ResInfo)
}

func GetPow(c *fiber.Ctx) error {
	if !ManagerAuth(c) {
		return fiber.ErrForbidden
	}
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "app id must be set")
	}
	info, err := getPermissionById(id)
	return util.ResultParam(c, err, true, info)
}

func ListPow(c *fiber.Ctx) error {
	if !ManagerAuth(c) {
		return fiber.ErrForbidden
	}
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputLog(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countPermissionWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findPermissionWithPage(c, input)
		})
}

func SimplePow(c *fiber.Ctx) error {
	get := c.Query("type")
	if "" == get {
		return util.ErrorInput(c, "type input error")
	}
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputLog(c, err, "main input error")
	}
	if len(param.Ids) == 0 {
		return util.ResultList(c, nil, 0, map[string]interface{}{})
	}
	token, err := getInputToken(c)
	if nil != err {
		return err
	}
	auth, err := getAuthFromToken(c, token)
	if nil != err {
		return err
	}
	if auth.IsManage {
		return util.ResultParam(c, nil, true, map[string]interface{}{
			util.GenerateMD5(get, token): map[string]interface{}{},
		})
	}

	idmap := make(map[uint64]bool, len(param.Ids))
	for _, id := range param.Ids {
		idmap[id] = true
	}
	list := (*auth.ResPower)[CommonPre+get]
	if nil == list {
		return util.ErrorInput(c, "type not found")
	}
	merge := make(map[string]uint, len(list))
	for _, one := range list {
		_, ok := idmap[one.LinkId]
		if ok {
			key := util.GenerateMD5(get, token, strconv.FormatUint(one.LinkId, 10))
			getPower, pok := merge[key]
			if pok {
				merge[key] = one.Power | getPower
			} else {
				merge[key] = one.Power
			}
		}
	}
	result := make(map[string]interface{})
	for key, power := range merge {
		result[key] = map[string]bool{
			"thisEdit":     power&AllowThisUpdate > 0,
			"thisPack":     power&AllowThisPackageDeploy > 0,
			"thisRollback": power&AllowThisRollbackDeploy > 0,
			"docView":      power&AllowThisDocumentView > 0,
			"propView":     power&AllowThisPropertyView > 0,
			"propEdit":     power&AllowThisPropertyUpdate > 0,
			"propAdd":      power&AllowThisPropertyCreate > 0,
			"podLog":       power&AllowThisPodLog > 0,
			"podExec":      power&AllowThisPodExec > 0,
			"podRestart":   power&AllowThisPodRestart > 0,
			"podSpace":     power&AllowThisPodSpace > 0,
		}
	}
	return util.ResultParam(c, nil, true, result)
}
