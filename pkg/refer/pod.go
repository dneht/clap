package refer

import (
	v1 "k8s.io/api/core/v1"
)

type PodInfo struct {
	PodName    string           `json:"podName"`
	Namespace  string           `json:"namespace"`
	Component  string           `json:"component"`
	AppId      string           `json:"appId"`
	AppKey     string           `json:"appKey"`
	AppName    string           `json:"appName"`
	AppType    string           `json:"appType"`
	AppEnv     string           `json:"appEnv"`
	AppSpace   string           `json:"appSpace"`
	Containers []ContainerInfo `json:"containers"`
	Status     v1.PodStatus    `json:"status"`
}

type ContainerInfo struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}
