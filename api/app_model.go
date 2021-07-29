package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

func getAppById(id uint64) (*model.Project, error) {
	value, ok := appMap[id]
	if ok {
		return value, nil
	}

	var info model.Project
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		appMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not get app by id")
	}
	return &info, err
}

func getBaseAppInfoById(id uint64) (*model.Project, *refer.AppInfo, error) {
	value, ok := appInfoMap[id]
	if ok {
		return appMap[id], value, nil
	}
	project, err := getAppById(id)
	if nil != err {
		return nil, nil, err
	}

	var info refer.AppInfo
	err = json.Unmarshal([]byte(project.AppInfo), &info)
	if nil != err {
		return nil, nil, err
	}
	appInfoMap[id] = &info
	return project, &info, nil
}

func invalidAppById(id uint64) {
	delete(appMap, id)
	delete(appInfoMap, id)
}

func invalidAppInfoById(id uint64) {
	delete(appInfoMap, id)
}

func findAllAppSimple(c *fiber.Ctx) (int, *[]model.Project, error) {
	var list []model.Project
	sql := base.Engine.Cols(model.IdInProject, model.AppKeyInProject, model.AppTypeInProject).
		Where("is_disable = 0")
	err := SelectAuth(c, model.ProjectTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = sql.Find(&list)
	return len(list), &list, err
}

func countAppWithPage(c *fiber.Ctx, input *util.MainInput) (int64, error) {
	var info model.Project
	sql := base.Engine.Cols(model.IdInProject)
	err := SelectAuth(c, model.ProjectTable, sql)
	if nil != err {
		return 0, err
	}
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findAppWithPage(c *fiber.Ctx, input *util.MainInput) (int, *[]model.Project, error) {
	var list []model.Project
	sql := base.Engine.Omit(model.AppDescInProject, model.AppInfoInProject, model.SourceInfoInProject, model.InjectInfoInProject)
	err := SelectAuth(c, model.ProjectTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), &list, err
}

func updateAppById(c *fiber.Ctx, session *xorm.Session, info *model.Project) (int64, error) {
	if nil == info || info.Id <= 0 {
		return -1, errors.New("input model error, id is empty")
	}
	return session.Omit(model.IdInProject, model.AppKeyInProject, model.AppNameInProject).Update(info)
}

func insertApp(c *fiber.Ctx, session *xorm.Session, info *model.Project) (int64, error) {
	return session.InsertOne(info)
}
