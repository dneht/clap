package api

import (
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"strconv"
)

func selectAppLabelList(envBase *model.Environment, appBase *model.Project) *[]string {
	return &[]string{refer.LabelAppEnv + "=" + envBase.Env,
		refer.LabelAppName + "=" + appBase.AppName,
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType)}
}

func selectDeployLabelList(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project) *[]string {
	return &[]string{refer.LabelAppEnv + "=" + envBase.Env,
		refer.LabelAppSpace + "=" + spaceBase.SpaceName,
		refer.LabelAppName + "=" + appBase.AppName,
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType)}
}

func selectJobLabelMap(jobName string) map[string]string {
	return map[string]string{
		refer.LabelAppType:      "build",
		refer.LabelAppName:      jobName,
		refer.LabelAppComponent: "deployment",
		refer.LabelAppManaged:   "Clap",
	}
}

func selectJobLabelList(jobName string) *[]string {
	return &[]string{refer.LabelAppType + "=build",
		refer.LabelAppName + "=" + jobName,
		refer.LabelAppComponent + "=deployment",
		refer.LabelAppManaged + "=Clap"}
}

func buildAppLabelString(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, component string) string {
	return refer.LabelAppId + "=" + strconv.FormatUint(appBase.Id, 10) + "," +
		refer.LabelAppKey + "=" + appBase.AppKey + "," +
		refer.LabelAppType + "=" + refer.ConvertAppType(appBase.AppType) + "," +
		refer.LabelAppName + "=" + appBase.AppName + "," +
		refer.LabelAppComponent + "=" + component + "," +
		refer.LabelAppEnv + "=" + envBase.Env + "," +
		refer.LabelAppSpace + "=" + spaceBase.SpaceName + "," +
		refer.LabelAppManaged + "=Clap"
}

func configAppLabelMap(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, templateName string) map[string]string {
	return map[string]string{
		refer.LabelAppId:        strconv.FormatUint(appBase.Id, 10),
		refer.LabelAppKey:       appBase.AppKey,
		refer.LabelAppType:      refer.ConvertAppType(appBase.AppType),
		refer.LabelAppName:      appBase.AppName,
		refer.LabelAppComponent: templateName,
		refer.LabelAppEnv:       envBase.Env,
		refer.LabelAppSpace:     spaceBase.SpaceName,
		refer.LabelAppManaged:   "Clap",
	}
}

func selectAppLabelMap(envBase *model.Environment, spaceBase *model.EnvironmentSpace, appBase *model.Project, templateName string) map[string]string {
	return map[string]string{
		refer.LabelAppType:      refer.ConvertAppType(appBase.AppType),
		refer.LabelAppName:      appBase.AppName,
		refer.LabelAppComponent: templateName,
		refer.LabelAppEnv:       envBase.Env,
		refer.LabelAppSpace:     spaceBase.SpaceName,
		refer.LabelAppManaged:   "Clap",
	}
}
