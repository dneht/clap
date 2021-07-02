package refer

import (
	v1 "k8s.io/api/core/v1"
	"sort"
)

const (
	LabelAppId        = "clap.cana.io/app-id"
	LabelAppKey       = "clap.cana.io/app-key"
	LabelAppType      = "clap.cana.io/app-type"
	LabelAppBranch    = "clap.cana.io/app-branch"
	LabelAppName      = "app.kubernetes.io/name"
	LabelAppEnv       = "app.kubernetes.io/part-env"
	LabelAppSpace     = "app.kubernetes.io/part-of"
	LabelAppComponent = "app.kubernetes.io/component"
	LabelAppManaged   = "app.kubernetes.io/managed-by"

	SelectPodInnerType  = "inner"
	SelectPodExecType   = "exec"
	SelectPodAttachType = "attach"

	PodStatusRunning = "Running"
)

func BuildOneFromPod(pod *v1.Pod) *PodInfo {
	containers := make([]ContainerInfo, len(pod.Spec.Containers))
	for idx, value := range pod.Spec.Containers {
		containers[idx] = ContainerInfo{
			Name:  value.Name,
			Image: value.Image,
		}
	}
	return &PodInfo{
		PodName:    pod.Name,
		Namespace:  pod.Namespace,
		Component:  pod.Labels[LabelAppComponent],
		AppId:      pod.Labels[LabelAppId],
		AppKey:     pod.Labels[LabelAppKey],
		AppName:    pod.Labels[LabelAppName],
		AppType:    pod.Labels[LabelAppType],
		AppEnv:     pod.Labels[LabelAppEnv],
		AppSpace:   pod.Labels[LabelAppSpace],
		Containers: containers,
		Status:     pod.Status,
	}
}

func BuildListFromPod(podList *v1.PodList) []*PodInfo {
	pods := podList.Items
	infos := make([]*PodInfo, len(pods))
	for idx, pod := range pods {
		infos[idx] = BuildOneFromPod(&pod)
	}
	sort.SliceStable(infos, func(i, j int) bool {
		if nil == infos[i].Status.StartTime {
			return true
		}
		if nil == infos[j].Status.StartTime {
			return false
		}
		return infos[i].Status.StartTime.After(infos[j].Status.StartTime.Time)
	})
	return infos
}
