package refer

import "cana.io/clap/pkg/model"

const NoneAppType = 0

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

var AppTypeMap = map[int]string {
	5: "Nginx",
	10: "Java",
	11: "Tomcat",
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

func ConvertAppType(appType int) string {
	get, ok := AppTypeMap[appType]
	if ok {
		return get
	}
	return "None"
}
