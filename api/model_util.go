package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"errors"
	"strconv"
)

func getBaseModels(deployId uint64) (*model.Project, *model.Deployment, error) {
	deployBase, err := getDeployById(deployId)
	if nil != err {
		return nil, nil, err
	}
	appBase, err := getAppById(deployBase.AppId)
	if nil != err {
		return nil, nil, err
	}
	return appBase, deployBase, nil
}

func getMoreModels(deployId uint64) (*model.Environment, *model.EnvironmentSpace, *model.Project, *model.Deployment, error) {
	appBase, deployBase, err := getBaseModels(deployId)
	if nil != err {
		return nil, nil, nil, nil, err
	}
	envBase, err := base.Env(deployBase.EnvId)
	if nil != err {
		return nil, nil, nil, nil, err
	}
	spaceBase, err := getSpaceById(deployBase.SpaceId)
	if nil != err {
		return nil, nil, nil, nil, err
	}
	return envBase, spaceBase, appBase, deployBase, err
}

func getMoreInfos(envId, projectId, deployId uint64) (*refer.DeployInfo, *refer.AppInfo, error) {
	baseDeployInfo, appDeployInfo, err := getDeployAppInfoById(deployId)
	if nil != err {
		return nil, nil, err
	}
	if envId != baseDeployInfo.EnvId {
		return nil, nil, errors.New("input env id " + strconv.FormatUint(envId, 10) +
			" not eq" + strconv.FormatUint(baseDeployInfo.EnvId, 10))
	}
	if projectId != baseDeployInfo.AppId {
		return nil, nil, errors.New("input app id " + strconv.FormatUint(projectId, 10) +
			" not eq" + strconv.FormatUint(baseDeployInfo.AppId, 10))
	}

	envDeploy, err := base.Deploy(baseDeployInfo.EnvId)
	if nil != err {
		return nil, nil, err
	}
	_, spaceBaseInfo, err := getBaseSpaceInfoById(baseDeployInfo.SpaceId)
	if nil != err {
		return nil, nil, err
	}
	_, appBaseInfo, err := getBaseAppInfoById(baseDeployInfo.AppId)
	if nil != err {
		return nil, nil, err
	}
	return envDeploy, mergeAppInfo(spaceBaseInfo, appBaseInfo, appDeployInfo), err
}

func getAndCheckTemplate(templateId uint64) (*model.Template, error) {
	templateBase, err := getTemplateById(templateId)
	if nil != err {
		return nil, err
	}
	if "" == templateBase.TemplateContent {
		return nil, errors.New("template is empty")
	}
	return templateBase, err
}

func mergeAppInfo(spaceInfo *refer.SpaceInfo, baseInfo, deployInfo *refer.AppInfo) *refer.AppInfo {
	if deployInfo.Type == 0 {
		deployInfo.Type = baseInfo.Type
	}
	if "" == deployInfo.Port {
		deployInfo.Port = baseInfo.Port
	}
	if "" == deployInfo.From {
		deployInfo.From = baseInfo.From
	}
	if deployInfo.Template == 0 {
		deployInfo.Template = baseInfo.Template
	}
	if "" == deployInfo.Component {
		deployInfo.Component = baseInfo.Component
	}
	if nil == deployInfo.Ready {
		deployInfo.Ready = baseInfo.Ready
	}
	if nil == deployInfo.Factor {
		deployInfo.Factor = baseInfo.Factor
	}
	deployInfo.Conf = mergeStringMap(baseInfo.Conf, deployInfo.Conf)
	deployInfo.Code = mergeStringMap(baseInfo.Code, deployInfo.Code)
	deployInfo.Repo = mergeStringMap(baseInfo.Repo, deployInfo.Repo)
	deployInfo.Param = mergeParamMap(mergeParamMap(spaceInfo.Param, baseInfo.Param), deployInfo.Param)
	return deployInfo
}

func mergeStringMap(base, next map[string]string) map[string]string {
	newMap := make(map[string]string, len(base)+4)
	for key, value := range base {
		newMap[key] = value
	}
	for key, value := range next {
		newMap[key] = value
	}
	return newMap
}

func mergeParamMap(base, next map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(base)+4)
	for key, value := range base {
		if nil != value {
			newMap[key] = value
		}
	}
	for key, value := range next {
		if nil != value {
			newMap[key] = value
		}
	}
	return newMap
}
