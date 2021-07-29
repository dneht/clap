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

func getDeployById(id uint64) (*model.Deployment, error) {
	value, ok := deployMap[id]
	if ok {
		return value, nil
	}

	var info model.Deployment
	result, err := base.Engine.ID(id).Get(&info)
	if nil == err {
		deployMap[id] = &info
	}
	if !result {
		return nil, errors.New("can not found deploy")
	}
	return &info, err
}

func getDeployAppInfoById(id uint64) (*model.Deployment, *refer.AppInfo, error) {
	value, ok := appDeployMap[id]
	if ok {
		return deployMap[id], value, nil
	}
	deploy, err := getDeployById(id)
	if nil != err {
		return nil, nil, err
	}

	var info refer.AppInfo
	err = json.Unmarshal([]byte(deploy.AppInfo), &info)
	if nil != err {
		return nil, nil, err
	}
	appDeployMap[id] = &info
	return deploy, &info, nil
}

func invalidDeployById(id uint64) {
	delete(deployMap, id)
	delete(appDeployMap, id)
}

func invalidDeployInfoById(id uint64) {
	delete(appDeployMap, id)
}

func countDeployWithPage(c *fiber.Ctx, input *util.MainInput, planId uint64) (int64, error) {
	var info model.Deployment
	sql := base.Engine.Cols(model.IdInDeployment).Where(model.PlanIdInDeployment + " = ?", planId)
	err := SelectAuth(c, model.DeploymentTable, sql)
	if nil != err {
		return 0, err
	}
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findDeployWithPage(c *fiber.Ctx, input *util.MainInput, planId uint64) (int, *[]model.Deployment, error) {
	var list []model.Deployment
	sql := base.Engine.Omit(model.AppInfoInDeployment).Where(model.PlanIdInDeployment + " = ?", planId)
	err := SelectAuth(c, model.DeploymentTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), &list, err
}

func updateDeployById(c *fiber.Ctx, session *xorm.Session, info *model.Deployment) (int64, error) {
	if nil == info || info.Id <= 0 {
		return -1, errors.New("input model error, id is empty")
	}
	result, err := session.Omit(model.IdInDeployment, model.AppIdInDeployment, model.EnvIdInDeployment,
		model.SpaceIdInDeployment, model.PlanIdInDeployment, model.DeployTagInDeployment).Update(info)
	if nil == err {
		invalidDeployById(info.Id)
	}
	return result, err
}

func insertDeploy(c *fiber.Ctx, session *xorm.Session, info *model.Deployment) (int64, error) {
	return session.InsertOne(info)
}

func updateDeployStatusById(c *fiber.Ctx, session *xorm.Session, id uint64, status int, tag string) (interface{}, error) {
	var info model.Deployment
	get, err := session.Cols(model.DeployStatusInDeployment).
		ForUpdate().ID(id).Get(&info)
	if nil != err {
		return -1, err
	}
	if !get {
		return 0, errors.New("deploy not exist")
	}
	if info.DeployStatus == status {
		return 0, nil
	}
	info.DeployStatus = status
	var result int64
	if "" == tag {
		result, err = session.Cols(model.DeployStatusInDeployment).
			Where(model.IdInDeployment + " = ?", id).Update(info)
	} else {
		info.DeployTag = tag
		result, err = session.Cols(model.DeployStatusInDeployment, model.DeployTagInDeployment).
			Where(model.IdInDeployment + " = ?", id).Update(info)
	}
	if nil != err {
		return result, err
	}
	return 1, nil
}
