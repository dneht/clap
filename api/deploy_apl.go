package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	"xorm.io/xorm"
)

const DeployApiPre = "/api/deploy"

func GetDeploy(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "deploy id must be set")
	}
	info, err := getDeployById(id)
	return util.ResultParam(c, err, true, info)
}

func CheckDeploy(c *fiber.Ctx) error {
	id, err := util.CheckIdInput(c, "id")
	if nil != err {
		return util.ErrorInput(c, "deploy id must be set")
	}
	info, err := getDeployById(id)
	return util.ResultParam(c, err, true, info)
}

func ListDeploy(c *fiber.Ctx) error {
	param, err := util.CheckReferInput(c)
	if nil != err {
		return err
	}
	return util.ResultPageOrList(c, param,
		func(input *util.MainInput) (int64, error) {
			return countDeployWithPage(c, input, 0)
		}, func(input *util.MainInput) (int, interface{}, error) {
			return findDeployWithPage(c, input, 0)
		})
}

func CreateDeploy(c *fiber.Ctx) error {
	info, err := checkDeployInput(c)
	if err != nil {
		return err
	}
	result, err := insertDeploy(c, base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func UpdateDeploy(c *fiber.Ctx) error {
	info, err := checkDeployInput(c)
	if err != nil {
		return err
	}
	result, err := updateDeployById(c, base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func checkDeployInput(c *fiber.Ctx) (*model.Deployment, error) {
	info := new(model.Deployment)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputLog(c, err, "input is error")
	}
	return info, nil
}

func updateDeployStatus(c *fiber.Ctx, id uint64, status int, tag string) error {
	if id <= 0 || status == 0 {
		return errors.New("input id or status is empty")
	}
	_, err := base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		return updateDeployStatusById(c, session, id, status, tag)
	})
	if nil != err {
		return err
	}
	deployMap[id].DeployStatus = status
	if "" != tag {
		deployMap[id].DeployTag = tag
	}
	return nil
}

func ExecDeploy(c *fiber.Ctx) error {
	deployId, err := util.CheckIdInput(c, "deploy")
	if nil != err {
		return err
	}
	err = DeploymentAuth(c, deployId, AllowThisPackageDeploy)
	if nil != err {
		return err
	}
	appBase, _, err := getBaseModels(deployId)
	if nil != err {
		return err
	}

	if appBase.AppType == refer.NoneAppType {
		status, err := createPlatformApp(deployId)
		return util.ResultParamMapOne(c, err, "status", status)
	} else {
		selectType := c.Params("type")
		if "check" == selectType {
			status, pods, err := checkBuildJob(deployId)
			if nil != err {
				return err
			}
			if status.Succeeded > 0 {
				err = updateDeployStatus(c, deployId, refer.DeployStatusBuildEnd, "")
			} else if status.Failed > 0 {
				err = updateDeployStatus(c, deployId, refer.DeployStatusBuildFail, "")
			}
			return util.ResultParamMapTwo(c, err, "pods", pods, "status", status)
		} else if "build" == selectType {
			tag, status, err := createBuildJob(deployId)
			if nil != err {
				return err
			}
			err = updateDeployStatus(c, deployId, refer.DeployStatusBuilding, tag)
			return util.ResultParamMapTwo(c, err, "tag", tag, "status", status)
		} else if "deploy" == selectType {
			status, err := createTemplateApp(deployId)
			if nil != err {
				return err
			}
			err = updateDeployStatus(c, deployId, refer.DeployStatusPackHear, "")
			return util.ResultParamMapOne(c, err, "status", status)
		}
		return errors.New("select type is not support")
	}
}
