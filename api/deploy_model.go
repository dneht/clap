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
	sql := base.Engine.Cols().Where(model.PlanIdInDeployment+" = ?", planId)
	err := SelectAuth(c, model.DeploymentTable, sql)
	if nil != err {
		return 0, err
	}
	return input.ApplyWithoutDisable(sql).Count(&info)
}

func findDeployWithPage(c *fiber.Ctx, input *util.MainInput, planId uint64) (int, []model.Deployment, error) {
	var list []model.Deployment
	sql := base.Engine.Omit(model.AppInfoInDeployment, model.CreatedAt, model.UpdatedAt).Where(model.PlanIdInDeployment+" = ?", planId)
	err := SelectAuth(c, model.DeploymentTable, sql)
	if nil != err {
		return 0, nil, err
	}
	err = input.ApplyWithoutDisable(sql).Find(&list)
	return len(list), list, err
}

func updateDeployById(session *xorm.Session, info *model.Deployment) (int64, error) {
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

func insertDeploy(session *xorm.Session, info *model.Deployment) (int64, error) {
	return session.InsertOne(info)
}

func updateDeployStatusById(session *xorm.Session, id uint64, status int, tag string) (interface{}, error) {
	result, info, err := selectDeployWithUpdate(session, id)
	if nil != err {
		return result, err
	}
	if info.DeployStatus == status {
		return 0, nil
	}
	info.DeployStatus = status
	if "" == tag {
		result, err = session.Cols(model.DeployStatusInDeployment).
			Where(model.IdInDeployment+" = ?", id).Update(info)
	} else {
		info.DeployTag = tag
		result, err = session.Cols(model.DeployStatusInDeployment, model.DeployTagInDeployment).
			Where(model.IdInDeployment+" = ?", id).Update(info)
	}
	if nil != err {
		return result, err
	}
	get, err := getDeployById(id)
	if nil != err {
		invalidDeployById(id)
	} else {
		get.DeployStatus = status
		if "" != tag {
			get.DeployTag = tag
		}
	}
	return 1, nil
}

func updateDeployBranchById(session *xorm.Session, id uint64, branch string) (interface{}, error) {
	result, info, err := selectDeployWithUpdate(session, id)
	if nil != err {
		return result, err
	}
	info.BranchName = branch
	result, err = session.Cols(model.BranchNameInDeployment).
		Where(model.IdInDeployment+" = ?", id).Update(info)
	if nil != err {
		return result, err
	}
	get, err := getDeployById(id)
	if nil != err {
		invalidDeployById(id)
	} else {
		get.BranchName = branch
	}
	return 1, nil
}

func updateDeployTagById(session *xorm.Session, id uint64, tag string) (interface{}, error) {
	result, info, err := selectDeployWithUpdate(session, id)
	if nil != err {
		return result, err
	}
	info.DeployTag = tag
	result, err = session.Cols(model.DeployTagInDeployment).
		Where(model.IdInDeployment+" = ?", id).Update(info)
	if nil != err {
		return result, err
	}
	get, err := getDeployById(id)
	if nil != err {
		invalidDeployById(id)
	} else {
		get.DeployTag = tag
	}
	return 1, nil
}

func selectDeployWithUpdate(session *xorm.Session, id uint64) (int64, *model.Deployment, error) {
	var info model.Deployment
	get, err := session.Cols(model.DeployStatusInDeployment).
		ForUpdate().ID(id).Get(&info)
	if nil != err {
		return -1, nil, err
	}
	if !get {
		return 0, nil, errors.New("deploy not exist")
	}
	return 0, &info, nil
}

func getDeploySnapById(id uint64) (*model.DeploymentSnap, error) {
	var info model.DeploymentSnap
	_, err := base.Engine.Omit(model.BranchNameInDeploymentSnap).
		ID(id).Get(&info)
	return &info, err
}

func getLatestDeploySnapByMain(id uint64) (*model.DeploymentSnap, error) {
	var info model.DeploymentSnap
	err := base.Engine.Omit(model.AppIdInDeploymentSnap, model.EnvIdInDeploymentSnap, model.SpaceIdInDeploymentSnap, model.DeployRenderInDeploymentSnap).
		Where(model.DeployIdInDeploymentSnap+" = ?", id).
		Desc(model.IdInDeploymentSnap).Limit(1).
		Find(&info)
	return &info, err
}

func findDeploySnapByMain(id uint64) ([]model.DeploymentSnap, error) {
	list := make([]model.DeploymentSnap, 0)
	err := base.Engine.Omit(model.AppIdInDeploymentSnap, model.EnvIdInDeploymentSnap, model.SpaceIdInDeploymentSnap, model.DeployRenderInDeploymentSnap).
		Where(model.DeployIdInDeploymentSnap+" = ?", id).
		Desc(model.IdInDeploymentSnap).Limit(10).
		Find(&list)
	return list, err
}

func insertDeploySnap(session *xorm.Session, info *model.DeploymentSnap) (int64, error) {
	return session.Omit(model.IdInPropertySnap, model.CreatedAt).
		InsertOne(info)
}
