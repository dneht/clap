package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"log"
	"strconv"
	"strings"
	"time"
)

const jobNamePre = "clap-build-"

func checkBuildJob(deployId uint64) (*batchv1.JobStatus, *[]*refer.PodInfo, error) {
	deployBase, err := getDeployById(deployId)
	if nil != err {
		return nil, nil, err
	}
	timeTag := deployBase.DeployTag
	if "" == timeTag {
		return nil, nil, errors.New("not package yet")
	}
	envId := deployBase.EnvId
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return nil, nil, err
	}

	nowEnv := base.Now()
	job, err := k8s.BatchV1().Jobs(nowEnv.Namespace).Get(context.TODO(), jobNamePre+timeTag, metav1.GetOptions{})
	if nil != err {
		return nil, nil, err
	}
	podList, err := listPodByLabel(envId, nowEnv.Namespace, selectJobLabelList(jobNamePre+timeTag))
	var pods []*refer.PodInfo
	if nil == err {
		pods = refer.BuildListFromPod(podList)
	}
	return &job.Status, &pods, err
}

func createBuildJob(deployId uint64) (string, *batchv1.JobStatus, error) {
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return "", nil, err
	}
	if appBase.AppType == refer.NoneAppType {
		return "", nil, errors.New("no need to pack")
	}
	k8s, _, err := base.K8S(envBase.Id)
	if nil != err {
		return "", nil, err
	}
	if deployBase.DeployStatus == refer.DeployStatusBuilding {
		// TODO check job status
		return "", nil, errors.New("packaging")
	}
	_, appInfo, err := getMoreInfos(envBase.Id, appBase.Id, deployId)
	if nil != err {
		return "", nil, err
	}
	repoUrl, codeUrl, err := checkProjectParam(spaceBase, appBase, appInfo)
	if nil != err {
		return "", nil, err
	}
	nowEnv := base.Now()
	timeTag := generateTimeTag()
	jobName := jobNamePre + timeTag
	deploymentProp := nowEnv.Package
	var imagePullSecret []corev1.LocalObjectReference
	if "" != deploymentProp.ImagePullSecret {
		imagePullSecret = []corev1.LocalObjectReference{{Name: deploymentProp.ImagePullSecret}}
	}
	var hostAliases []corev1.HostAlias
	appHosts, ok := appInfo.Param["hostAliases"]
	if ok {
		hostsBytes, err := json.Marshal(appHosts)
		if nil == err {
			_ = json.Unmarshal(hostsBytes, &hostAliases)
		}
	}
	extraBytes, err := json.Marshal(appInfo)
	if nil != err {
		return "", nil, err
	}
	setLabels := selectJobLabelMap(jobName)
	job, err := k8s.BatchV1().Jobs(nowEnv.Namespace).Create(context.TODO(), &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &deploymentProp.BackoffLimit,
			TTLSecondsAfterFinished: &deploymentProp.CleanAfterFinished,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: setLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            jobName,
							Image:           deploymentProp.BuildJobImage,
							ImagePullPolicy: corev1.PullPolicy(deploymentProp.ImagePullPolicy),
							Env: []corev1.EnvVar{
								{Name: "APP_ID", Value: strconv.FormatUint(appBase.Id, 10)},
								{Name: "APP_TYPE", Value: strings.ToLower(refer.ConvertAppType(appBase.AppType))},
								{Name: "APP_KEY", Value: appBase.AppKey},
								{Name: "APP_NAME", Value: appBase.AppName},
								{Name: "APP_ENV", Value: envBase.Env},
								{Name: "APP_SPACE", Value: spaceBase.SpaceName},
								{Name: "APP_LABELS", Value: buildAppLabelString(envBase, spaceBase, appBase, appInfo.Component)},
								{Name: "APP_EXTRA", Value: string(extraBytes)},
								{Name: "ENV_BASE", Value: envBase.DeployInfo},
								{Name: "ENV_SPEC", Value: envBase.FormatInfo},
								{Name: "ENV_SPACE", Value: spaceBase.SpaceInfo},
								{Name: "REPO_URL", Value: repoUrl},
								{Name: "CODE_URL", Value: codeUrl},
								{Name: "CODE_BRANCH", Value: deployBase.BranchName},
								{Name: "TIME_TAG", Value: timeTag},
								{Name: "SKIP_TEST", Value: util.ConvertBoolToYesOrNo(deploymentProp.MavenSkipTests)},
							},
						},
					},
					HostAliases:      hostAliases,
					RestartPolicy:    corev1.RestartPolicyNever,
					ImagePullSecrets: imagePullSecret,
				},
			},
		},
	}, metav1.CreateOptions{})
	if nil != err {
		return "", nil, err
	}
	return timeTag, &job.Status, err
}

func createPlatformApp(deployId uint64) (refer.DeployStatus, error) {
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	_, appInfo, err := getMoreInfos(envBase.Id, appBase.Id, deployId)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	repoFullUrl, err := checkPlatformParam(appInfo)
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	appKind, jsonStr, err := renderJsonnet(envBase, spaceBase, appBase, deployBase, repoFullUrl)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	namespace := spaceBase.SpaceKeep
	appDeployName := refer.GetAppName(appBase, spaceBase)
	if appKind == refer.KindDeploy {
		return createDeployment(envBase.Id, namespace, jsonStr, appDeployName)
	} else if appKind == refer.KindStateful {
		return createStatefulSet(envBase.Id, namespace, jsonStr, appDeployName)
	} else {
		return refer.FailedDeployStatus, errors.New("can not support this kind: " + appKind)
	}
}

func createTemplateApp(deployId uint64) (refer.DeployStatus, error) {
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	timeTag := deployBase.DeployTag
	if "" == timeTag {
		return refer.FailedDeployStatus, errors.New("not package yet")
	}
	_, appInfo, err := getMoreInfos(envBase.Id, appBase.Id, deployId)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	repoUrl, _, err := checkProjectParam(spaceBase, appBase, appInfo)
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	appKind, jsonStr, err := renderJsonnet(envBase, spaceBase, appBase, deployBase, repoUrl+":"+timeTag)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	namespace := spaceBase.SpaceKeep
	appDeployName := refer.GetAppName(appBase, spaceBase)

	util.DingDingMessage(map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": fmt.Sprintf("正在[%s]环境的[%s]空间发布[%s|%s]项目",
				envBase.Env, spaceBase.SpaceName, appBase.AppName, appBase.AppDesc),
		},
	})
	if appKind == refer.KindDeploy {
		return createDeployment(envBase.Id, namespace, jsonStr, appDeployName)
	} else if appKind == refer.KindStateful {
		return createStatefulSet(envBase.Id, namespace, jsonStr, appDeployName)
	} else {
		return refer.FailedDeployStatus, errors.New("can not support this kind: " + appKind)
	}
}

// generate tag and create job
func generateTimeTag() string {
	location, err := time.LoadLocation(base.Now().Timezone)
	if nil != err {
		location = time.Local
	}
	return strings.Replace(time.Now().In(location).
		Format("20060102150405.999999999"), ".", "", 1)
}

func checkPlatformParam(appInfo *refer.AppInfo) (string, error) {
	if nil == appInfo.Ready || "" == appInfo.Ready.Url || "" == appInfo.Ready.Version {
		return "", errors.New("this app type must provide image url and version")
	}
	return appInfo.Ready.Url + ":" + appInfo.Ready.Version, nil
}

func checkProjectParam(spaceBase *model.EnvironmentSpace, appBase *model.Project, appInfo *refer.AppInfo) (string, string, error) {
	var spaceInfo = new(refer.SpaceInfo)
	err := json.Unmarshal([]byte(spaceBase.SpaceInfo), spaceInfo)
	if nil != err {
		return "", "", err
	}
	repoUrl, ok := appInfo.Repo["url"]
	if !ok {
		repoUrl = util.IfNotExistTailAdd(spaceInfo.Repo.Base, "/") + appBase.AppName
	}
	repoIdx := strings.Index(repoUrl, "://")
	if repoIdx > 0 {
		repoUrl = repoUrl[repoIdx+3:]
	}
	codeUrl, ok := appInfo.Code["url"]
	if !ok {
		codeUrl = util.IfNotExistTailAdd(spaceInfo.Code.Base, "/") + appBase.AppName
	}
	return repoUrl, codeUrl, nil
}

func createDeployment(envId uint64, namespace, jsonStr, appDeployName string) (refer.DeployStatus, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	var median = new(refer.Deployment)
	err = json.Unmarshal([]byte(jsonStr), median)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	get, err := k8s.AppsV1().Deployments(namespace).Get(context.TODO(), appDeployName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return refer.FailedDeployStatus, err
		}
	}
	if nil == get {
		get, err = k8s.AppsV1().Deployments(namespace).Create(context.TODO(), median.Main, metav1.CreateOptions{})
	} else {
		get, err = k8s.AppsV1().Deployments(namespace).Update(context.TODO(), median.Main, metav1.UpdateOptions{})
	}
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	status := refer.DefaultDeployStatus
	status, err = createService(envId, &median.CommonExtend, namespace, status)
	status, err = createContour(envId, &median.CommonExtend, namespace, status)
	status, err = createSecret(envId, &median.CommonExtend, namespace, status)
	status, err = createConfig(envId, &median.CommonExtend, namespace, status)
	return status, err
}

func createStatefulSet(envId uint64, namespace, jsonStr, appDeployName string) (refer.DeployStatus, error) {
	k8s, _, err := base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	var median = new(refer.StatefulSet)
	err = json.Unmarshal([]byte(jsonStr), median)
	if nil != err {
		return refer.FailedDeployStatus, err
	}
	get, err := k8s.AppsV1().StatefulSets(namespace).Get(context.TODO(), appDeployName, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			get = nil
			err = nil
		} else {
			return refer.FailedDeployStatus, err
		}
	}
	if nil == get {
		get, err = k8s.AppsV1().StatefulSets(namespace).Create(context.TODO(), median.Main, metav1.CreateOptions{})
	} else {
		get, err = k8s.AppsV1().StatefulSets(namespace).Update(context.TODO(), median.Main, metav1.UpdateOptions{})
	}
	if nil != err {
		return refer.FailedDeployStatus, err
	}

	status := refer.DefaultDeployStatus
	status, err = createService(envId, &median.CommonExtend, namespace, status)
	status, err = createContour(envId, &median.CommonExtend, namespace, status)
	status, err = createSecret(envId, &median.CommonExtend, namespace, status)
	status, err = createConfig(envId, &median.CommonExtend, namespace, status)
	return status, err
}

func createService(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Services {
		return status, nil
	}
	serviceList := *extend.Services
	if len(serviceList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for serviceName, serviceData := range serviceList {
		var get *corev1.Service
		get, err = k8s.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
		if nil != err {
			if k8serror.IsNotFound(err) {
				get = nil
				err = nil
			} else {
				//TODO add status
				log.Printf("[error] create svc: %v", err)
				continue
			}
		}
		if nil == get {
			get, err = k8s.CoreV1().Services(namespace).Create(context.TODO(), serviceData, metav1.CreateOptions{})
		} else {
			get, err = k8s.CoreV1().Services(namespace).Update(context.TODO(), serviceData, metav1.UpdateOptions{})
		}
	}

	return status, err
}

func createContour(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Contours {
		return status, nil
	}
	contourList := *extend.Contours
	if len(contourList) <= 0 {
		return status, nil
	}

	crd, _, err := base.K8D(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for contourName, contourData := range contourList {
		var get *unstructured.Unstructured
		get, err = crd.Resource(refer.ContourGvr).Namespace(namespace).Get(context.TODO(), contourName, metav1.GetOptions{})
		if nil != err {
			if k8serror.IsNotFound(err) {
				get = nil
				err = nil
			} else {
				//TODO add status
				log.Printf("[error] create contour: %v", err)
				continue
			}
		}

		if nil == get {
			get, err = crd.Resource(refer.ContourGvr).Namespace(namespace).Create(context.TODO(), contourData, metav1.CreateOptions{})
		} else {
			contourData.SetResourceVersion(get.GetResourceVersion())
			get, err = crd.Resource(refer.ContourGvr).Namespace(namespace).Update(context.TODO(), contourData, metav1.UpdateOptions{})
		}
	}

	return status, err
}

func createSecret(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Secrets {
		return status, nil
	}
	secretList := *extend.Secrets
	if len(secretList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for secretName, secretData := range secretList {
		if nil == secretData {
			continue
		}
		var get *corev1.Secret
		get, err = k8s.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
		if nil != err {
			if k8serror.IsNotFound(err) {
				get = nil
				err = nil
			} else {
				//TODO add status
				log.Printf("[error] create secret: %v", err)
				continue
			}
		}
		if nil == get {
			get, err = k8s.CoreV1().Secrets(namespace).Create(context.TODO(), secretData, metav1.CreateOptions{})
		} else {
			get, err = k8s.CoreV1().Secrets(namespace).Update(context.TODO(), secretData, metav1.UpdateOptions{})
		}
	}

	return status, err
}

func createConfig(envId uint64, extend *refer.CommonExtend, namespace string, status refer.DeployStatus) (refer.DeployStatus, error) {
	if nil == extend.Configs {
		return status, nil
	}
	configList := *extend.Configs
	if len(configList) <= 0 {
		return status, nil
	}

	var err error
	var k8s *kubernetes.Clientset
	k8s, _, err = base.K8S(envId)
	if nil != err {
		return refer.ConnectErrorDeployStatus, err
	}

	for configName, configData := range configList {
		if nil == configData {
			continue
		}
		var get *corev1.ConfigMap
		get, err = k8s.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configName, metav1.GetOptions{})
		if nil != err {
			if k8serror.IsNotFound(err) {
				get = nil
				err = nil
			} else {
				//TODO add status
				log.Printf("[error] create config: %v", err)
				continue
			}
		}
		if nil == get {
			get, err = k8s.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configData, metav1.CreateOptions{})
		} else {
			get, err = k8s.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configData, metav1.UpdateOptions{})
		}
	}

	return status, err
}

//TODO createBudget
func createBudget() {

}

//TODO createPolicy
func createPolicy() {

}
