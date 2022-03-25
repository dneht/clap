package api

import (
	"cana.io/clap/pkg/base"
	"cana.io/clap/pkg/log"
	"cana.io/clap/pkg/model"
	"cana.io/clap/pkg/refer"
	"cana.io/clap/util"
	"context"
	"encoding/json"
	"errors"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	conf := base.Now()
	job, err := k8s.BatchV1().Jobs(conf.Namespace).Get(context.TODO(), jobNamePre+timeTag, metav1.GetOptions{})
	if nil != err {
		if k8serror.IsNotFound(err) {
			return &batchv1.JobStatus{Succeeded: 1}, &[]*refer.PodInfo{}, nil
		} else {
			return nil, nil, err
		}
	}
	podList, err := listPodByLabel(envId, conf.Namespace, util.SelectJobLabelList(jobNamePre+timeTag))
	var pods []*refer.PodInfo
	if nil == err {
		pods = refer.BuildListFromPod(podList)
	}
	return &job.Status, &pods, err
}

func createBuildJob(deployId uint64, branchName string) (string, *batchv1.JobStatus, error) {
	envBase, spaceBase, appBase, deployBase, err := getMoreModels(deployId)
	if nil != err {
		return "", nil, err
	}
	if "" == branchName {
		branchName = deployBase.BranchName
	} else {
		err = updateDeployBranch(deployId, branchName)
		if nil != err {
			return "", nil, err
		}
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
	conf := base.Now()
	timeTag := generateTimeTag()
	jobName := jobNamePre + timeTag
	deploymentProp := conf.Package
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
	setLabels := util.SelectJobLabelMap("build", jobName)
	job, err := k8s.BatchV1().Jobs(conf.Namespace).Create(context.TODO(), &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:   jobName,
			Labels: setLabels,
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
								{Name: "APP_NAME", Value: appBase.AppName},
								{Name: "APP_ENV", Value: envBase.Env},
								{Name: "APP_SPACE", Value: spaceBase.SpaceName},
								{Name: "APP_LABELS", Value: util.BuildAppLabelString(envBase, spaceBase, appBase, appInfo.Component)},
								{Name: "APP_EXTRA", Value: string(extraBytes)},
								{Name: "ENV_BASE", Value: envBase.DeployInfo},
								{Name: "ENV_SPEC", Value: envBase.FormatInfo},
								{Name: "ENV_SPACE", Value: spaceBase.SpaceInfo},
								{Name: "REPO_URL", Value: repoUrl},
								{Name: "CODE_URL", Value: codeUrl},
								{Name: "CODE_BRANCH", Value: branchName},
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

func deleteBuildJob(deployId uint64) {
	deployBase, err := getDeployById(deployId)
	if nil != err {
		log.Infof("get deploy detail %v failed: %v", deployId, err)
		return
	}
	timeTag := deployBase.DeployTag
	if "" == timeTag {
		return
	}
	conf := base.Now()
	k8s, _, err := base.K8S(deployBase.EnvId)
	if nil != err {
		return
	}
	foreground := metav1.DeletePropagationForeground
	err = k8s.BatchV1().Jobs(conf.Namespace).Delete(context.TODO(), jobNamePre+timeTag,
		metav1.DeleteOptions{PropagationPolicy: &foreground})
	if nil != err {
		log.Infof("delete build job %v failed: %v", jobNamePre+timeTag, err)
	}
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
		return handleDeployment(envBase.Id, namespace, jsonStr, appDeployName)
	} else if appKind == refer.KindStateful {
		return handleStatefulSet(envBase.Id, namespace, jsonStr, appDeployName)
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

	util.DingDingDeployMessage(envBase.Env, spaceBase.SpaceName, appBase.AppName, appBase.AppDesc)
	if appKind == refer.KindDeploy {
		return handleDeployment(envBase.Id, namespace, jsonStr, appDeployName)
	} else if appKind == refer.KindStateful {
		return handleStatefulSet(envBase.Id, namespace, jsonStr, appDeployName)
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
