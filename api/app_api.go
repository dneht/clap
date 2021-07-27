package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/util"
	"github.com/gofiber/fiber/v2"
)

const AppApiPre = "/api/app"

func GetApp(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "app id must be set")
	}
	info, err := getAppById(id)
	return util.ResultParam(c, err, true, info)
}

func ListApp(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return util.ErrorInputErrorMessage(c, err, "main input error")
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countAppWithPage(c, input)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findAppWithPage(c, input)
		})
}

func SimpleApp(c *fiber.Ctx) error {
	size, list, err := findAllAppSimple(c)
	if nil != err {
		return util.ErrorInputErrorMessage(c, err, "can not get app list")
	}
	return util.ResultList(c, err, size, list)
}

func CreateApp(c *fiber.Ctx) error {
	info, err := checkProjectInput(c)
	if err != nil {
		return err
	}
	result, err := insertApp(c, base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func UpdateApp(c *fiber.Ctx) error {
	info, err := checkProjectInput(c)
	if err != nil {
		return err
	}
	result, err := updateAppById(c, base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func checkProjectInput(c *fiber.Ctx) (*model.Project, error) {
	info := new(model.Project)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputErrorMessage(c, err, "input is error")
	}
	if info.Id <= 10000 {
		return nil, util.ErrorInputShowMessage(c, "id set error")
	}
	if "" == info.AppKey {
		return nil, util.ErrorInputShowMessage(c, "app key set error")
	}
	if "" == info.AppName {
		return nil, util.ErrorInputShowMessage(c, "app name set error")
	}
	if "" == info.AppDesc {
		return nil, util.ErrorInputShowMessage(c, "app desc set error")
	}
	if info.AppType <= 0 {
		return nil, util.ErrorInputShowMessage(c, "app type set error")
	}
	if "" == info.AppInfo {
		return nil, util.ErrorInputShowMessage(c, "app info set error")
	}
	if "" == info.SourceInfo || "{}" == info.SourceInfo {
		return nil, util.ErrorInputShowMessage(c, "app enc set error")
	}
	return info, nil
}

func DeleteApp(c *fiber.Ctx) error {
	param, err := util.CheckMainInput(c)
	if nil != err {
		return err
	}
	if len(param.Ids) <= 0 {
		return util.ResultEmpty(c)
	}
	//return ResultParamWithMessage(c, , , "update app error")
	return util.ResultEmpty(c)
}
