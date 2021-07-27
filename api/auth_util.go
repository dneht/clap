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
	"cana.io/clap/util"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"log"
	"strings"
	"xorm.io/xorm"
)

const (
	TokenHeader               = "X-Clap-Access-Token"
	MenuPre                   = "menu:"
	CommonPre                 = "base:"
	CreatePre                 = "add:"
	ListSuffix                = "/list"
	SimpleSuffix              = "/simple"
	AllowThisSelect           = 1
	AllowThisUpdate           = 1 << 1
	AllowThisCreate           = 1 << 2
	AllowThisDelete           = 1 << 3
	AllowThisGrant            = 1 << 8
	AllowThisAndSubGrant      = 1 << 9
	AllowThisDocumentView     = 1 << 12
	AllowThisPackageDeploy    = 1 << 13
	AllowThisPropertyView     = 1 << 16
	AllowThisPropertyCreate   = 1 << 17
	AllowThisPropertyUpdate   = 1 << 18
	AllowThisDeployPlanCreate = 1 << 20
	AllowThisPodLog           = 1 << 22
	AllowThisPodExec          = 1 << 23
	AllowThisPodRestart       = 1 << 24
	AllowThisPodRollback      = 1 << 25
	AllowThisPodSpace         = 1 << 28
)

func CheckAuth(c *fiber.Ctx) error {
	token, err := getInputToken(c)
	if nil != err {
		return err
	}
	c.Request().Header.Set(TokenHeader, token)

	_, ok := authInfoMap[token]
	if !ok {
		err = getAllAuthInfo(c, token)
		if nil != err {
			return err
		}
	}
	err = RequestAuth(c)
	if nil != err {
		return err
	}
	return c.Next()
}

func getInputToken(c *fiber.Ctx) (string, error) {
	token := c.Get("Authorization")
	if len(token) >= 30 {
		token = token[7:]
	}
	if len(token) <= 22 {
		return "", fiber.ErrUnauthorized
	}
	return utils.CopyString(token), nil
}

func getAllAuthInfo(c *fiber.Ctx, token string) error {
	user, err := getUserByToken(token)
	if nil != err {
		log.Print(err)
		return fiber.ErrUnauthorized
	}
	roleIds, err := getUserRoleIds(user)
	if nil != err {
		log.Print(err)
		return fiber.ErrUnauthorized
	}
	if len(roleIds) == 0 {
		return fiber.ErrForbidden
	}
	authInfo := refer.AuthInfo{
		UserId: user.Id,
	}
	roleInfos, err := findRoleSimpleByIds(c, roleIds)
	if nil != err {
		return err
	}
	for _, roleInfo := range *roleInfos {
		if roleInfo.IsSuper == 1 {
			authInfo.IsSuper = true
			authInfo.IsManage = true
		}
		if roleInfo.IsManage == 1 {
			authInfo.IsManage = true
		}
		if authInfo.IsSuper || authInfo.IsManage {
			authInfo.ResInfo = getAdminResInfo()
			authInfoMap[token] = &authInfo
			return nil
		}
	}
	authPower := make(map[string][]refer.AuthPower)
	powerInfos, err := findPermissionByRole(c, roleIds)
	if nil != err {
		log.Print(err)
		return fiber.ErrUnauthorized
	}
	if nil == powerInfos {
		return fiber.ErrForbidden
	}
	resIds := make([]uint64, 0, len(*powerInfos))
	for _, powerInfo := range *powerInfos {
		resInfo, err := getResourceById(powerInfo.ResId)
		if nil != err {
			return err
		}
		resIds = append(resIds, resInfo.Id)
		powerOne, ok := authPower[resInfo.ResName]
		if powerInfo.ResPower <= 0 {
			continue
		}
		if ok {
			powerOne = append(powerOne, refer.AuthPower{
				Power:  powerInfo.ResPower,
				RoleId: powerInfo.RoleId,
				LinkId: powerInfo.LinkId,
			})
		} else {
			powerOne = []refer.AuthPower{
				{
					Power:  powerInfo.ResPower,
					RoleId: powerInfo.RoleId,
					LinkId: powerInfo.LinkId,
				},
			}
		}
		authPower[resInfo.ResName] = powerOne
	}
	authInfo.ResInfo = getUserResInfo(resIds)
	authInfo.ResPower = &authPower
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

func getAdminResInfo() *map[string]interface{} {
	allInfo, allRes := base.Resources()
	resOrder := make(map[string]int)
	resInfo := make(map[string]interface{})
	for id, oneRes := range *allRes {
		oneInfo, ok := (*allInfo)[id]
		if ok {
			for key, value := range *oneRes {
				order, sok := resOrder[key]
				if !sok || oneInfo.ResOrder > order {
					resOrder[key] = oneInfo.ResOrder
					resInfo[key] = value
				}
			}
		}
	}
	return &resInfo
}

func getUserResInfo(ids []uint64) *map[string]interface{} {
	resInfo := make(map[string]interface{})
	if len(ids) == 0 {
		return &resInfo
	}

	for _, id := range ids {
		_, oneRes := base.Resource(id)
		for key, value := range *oneRes {
			resInfo[key] = value
		}
	}
	return &resInfo
}

func getAuthFromHeader(c *fiber.Ctx) (*refer.AuthInfo, error) {
	token := c.Get(TokenHeader)
	if "" == token {
		return nil, fiber.ErrUnauthorized
	}
	return getAuthFromToken(token)
}

func getAuthFromToken(token string) (*refer.AuthInfo, error) {
	auth, ok := authInfoMap[token]
	if !ok {
		return nil, fiber.ErrForbidden
	}
	return auth, nil
}

func ManagerAuth(c *fiber.Ctx) bool {
	auth, err := getAuthFromHeader(c)
	if nil != err {
		return false
	}
	if auth.IsManage {
		return true
	}
	return false
}

func RequestAuth(c *fiber.Ctx) error {
	auth, err := getAuthFromHeader(c)
	if nil != err {
		return err
	}
	if auth.IsManage {
		return nil
	}
	if nil == auth.ResPower {
		return fiber.ErrForbidden
	}

	path := c.Path()
	if path == ConfigApi || path == CleanApi || path == StaticApi {
		return nil
	}
	if strings.HasSuffix(path, ListSuffix) || strings.HasSuffix(path, SimpleSuffix) {
		return nil
	}
	switch c.Method() {
	case fiber.MethodOptions:
		return nil
	case fiber.MethodConnect:
		return nil
	case fiber.MethodGet:
		return requestAuth(c, auth.ResPower, path, CommonPre, AllowThisSelect)
	case fiber.MethodPut:
		return requestAuth(c, auth.ResPower, path, CreatePre, AllowThisCreate)
	case fiber.MethodPost:
		return requestAuth(c, auth.ResPower, path, CommonPre, AllowThisUpdate)
	case fiber.MethodDelete:
		return requestAuth(c, auth.ResPower, path, CommonPre, AllowThisDelete)
	}
	return fiber.ErrForbidden
}

func requestAuth(c *fiber.Ctx, pow *map[string][]refer.AuthPower, path, pre string, allow uint) error {
	table, err := getTableByPath(path)
	if nil != err {
		return err
	}

	list, ok := (*pow)[pre+table]
	if ok {
		if pre == CommonPre {
			id, err := util.CheckIdInput(c, "id")
			if nil != err {
				log.Print(err)
				return fiber.ErrForbidden
			}
			for _, one := range list {
				if one.LinkId == id && one.Power&allow > 0 {
					return nil
				}
			}
		} else {
			for _, one := range list {
				if one.Power&allow > 0 {
					return nil
				}
			}
		}
	}
	return fiber.ErrForbidden
}

func getTableByPath(path string) (string, error) {
	if strings.HasPrefix(path, EnvApiPre) {
		return model.EnvironmentTable, nil
	} else if strings.HasPrefix(path, SpaceApiPre) {
		return model.EnvironmentSpaceTable, nil
	} else if strings.HasPrefix(path, DeployApiPre) {
		return model.DeploymentTable, nil
	} else if strings.HasPrefix(path, RenderApiPre) {
		return model.TemplateTable, nil
	} else if strings.HasPrefix(path, UserApiPre) {
		return model.UserInfoTable, nil
	} else if strings.HasPrefix(path, RoleApiPre) {
		return model.RoleInfoTable, nil
	} else if strings.HasPrefix(path, ResApiPre) {
		return model.ResourceTable, nil
	} else if strings.HasPrefix(path, PowApiPre) {
		return model.PermissionTable, nil
	}
	return "", fiber.ErrForbidden
}

func SelectAuth(c *fiber.Ctx, table string, sql *xorm.Session) error {
	auth, err := getAuthFromHeader(c)
	if nil != err {
		return err
	}
	if auth.IsSuper {
		return nil
	}
	if nil == auth.ResPower {
		return fiber.ErrForbidden
	}

	now := CommonPre + table
	list, ok := (*auth.ResPower)[now]
	if ok {
		ids := make([]uint64, 0, len(list))
		hash := make(map[uint64]bool)
		for _, one := range list {
			if one.LinkId > 0 && (one.Power&AllowThisSelect > 0 || one.Power&AllowThisDocumentView > 0) {
				_, get := hash[one.LinkId]
				if !get {
					ids = append(ids, one.LinkId)
					hash[one.LinkId] = true
				}
			}
		}
		if len(ids) == 0 {
			return fiber.ErrForbidden
		}
		sql.In("id", ids)
		return nil
	}
	return fiber.ErrForbidden
}

func DeploymentAuth(c *fiber.Ctx, id uint64, allow uint) error {
	token, err := getInputToken(c)
	if nil != err {
		return err
	}
	auth, err := getAuthFromToken(token)
	if nil != err {
		return err
	}
	return checkInnerAuth(auth, id, allow)
}

func WebsocketAuth(token string, id uint64, allow uint) error {
	auth, err := getAuthFromToken(token)
	if nil != err {
		return err
	}
	return checkInnerAuth(auth, id, allow)
}

func checkInnerAuth(auth *refer.AuthInfo, id uint64, allow uint) error {
	if auth.IsSuper {
		return nil
	}
	if nil == auth.ResPower {
		return fiber.ErrForbidden
	}
	if id <= 0 {
		return fiber.ErrForbidden
	}

	list, ok := (*auth.ResPower)[CommonPre+model.DeploymentTable]
	if ok {
		for _, one := range list {
			if one.LinkId == id && one.Power&allow > 0 {
				return nil
			}
		}
	}
	return fiber.ErrForbidden
}
