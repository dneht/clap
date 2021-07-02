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
	RevisionHistoryLimit          string                 `json:"revisionHistoryLimit,omitempty"`
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
	Name string `json:"name"`
}

type HostAliasesInfo struct {
	Ip        string   `json:"ip"`
	Hostnames []string `json:"hostnames"`
}

type ContainerSpecInfo struct {
	Name               string                            `json:"name,omitempty"`
	Image              string                            `json:"image,omitempty"`
	Format             string                            `json:"format"`
	Command            []string                          `json:"command"`
	Args               []string                          `json:"args"`
	Env                []ContainerEnvInfo                `json:"env,omitempty"`
	EnvFrom            []map[string]ContainerEnvFormInfo `json:"envFrom,omitempty"`
	Ports              []ContainerPortInfo               `json:"ports"`
	SecurityContext    SecurityContextInfo               `json:"securityContext,omitempty"`
	LifecyclePostStart ContainerHandleInfo               `json:"lifecyclePostStart,omitempty"`
	LifecyclePreStop   ContainerHandleInfo               `json:"lifecyclePreStop,omitempty"`
	StartupProbe       ContainerProbeInfo                `json:"startupProbe,omitempty"`
	ReadinessProbe     ContainerProbeInfo                `json:"readinessProbe,omitempty"`
	LivenessProbe      ContainerProbeInfo                `json:"livenessProbe,omitempty"`
	VolumeNames        []string                          `json:"volumeNames"`
	FormatSpecInfo
}

type ContainerHandleInfo struct {
	Exec      map[string]interface{} `json:"exec"`
	HttpGet   map[string]interface{} `json:"httpGet"`
	TcpSocket map[string]interface{} `json:"tcpSocket"`
}

type ContainerProbeInfo struct {
	ContainerHandleInfo
	InitialDelaySeconds int `json:"initialDelaySeconds"`
	TimeoutSeconds      int `json:"timeoutSeconds"`
	PeriodSeconds       int `json:"periodSeconds"`
	SuccessThreshold    int `json:"successThreshold"`
	FailureThreshold    int `json:"failureThreshold"`
}

// fieldRef or secretKeyRef or configMapKeyRef
type ContainerEnvInfo struct {
	Name      string                            `json:"name"`
	Value     string                            `json:"value"`
	ValueFrom map[string]ContainerEnvKeyRefInfo `json:"valueFrom"`
}

// configMapRef
type ContainerEnvFormInfo struct {
	Name string `json:"name"`
}

type ContainerEnvKeyRefInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type ContainerPortInfo struct {
	ContainerPort string `json:"containerPort"`
	HostPort      string `json:"hostPort"`
}

type SecurityContextInfo struct {
	FsGroup      int            `json:"fsGroup"`
	RunAsUser    int            `json:"runAsUser"`
	RunAsGroup   int            `json:"runAsGroup"`
	Capabilities CapabilityInfo `json:"capabilities"`
}

type CapabilityInfo struct {
	Add  []string `json:"add"`
	Drop []string `json:"drop"`
}

type ConfigDataInfo struct {
	Name        string            `json:"name"`
	Data        map[string]string `json:"data"`
	Items       []MountItemInfo   `json:"items"`
	MountPath   string            `json:"mountPath"`
	DefaultMode int               `json:"defaultMode"`
	ReadOnly    bool              `json:"readOnly"`
}

type MountItemInfo struct {
	Key  string `json:"key"`
	Path string `json:"path"`
	Mode int    `json:"mode"`
}

/**
Type like Secret, Config, HostPath, VolumeClaim, VolumeClaimTemplate
Data, Items, DefaultMode is only for Secret or Config type
*/
type VolumeMountInfo struct {
	Name          string            `json:"name"`
	Type          string            `json:"type"`
	Data          map[string]string `json:"data"`
	Items         []MountItemInfo   `json:"items"`
	MountPath     string            `json:"mountPath"`
	DefaultMode   int               `json:"defaultMode"`
	ReadOnly      bool              `json:"readOnly"`
	HostPath      string            `json:"hostPath"`
	ClaimName     string            `json:"claimName"`
	ClaimTemplate ClaimTemplateInfo `json:"claimTemplate"`
}

type ClaimTemplateInfo struct {
	AccessModes      []string `json:"accessModes"`
	StorageClassName string   `json:"storageClassName"`
	RequestStorage   string   `json:"requestStorage"`
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
	Routers               []ServiceRouterInfo `json:"routers"`
}

type ServiceDataInfo struct {
	Name                  string            `json:"name,omitempty"`
	Type                  string            `json:"type,omitempty"`
	ClusterIP             string            `json:"clusterIP,omitempty"`
	SessionAffinity       string            `json:"sessionAffinity,omitempty"`
	ExternalTrafficPolicy string            `json:"externalTrafficPolicy,omitempty"`
	Ports                 []ServicePortInfo `json:"ports"`
}

type ServicePortInfo struct {
	Name       string `json:"name,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Port       int32  `json:"port"`
	TargetPort int32  `json:"targetPort"`
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

type ContourDataInfo struct {
	Fqdn       string                `json:"fqdn,omitempty"`
	Tls        ContourTlsInfo        `json:"tls,omitempty"`
	CorsPolicy ContourCorsInfo       `json:"corsPolicy,omitempty"`
	Includes   []ContourIncludeInfo  `json:"includes,omitempty"`
	Routers    []ContourRouterInfo   `json:"routes,omitempty"`
	Tcpproxy   []ContourTcpproxyInfo `json:"tcpproxy,omitempty"`
}

type ContourTlsInfo struct {
	SecretName  string `json:"secretName,omitempty"`
	Passthrough bool   `json:"passthrough,omitempty"`
}

type ContourIncludeInfo struct {
	Name       string               `json:"name"`
	Namespace  string               `json:"namespace"`
	Conditions ContourConditionInfo `json:"conditions,omitempty"`
}

/**
LoadBalanceStrategy: RoundRobin, WeightedLeastRequest, Cookie
*/
type ContourRouterInfo struct {
	Conditions          ContourConditionInfo   `json:"conditions"`
	Services            []ContourServiceInfo   `json:"services"`
	InnerService        ContourServiceInfo     `json:"innerService"`
	RespTimeout         string                 `json:"respTimeout,omitempty"`
	IdleTimeout         string                 `json:"idleTimeout,omitempty"`
	RetryCount          uint64                 `json:"retryCount,omitempty"`
	RetryPerTryTimeout  string                 `json:"retryPerTryTimeout,omitempty"`
	EnableWebsockets    bool                   `json:"enableWebsockets,omitempty"`
	LoadBalanceStrategy string                 `json:"loadBalanceStrategy,omitempty"`
	HealthCheckPolicy   ContourHealthCheckInfo `json:"healthCheckPolicy,omitempty"`
}

type ContourTcpproxyInfo struct {
	Services []ContourServiceInfo `json:"services"`
}

type ContourConditionInfo struct {
	Prefix []string            `json:"prefix"`
	Header []ContourHeaderInfo `json:"header"`
}

type ContourHeaderInfo struct {
	Name     string `json:"name"`
	Contains string `json:"contains"`
}

type ContourServiceInfo struct {
	Name   string `json:"name"`
	Port   int32  `json:"port"`
	Weight uint32 `json:"weight"`
}

type ContourCorsInfo struct {
	AllowCredentials bool     `json:"allowCredentials,omitempty"`
	AllowOrigin      []string `json:"allowOrigin,omitempty"`
	AllowMethods     []string `json:"allowMethods,omitempty"`
	AllowHeaders     []string `json:"allowHeaders,omitempty"`
	ExposeHeaders    []string `json:"exposeHeaders,omitempty"`
	MaxAge           string   `json:"maxAge,omitempty"`
}

type ContourHealthCheckInfo struct {
	Path                    string `json:"path,omitempty"`
	IntervalSeconds         uint64 `json:"intervalSeconds,omitempty"`
	TimeoutSeconds          uint64 `json:"timeoutSeconds,omitempty"`
	UnhealthyThresholdCount uint64 `json:"unhealthyThresholdCount,omitempty"`
	HealthyThresholdCount   uint64 `json:"healthyThresholdCount,omitempty"`
}
