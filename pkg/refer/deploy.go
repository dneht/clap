package refer

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type DeployStatus struct {
	Status string `json:"status,omitempty"`
}

type Deployment struct {
	Main *appsv1.Deployment `json:"main,omitempty"`
	CommonExtend
}

type StatefulSet struct {
	Main *appsv1.StatefulSet `json:"main,omitempty"`
	CommonExtend
}

type CommonExtend struct {
	Secrets  *map[string]*v1.Secret                 `json:"secrets,omitempty"`
	Configs  *map[string]*v1.ConfigMap              `json:"configs,omitempty"`
	Services *map[string]*v1.Service                `json:"services,omitempty"`
	Contours *map[string]*unstructured.Unstructured `json:"contours,omitempty"`
	Policy   *policyv1.PodSecurityPolicy            `json:"policy,omitempty"`
}
