package refer

type AppInfo struct {
	Type      int    `json:"type,omitempty"`
	Port      string `json:"port,omitempty"`
	From      string `json:"from,omitempty"`
	Template  uint64 `json:"template,omitempty"`
	Component string `json:"component,omitempty"`
	// if app_type is none or 0, please provide ready information
	Ready  *AppReadyInfo  `json:"ready,omitempty"`
	Factor *AppFactorInfo `json:"factor,omitempty"`
	// Conf is AppConfInfo
	Conf map[string]string `json:"conf,omitempty"`
	// Code is AppCodeInfo
	Code map[string]string `json:"code,omitempty"`
	// Repo is AppRepoInfo
	Repo map[string]string `json:"repo,omitempty"`
	// Param is AppParamInfo
	Param map[string]interface{} `json:"param,omitempty"`
}

type AppReadyInfo struct {
	Url     string   `json:"url,omitempty"`
	Version string   `json:"version,omitempty"`
	Ingress []string `json:"ingress,omitempty"`
}

type AppFactorInfo struct {
	// see pkg/cloud, only: none, aliyun
	CDNProvider string `json:"cdnProvider,omitempty"`
	// properties file path, if exist, default is /app
	ConfigMouthPath string `json:"configMouthPath,omitempty"`
}

type AppConfInfo struct {
	ConfInfo
	Url string `json:"url,omitempty"`
}

type AppCodeInfo struct {
	CodeInfo
	Url string `json:"url,omitempty"`
}

type AppRepoInfo struct {
	RepoInfo
	Url string `json:"url,omitempty"`
}

/**
ServiceName is only for StatefulSet, if not set will use first AccessPortals name
if container is not set VolumeNames and is first container, all VolumeMounts will work
*/
type AppParamInfo struct {
	CommonParamInfo
	ServiceName   string             `json:"serviceName,omitempty"`
	VolumeMounts  []VolumeMountInfo  `json:"volumeMounts,omitempty"`
	AccessPortals []AccessPortalInfo `json:"accessPortals,omitempty"`
}
