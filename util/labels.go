package util

import (
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"strconv"
)

func SelectAppLabelList(envBase *model.Environment, appBase *model.Project) *[]string {
	return &[]string{refer.LabelAppEnv + "=" + envBase.Env,
		refer.LabelAppName + "=" + appBase.AppName,
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType)}
}

func SelectDeployLabelList(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project) *[]string {
	return &[]string{refer.LabelAppEnv + "=" + envBase.Env,
		refer.LabelAppSpace + "=" + spaceBase.SpaceName,
		refer.LabelAppName + "=" + appBase.AppName,
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType)}
}

func SelectJobLabelMap(appType, jobName string) map[string]string {
	return map[string]string{
		refer.LabelAppType:      appType,
		refer.LabelAppName:      jobName,
		refer.LabelAppComponent: "deployment",
		refer.LabelAppManaged:   "Clap",
	}
}

func SelectJobLabelList(jobName string) *[]string {
	return &[]string{refer.LabelAppType + "=build",
		refer.LabelAppName + "=" + jobName,
		refer.LabelAppComponent + "=deployment",
		refer.LabelAppManaged + "=Clap"}
}

func BuildAppLabelString(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, component string) string {
	return refer.LabelAppId + "=" + strconv.FormatUint(appBase.Id, 10) + "," +
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType) + "," +
		refer.LabelAppName + "=" + appBase.AppName + "," +
		refer.LabelAppComponent + "=" + component + "," +
		refer.LabelAppEnv + "=" + envBase.Env + "," +
		refer.LabelAppSpace + "=" + spaceBase.SpaceName + "," +
		refer.LabelAppManaged + "=Clap"
}

func ConfigAppLabelMap(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, templateName string) map[string]string {
	return map[string]string{
		refer.LabelAppId:        strconv.FormatUint(appBase.Id, 10),
		refer.LabelAppType:      refer.ConvertAppType(appBase.AppType),
		refer.LabelAppName:      appBase.AppName,
		refer.LabelAppComponent: templateName,
		refer.LabelAppEnv:       envBase.Env,
		refer.LabelAppSpace:     spaceBase.SpaceName,
		refer.LabelAppManaged:   "Clap",
	}
}

func SelectAppLabelMap(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, templateName string) map[string]string {
	return map[string]string{
		refer.LabelAppType:      refer.ConvertAppType(appBase.AppType),
		refer.LabelAppName:      appBase.AppName,
		refer.LabelAppComponent: templateName,
		refer.LabelAppEnv:       envBase.Env,
		refer.LabelAppSpace:     spaceBase.SpaceName,
		refer.LabelAppManaged:   "Clap",
	}
}
