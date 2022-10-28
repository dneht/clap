package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"errors"
	"github.com/gofiber/fiber/v2"
	batchv1 "k8s.io/api/batch/v1"
	"time"
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
	result, err := insertDeploy(base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func UpdateDeploy(c *fiber.Ctx) error {
	info, err := checkDeployInput(c)
	if err != nil {
		return err
	}
	result, err := updateDeployById(base.Engine.NewSession(), info)
	return util.ResultParamWithMessage(c, err, result > 0, "create app error", info.Id)
}

func checkDeployInput(c *fiber.Ctx) (*model.Deployment, error) {
	info := new(model.Deployment)
	if err := c.BodyParser(info); err != nil {
		return nil, util.ErrorInputLog(c, err, "input is error")
	}
	return info, nil
}

func updateDeployStatus(id uint64, status int, tag string) error {
	if id <= 0 || status == 0 {
		return errors.New("input id or status is empty")
	}
	_, err := base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		return updateDeployStatusById(session, id, status, tag)
	})
	if nil != err {
		return err
	}
	get, ok := deployMap[id]
	if ok {
		get.DeployStatus = status
		if "" != tag {
			get.DeployTag = util.StringClone(tag)
		}
	}
	return nil
}

func updateDeployBranch(id uint64, branch string) error {
	if id <= 0 || "" == branch {
		return errors.New("input id or status is empty")
	}
	_, err := base.Engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		return updateDeployBranchById(session, id, branch)
	})
	if nil != err {
		return err
	}
	get, ok := deployMap[id]
	if ok {
		get.BranchName = util.StringClone(branch)
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
	appBase, deployBase, err := getBaseModels(deployId)
	if nil != err {
		return err
	}
	userId := contextUserId(c)

	if appBase.AppType == refer.NoneAppType {
		status, err := createPlatformApp(userId, deployId)
		return util.ResultParamMapOne(c, err, "status", status)
	} else {
		selectType := c.Params("type")
		branchName := ""
		if deployBase.IsBranch != 0 {
			branchName = c.Query("branch")
		}
		if "check" == selectType {
			_, pods, status, err := execCheckDeploy(deployId)
			return util.ResultParamMapTwo(c, err, "pods", pods, "status", status)
		} else if "build" == selectType {
			tag, status, err := execBuildDeploy(deployId, branchName, deployBase)
			return util.ResultParamMapTwo(c, err, "tag", tag, "status", status)
		} else if "cancel" == selectType {
			return util.ResultParamEmpty(c, execCancelDeploy(deployId))
		} else if "deploy" == selectType {
			status, err := execDeployDeploy(userId, deployId)
			return util.ResultParamMapOne(c, err, "status", status)
		} else if "auto_deploy" == selectType {
			tag, status, err := execBuildDeploy(deployId, branchName, deployBase)
			if nil != err {
				return util.ErrorInternal(c, err)
			}
			go execAutoCheckDeploy(userId, deployId)
			return util.ResultParamMapTwo(c, nil, "tag", tag, "status", status)
		}
		return errors.New("select type is not support")
	}
}

func execCheckDeploy(deployId uint64) (int, []*refer.PodInfo, string, error) {
	status, pods, err := checkBuildJob(deployId)
	if nil != err {
		return refer.DeployStatusBuilding, nil, "", err
	}
	get, err := getDeployById(deployId)
	if nil != err {
		return refer.DeployStatusBuilding, nil, "", err
	}
	if get.DeployStatus == refer.DeployStatusPackHear {
		return refer.DeployStatusPackHear, pods, refer.CompleteDeployStatus.Status, nil
	}

	res := refer.DeployStatusBuilding
	ds := refer.DefaultDeployStatus
	limit := base.Now().Package.BackoffLimit
	if status.Succeeded > 0 {
		res = refer.DeployStatusBuildEnd
		ds = refer.SuccessDeployStatus
		err = updateDeployStatus(deployId, refer.DeployStatusBuildEnd, "")
	} else if status.Failed > limit {
		log.Infof("build job failed: %v, %v, %v", deployId, limit, status)
		res = refer.DeployStatusBuildFail
		ds = refer.FailedDeployStatus
		err = updateDeployStatus(deployId, refer.DeployStatusBuildFail, "")
	}
	return res, pods, ds.Status, err
}

func execBuildDeploy(deployId uint64, branchName string, deployBase *model.Deployment) (string, *batchv1.JobStatus, error) {
	if deployBase.IsPackage == 0 {
		return "", nil, errors.New("this deploy can not package")
	}
	tag, status, err := createBuildJob(deployId, branchName)
	if nil != err {
		return "", nil, err
	}
	return tag, status, updateDeployStatus(deployId, refer.DeployStatusBuilding, tag)
}

func execCancelDeploy(deployId uint64) error {
	deleteBuildJob(deployId)
	deployGet, ok := deployMap[deployId]
	if ok {
		deployGet.DeployStatus = refer.DeployStatusPackHear
	}
	return updateDeployStatus(deployId, refer.DeployStatusPackHear, "")
}

func execDeployDeploy(userId, deployId uint64) (*refer.DeployStatus, error) {
	deleteBuildJob(deployId)
	status, err := createTemplateApp(userId, deployId)
	if nil != err {
		return nil, err
	}
	return &status, updateDeployStatus(deployId, refer.DeployStatusPackHear, "")
}

func execAutoCheckDeploy(userId, deployId uint64) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	num, limit := 0, 300
	for range ticker.C {
		res, _, _, _ := execCheckDeploy(deployId)
		if res == refer.DeployStatusBuildEnd {
			log.Infof("auto deploy[%v] package success", deployId)
			break
		} else if res == refer.DeployStatusBuildFail {
			log.Warnf("auto deploy[%v] package failed", deployId)
			return
		}
		num += 1
		if num > limit {
			log.Warnf("auto deploy[%v] package to long", deployId)
			return
		}
	}
	status, err := execDeployDeploy(userId, deployId)
	if nil != err {
		log.Error("auto deploy[%v] failed: %v, %v", deployId, status, err)
	}
}
