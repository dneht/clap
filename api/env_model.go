package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func getEnvById(id uint64) (*model.Environment, error) {
	return base.Env(id)
}

func findAllEnvSimple(c *fiber.Ctx) (int, []model.Environment, error) {
	var list []model.Environment
	sql := base.Engine.Cols(model.IdInEnvironment, model.EnvInEnvironment, model.EnvNameInEnvironment).
		Where(model.IsDisableInEnvironment + " = 0")
	err := SelectAuth(c, model.EnvironmentTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = sql.OrderBy(model.IdInEnvironmentSpace).Find(&list)
	return len(list), list, err
}

func countEnvWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.Environment
	sql := base.Engine.Cols()
	err := SelectAuth(c, model.EnvironmentTable, sql)
	if nil != err {
		return 0, err
	}
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findEnvWithPage(c *fiber.Ctx, input *util.MainInput) (int, []model.Environment, error) {
	var list []model.Environment
	sql := base.Engine.Omit(model.SyncInfoInEnvironment, model.DeployInfoInEnvironment, model.FormatInfoInEnvironment,
		model.CreatedAt, model.UpdatedAt)
	err := SelectAuth(c, model.EnvironmentTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), list, err
}

func getSpaceById(id uint64) (*model.EnvironmentSpace, error) {
	value, ok := spaceMap[id]
	if ok {
		return value, nil
	}

	var info model.EnvironmentSpace
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		spaceMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not get space by id")
	}
	return &info, err
}

func getBaseSpaceInfoById(id uint64) (*model.EnvironmentSpace, *refer.SpaceInfo, error) {
	value, ok := spaceInfoMap[id]
	if ok {
		return spaceMap[id], value, nil
	}
	space, err := getSpaceById(id)
	if nil != err {
		return nil, nil, err
	}

	var info refer.SpaceInfo
	err = json.Unmarshal([]byte(space.SpaceInfo), &info)
	if nil != err {
		return nil, nil, err
	}
	spaceInfoMap[id] = &info
	return space, &info, nil
}

func invalidSpaceById(id uint64) {
	delete(spaceMap, id)
	delete(spaceInfoMap, id)
}

func findAllSpaceSimple(c *fiber.Ctx, envId uint64) (int, []model.EnvironmentSpace, error) {
	var list []model.EnvironmentSpace
	sql := base.Engine.Cols(model.IdInEnvironmentSpace, model.EnvIdInEnvironmentSpace, model.SpaceNameInEnvironmentSpace).
		Where(model.EnvIdInEnvironmentSpace+" = ?", envId).And(model.IsViewInEnvironmentSpace + " = 0").And(model.IsDisableInEnvironmentSpace + " = 0")
	err := SelectAuth(c, model.EnvironmentSpaceTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = sql.OrderBy(model.IdInEnvironmentSpace).Find(&list)
	return len(list), list, err
}

func countSpaceWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.EnvironmentSpace
	sql := base.Engine.Cols()
	err := SelectAuth(c, model.EnvironmentSpaceTable, sql)
	if nil != err {
		return 0, err
	}
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findSpaceWithPage(c *fiber.Ctx, input *util.MainInput) (int, []model.EnvironmentSpace, error) {
	var list []model.EnvironmentSpace
	sql := base.Engine.Omit(model.SpaceInfoInEnvironmentSpace, model.CreatedAt, model.UpdatedAt)
	err := SelectAuth(c, model.EnvironmentSpaceTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), list, err
}
