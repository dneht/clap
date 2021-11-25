package refer

import (
	v1 "k8s.io/api/core/v1"
)

type PodInfo struct {
	PodName    string           `json:"podName,omitempty"`
	Namespace  string           `json:"namespace,omitempty"`
	Component  string           `json:"component,omitempty"`
	AppId      string           `json:"appId,omitempty"`
	AppName    string           `json:"appName,omitempty"`
	AppType    string           `json:"appType,omitempty"`
	AppEnv     string           `json:"appEnv,omitempty"`
	AppSpace   string           `json:"appSpace,omitempty"`
	Containers []ContainerInfo `json:"containers,omitempty"`
	Status     v1.PodStatus    `json:"status,omitempty"`
}

type ContainerInfo struct {
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}
