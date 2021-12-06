/*
Copyright 2020 Dasheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package refer

type CommonParamInfo struct {
	Replicas                      uint                   `json:"replicas,omitempty"`
	AccountName                   string                 `json:"accountName,omitempty"`
	HostNetwork                   bool                   `json:"hostNetwork,omitempty"`
	HostAliases                   []HostAliasesInfo      `json:"hostAliases,omitempty"`
	GeneralEnvs                   []ContainerEnvInfo     `json:"generalEnvs,omitempty"`
	DnsPolicy                     string                 `json:"dnsPolicy,omitempty"`
	RestartPolicy                 string                 `json:"restartPolicy,omitempty"`
	TerminationGracePeriodSeconds uint                   `json:"terminationGracePeriodSeconds,omitempty"`
	RevisionHistoryLimit          uint                   `json:"revisionHistoryLimit,omitempty"`
	MinReadySeconds               uint                   `json:"minReadySeconds,omitempty"`
	SecurityContext               SecurityContextInfo    `json:"securityContext,omitempty"`
	ImagePullSecrets              []ImagePullSecretsInfo `json:"imagePullSecrets,omitempty"`
	ImagePullPolicy               string                 `json:"imagePullPolicy,omitempty"`
	MaxSurge                      uint                   `json:"maxSurge,omitempty"`
	MaxUnavailable                uint                   `json:"maxUnavailable,omitempty"`
	Affinity                      string                 `json:"affinity,omitempty"`
	NodeSelector                  map[string]string      `json:"nodeSelector,omitempty"`
	Containers                    []ContainerSpecInfo    `json:"containers,omitempty"`
	InitContainers                []ContainerSpecInfo    `json:"initContainers,omitempty"`
	TimezonePath                  string                 `json:"timezonePath,omitempty"`
}

type ImagePullSecretsInfo struct {
	Name string `json:"name,omitempty"`
}

type HostAliasesInfo struct {
	Ip        string   `json:"ip,omitempty"`
	Hostnames []string `json:"hostnames,omitempty"`
}

type ContainerSpecInfo struct {
	Name               string                            `json:"name,omitempty"`
	Image              string                            `json:"image,omitempty"`
	Format             string                            `json:"format,omitempty"`
	Command            []string                          `json:"command,omitempty"`
	Args               []string                          `json:"args,omitempty"`
	Env                []ContainerEnvInfo                `json:"env,omitempty"`
	EnvFrom            []map[string]ContainerEnvFormInfo `json:"envFrom,omitempty"`
	Ports              []ContainerPortInfo               `json:"ports,omitempty"`
	SecurityContext    SecurityContextInfo               `json:"securityContext,omitempty"`
	LifecyclePostStart ContainerHandleInfo               `json:"lifecyclePostStart,omitempty"`
	LifecyclePreStop   ContainerHandleInfo               `json:"lifecyclePreStop,omitempty"`
	StartupProbe       ContainerProbeInfo                `json:"startupProbe,omitempty"`
	ReadinessProbe     ContainerProbeInfo                `json:"readinessProbe,omitempty"`
	LivenessProbe      ContainerProbeInfo                `json:"livenessProbe,omitempty"`
	VolumeNames        []string                          `json:"volumeNames,omitempty"`
	FormatSpecInfo
}

type ContainerHandleInfo struct {
	Exec      map[string]interface{} `json:"exec,omitempty"`
	HttpGet   map[string]interface{} `json:"httpGet,omitempty"`
	TcpSocket map[string]interface{} `json:"tcpSocket,omitempty"`
}

type ContainerProbeInfo struct {
	ContainerHandleInfo
	InitialDelaySeconds int `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      int `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       int `json:"periodSeconds,omitempty"`
	SuccessThreshold    int `json:"successThreshold,omitempty"`
	FailureThreshold    int `json:"failureThreshold,omitempty"`
}

// fieldRef or secretKeyRef or configMapKeyRef
type ContainerEnvInfo struct {
	Name      string                            `json:"name,omitempty"`
	Value     string                            `json:"value,omitempty"`
	ValueFrom map[string]ContainerEnvKeyRefInfo `json:"valueFrom,omitempty"`
}

// configMapRef
type ContainerEnvFormInfo struct {
	Name string `json:"name,omitempty"`
}

type ContainerEnvKeyRefInfo struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

type ContainerPortInfo struct {
	ContainerPort string `json:"containerPort,omitempty"`
	HostPort      string `json:"hostPort,omitempty"`
}

type SecurityContextInfo struct {
	FsGroup      int            `json:"fsGroup,omitempty"`
	RunAsUser    int            `json:"runAsUser,omitempty"`
	RunAsGroup   int            `json:"runAsGroup,omitempty"`
	Capabilities CapabilityInfo `json:"capabilities,omitempty"`
}

type CapabilityInfo struct {
	Add  []string `json:"add,omitempty"`
	Drop []string `json:"drop,omitempty"`
}

type ConfigDataInfo struct {
	Name        string            `json:"name,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
	Items       []MountItemInfo   `json:"items,omitempty"`
	MountPath   string            `json:"mountPath,omitempty"`
	DefaultMode int               `json:"defaultMode,omitempty"`
	ReadOnly    bool              `json:"readOnly,omitempty"`
}

type MountItemInfo struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
	Mode int    `json:"mode,omitempty"`
}

/**
Type like Secret, Config, HostPath, VolumeClaim, VolumeClaimTemplate
Data, Items, DefaultMode is only for Secret or Config type
*/
type VolumeMountInfo struct {
	Name          string            `json:"name,omitempty"`
	Type          string            `json:"type,omitempty"`
	Data          map[string]string `json:"data,omitempty"`
	Items         []MountItemInfo   `json:"items,omitempty"`
	MountPath     string            `json:"mountPath,omitempty"`
	DefaultMode   int               `json:"defaultMode,omitempty"`
	ReadOnly      bool              `json:"readOnly,omitempty"`
	HostPath      string            `json:"hostPath,omitempty"`
	ClaimName     string            `json:"claimName,omitempty"`
	ClaimTemplate ClaimTemplateInfo `json:"claimTemplate,omitempty"`
}

type ClaimTemplateInfo struct {
	AccessModes      []string `json:"accessModes,omitempty"`
	StorageClassName string   `json:"storageClassName,omitempty"`
	RequestStorage   string   `json:"requestStorage,omitempty"`
}

/**
Type like ClusterIP, NodePort, Contour(ClusterIP with ingress, must set EnableHttpProxy or EnableTcpProxy or EnableWebsockets)
*/
type AccessPortalInfo struct {
	Name                  string              `json:"name,omitempty"`
	Type                  string              `json:"type,omitempty"`
	ClusterIP             string              `json:"clusterIP,omitempty"`
	SessionAffinity       string              `json:"sessionAffinity,omitempty"`
	ExternalTrafficPolicy string              `json:"externalTrafficPolicy,omitempty"`
	Routers               []ServiceRouterInfo `json:"routers,omitempty"`
}

type ServiceDataInfo struct {
	Name                  string            `json:"name,omitempty"`
	Type                  string            `json:"type,omitempty"`
	ClusterIP             string            `json:"clusterIP,omitempty"`
	SessionAffinity       string            `json:"sessionAffinity,omitempty"`
	ExternalTrafficPolicy string            `json:"externalTrafficPolicy,omitempty"`
	Ports                 []ServicePortInfo `json:"ports,omitempty"`
}

type ServicePortInfo struct {
	Name       string `json:"name,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Port       int32  `json:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	NodePort   int32  `json:"nodePort,omitempty"`
}

type ServiceRouterInfo struct {
	ServicePortInfo
	Tls                 ContourTlsInfo         `json:"tls,omitempty"`
	PreDomain           string                 `json:"preDomain,omitempty"`
	FullDomain          string                 `json:"fullDomain,omitempty"`
	RespTimeout         string                 `json:"respTimeout,omitempty"`
	IdleTimeout         string                 `json:"idleTimeout,omitempty"`
	RetryCount          int32                  `json:"retryCount,omitempty"`
	RetryPerTryTimeout  string                 `json:"retryPerTryTimeout,omitempty"`
	HTTPPrefix          []string               `json:"httpPrefix,omitempty"`
	HTTPHeader          []ContourHeaderInfo    `json:"httpHeader,omitempty"`
	WSPrefix            []string               `json:"wsPrefix,omitempty"`
	WSHeader            []ContourHeaderInfo    `json:"wsHeader,omitempty"`
	TCPEnable           bool                   `json:"tcpEnable,omitempty"`
	LoadBalanceStrategy string                 `json:"loadBalanceStrategy,omitempty"`
	CorsPolicy          ContourCorsInfo        `json:"corsPolicy,omitempty"`
	HealthCheckPolicy   ContourHealthCheckInfo `json:"healthCheckPolicy,omitempty"`
}
