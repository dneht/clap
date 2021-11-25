package refer

import "cana.io/clap/pkg/model"

const (
	KindConfigMap           = "ConfigMap"
	KindSecret              = "Secret"
	KindDeploy              = "Deployment"
	KindStateful            = "StatefulSet"
	KindJob                 = "Job"
	KindCron                = "CronJob"
	KindService             = "Service"
	KindIngress             = "Ingress"
	KindHttpProxy           = "HTTPProxy"
	KindClusterRole         = "ClusterRole"
	KindClusterRoleBinding  = "ClusterRoleBinding"
	KindRole                = "Role"
	KindRoleBinding         = "RoleBinding"
	KindServiceAccount      = "ServiceAccount"
	KindPodAutoscalert      = "HorizontalPodAutoscaler"
	KindPodDisruptionBudget = "PodDisruptionBudget"
	KindPodSecurityPolicy   = "PodSecurityPolicy"
)

const (
	NoneAppType     = 0
	NoneAppTypeName = "None"
)

var AppTypeMap = map[int]string{
	5:  "Nginx",
	10: "Java",
	11: "Gradle",
	16: "Tomcat",
	20: "Go",
	60: "Python",
	90: "Node",
}

func GetAppName(app *model.Project, space *model.EnvironmentSpace) string {
	if space.IsControl == 1 {
		return app.AppName
	} else {
		return app.AppName + "-" + space.SpaceName
	}
}

func GetConfigName(app *model.Project, space *model.EnvironmentSpace, name string) string {
	return GetAppName(app, space) + "-" + name
}

func CheckAppTypeName(appTypeName string) int {
	if "" == appTypeName {
		return -1
	}
	for key, value := range AppTypeMap {
		if appTypeName == value {
			return key
		}
	}
	return -1
}

func ConvertAppType(appType int) string {
	get, ok := AppTypeMap[appType]
	if ok {
		return get
	}
	return NoneAppTypeName
}
