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
	"cana.io/clap/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ListSnaps(c *fiber.Ctx) error {
	get, id, err := checkSnapType(c)
	if err != nil {
		return err
	}
	switch get {
	case "deployment":

	}
	return nil
}

func RollSnaps(c *fiber.Ctx) error {
	return nil
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
	list, err := findPropDeployByMain(id)
	return util.ResultList(c, err, 10, list)
}

func listSnapProp(c *fiber.Ctx, id uint64) error {
	list, err := findPropSnapByMain(id)
	return util.ResultList(c, err, 10, list)
}
