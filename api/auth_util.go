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
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strings"
	"xorm.io/xorm"
)

var (
	backAuth = map[string]*[]string{
		"env": {
			"id in (1, 2, 8, 9)",
		},
		"space": {
			"id in (1000001, 1000011, 1000021, 1000031, 1000041, 2000001, 2000011, 2000021, 8000001, 9000001)",
		},
		"app": {
			"id > 10000",
			"id < 1000000",
		},
		"deploy": {
			"space_id in (1000001, 1000011, 1000021, 1000031, 1000041, 2000001, 2000011, 2000021, 8000001, 9000001)",
			"app_id > 10000",
			"app_id < 1000000",
		},
	}
	frontAuth = map[string]*[]string{
		"env": {
			"id in (1, 2, 8, 9)",
		},
		"space": {
			"id in (1000002, 1000012, 1000022, 1000032, 1000042, 2000002, 2000012, 2000022, 8000002, 9000002)",
		},
		"app": {
			"id >= 3000000",
			"id < 4000000",
		},
		"deploy": {
			"space_id in (1000002, 1000012, 1000022, 1000032, 1000042, 2000002, 2000012, 2000022, 8000002, 9000002)",
			"app_id >= 3000000",
			"app_id < 4000000",
		},
	}
)

const (
	TokenHeader = "X-Clap-Access-Token"
	DataBasePre = "clap:mysql:where"
)

var authInfoMap = make(map[string]*refer.AuthInfo)

func CheckAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if "" != token {
		token = token[7:]
	}
	if "" == token {
		return fiber.ErrUnauthorized
	}

	_, ok := authInfoMap[token]
	c.Request().Header.Set(TokenHeader, token)
	if ok {
		return c.Next()
	} else {
		err := getAllAuthInfo(c, token)
		if nil != err {
			return err
		}
		return c.Next()
	}
}

func getAllAuthInfo(c *fiber.Ctx, token string) error {
	user, err := getUserByToken(token)
	if nil != err {
		return err
	}
	input := &util.MainInput{}
	roleIds, err := getUserRoleIds(user)
	if nil != err {
		return err
	}
	if len(roleIds) == 0 {
		return fiber.ErrForbidden
	}
	authInfo := refer.AuthInfo{}
	input.Ids = roleIds
	_, roleInfos, err := findRoleWithPage(c, input)
	if nil != err {
		return err
	}
	for _, roleInfo := range *roleInfos {
		if roleInfo.IsSuper == 1 {
			authInfo.IsSuper = true
			authInfo.IsManage = true
			authInfoMap[token] = &authInfo
			return nil
		}
		if roleInfo.IsManage == 1 {
			authInfo.IsManage = true
		}
	}
	authResMap := make(map[uint64]*refer.AuthResInfo)
	input.Ids = make([]uint64, 0)
	_, powerInfos, err := findPermissionByRole(c, input, roleIds)
	if nil != err {
		return err
	}
	if nil == powerInfos {
		return fiber.ErrForbidden
	}
	for _, powerInfo := range *powerInfos {
		powerList, err := getUserPowerInfo(&powerInfo)
		if nil != err {
			return err
		}
		resInfo, err := getResourceById(powerInfo.ResId)
		if nil != err {
			return err
		}
		authResOne, ok := authResMap[resInfo.Id]
		if ok {
			*authResOne.PowerInfo = append(*authResOne.PowerInfo, refer.AuthPowerInfo{
				PowerData: powerInfo.ResPower,
				PowerList: &powerList,
			})
		} else {
			authResOne = &refer.AuthResInfo{
				ResName:   resInfo.ResName,
				ResExtend: resInfo.ResInfo,
				PowerInfo: &[]refer.AuthPowerInfo{
					{
						PowerData: powerInfo.ResPower,
						PowerList: &powerList,
					},
				},
			}
		}
		authResMap[resInfo.Id] = authResOne
	}
	authResInfo := make([]refer.AuthResInfo, 0, len(*roleInfos)*2)
	for _, authResOne := range authResMap {
		authResInfo = append(authResInfo, *authResOne)
	}
	authInfo.ResInfo = &authResInfo
	authInfoMap[token] = &authInfo
	return nil
}

func getUserRoleIds(user *model.UserInfo) ([]uint64, error) {
	var arr []uint64
	role := user.RoleList
	if "" != role && "[]" != role {
		err := json.Unmarshal([]byte(role), &arr)
		if nil != err {
			return nil, err
		}
	} else {
		arr = make([]uint64, 0)
	}
	return arr, nil
}

func getUserPowerInfo(power *model.Permission) ([]string, error) {
	var arr []string
	info := power.PowerInfo
	if "" != info && "[]" != info {
		err := json.Unmarshal([]byte(info), &arr)
		if nil != err {
			return nil, err
		}
	} else {
		arr = make([]string, 0)
	}
	return arr, nil
}

func DatabaseAuth(t string, c *fiber.Ctx, sql *xorm.Session) error {
	token := c.Get(TokenHeader)
	if "" == token {
		return fiber.ErrUnauthorized
	}
	auth, ok := authInfoMap[token]
	if !ok {
		return fiber.ErrForbidden
	}

	if auth.IsSuper {
		return nil
	}
	if nil == auth.ResInfo {
		return fiber.ErrForbidden
	}
	for _, res := range *auth.ResInfo {
		if res.ResName == DataBasePre {
			return nil
		} else if res.ResName == DataBasePre+":"+t {
			if nil == res.PowerInfo {
				return fiber.ErrForbidden
			}
			for _, pow := range *res.PowerInfo {
				if !checkDBPower(c, pow.PowerData) || nil == pow.PowerList {
					return fiber.ErrForbidden
				}
				for _, val := range *pow.PowerList {
					if "" != val {
						sql.Where(val)
					}
				}
			}
			return nil
		}
	}
	return nil
}

func checkDBPower(c *fiber.Ctx, p uint) bool {
	if p&(1<<5) > 0 {
		return true
	}
	return false
}

func checkURIPower(c *fiber.Ctx, p uint) bool {
	if p&(1<<4) > 0 {
		return true
	}
	method := strings.ToUpper(c.Method())
	if method == "GET" {
		return p&1 > 0
	} else if method == "POST" {
		if strings.HasSuffix(c.Path(), "list") {
			return p&1 > 0
		} else {
			return p&(1<<1) > 0
		}
	} else if method == "PUT" {
		return p&(1<<2) > 0
	} else if method == "DELETE" {
		return p&(1<<3) > 0
	}
	return false
}
