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
	"github.com/gofiber/fiber/v2"
	"strconv"
	"xorm.io/xorm"
)

func ListSnaps(c *fiber.Ctx) error {
	get, id, err := checkSnapType(c)
	if err != nil {
		return err
	}
	switch get {
	case "deploy":
		return listSnapDeploy(c, id)
	case "config":
		return listSnapConfig(c, id)
	}
	return util.ErrorInput(c, "snap type not support")
}

func RollSnaps(c *fiber.Ctx) error {
	get, id, err := checkSnapType(c)
	if err != nil {
		return err
	}
	switch get {
	case "deploy":
		return rollSnapDeploy(c, id)
	case "config":
		return rollSnapConfig(c, id)
	}
	return util.ErrorInput(c, "snap type not support")
}

func checkSnapType(c *fiber.Ctx) (string, uint64, error) {
	get := c.Params("type")
	if "" == get {
		return "", 0, util.ErrorInput(c, "type input error")
	}
	snap := c.Params("snap")
	if "" == snap {
		return "", 0, util.ErrorInput(c, "snap input error")
	}
	id, err := strconv.ParseUint(snap, 10, 64)
	if nil != err {
		return "", 0, util.ErrorInternal(c, err)
	}
	return get, id, nil
}

func listSnapDeploy(c *fiber.Ctx, id uint64) error {
	list, err := findDeploySnapByMain(id)
	return util.ResultList(c, err, 10, list)
}

func listSnapConfig(c *fiber.Ctx, id uint64) error {
	list, err := findConfigSnapByMain(id)
	return util.ResultList(c, err, 10, list)
}

func rollSnapDeploy(c *fiber.Ctx, id uint64) error {
	info, err := getDeploySnapById(id)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	if info.Id <= 0 {
		return util.ErrorInput(c, "input snap not exist")
	}

	spaceBase, err := getSpaceById(info.SpaceId)
	if nil != err {
		return util.ErrorInternal(c, err)
	}

	_, err = base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		return updateDeployTagById(session, info.DeployId, info.DeployTag)
	})
	if nil != err {
		return util.ErrorInternal(c, err)
	}

	status, err := handleWithKind(info.EnvId, info.DeployKind,
		spaceBase.SpaceKeep, info.DeployRender, info.DeployName)
	return util.ResultParamMapOne(c, err, "status", status)
}

func rollSnapConfig(c *fiber.Ctx, id uint64) error {
	info, err := getConfigSnapById(id)
	if nil != err {
		return util.ErrorInternal(c, err)
	}
	if info.Id <= 0 {
		return util.ErrorInput(c, "input snap not exist")
	}

	return updateProp(c, info.PropId, info.UserId, &model.PropertyFile{
		FileContent: info.FileContent,
		CreatedAt:   info.CreatedAt,
	}, true)
}
