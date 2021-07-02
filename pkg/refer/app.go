package refer

type AppInfo struct {
	Type      int    `json:"type"`
	Port      string `json:"port"`
	From      string `json:"from"`
	Template  uint64 `json:"template"`
	Component string `json:"component"`
	// if app_type is none or 0, please provide ready information
	Ready  *AppReadyInfo  `json:"ready"`
	Factor *AppFactorInfo `json:"factor"`
	// Conf is AppConfInfo
	Conf map[string]string `json:"conf"`
	// Code is AppCodeInfo
	Code map[string]string `json:"code"`
	// Repo is AppRepoInfo
	Repo map[string]string `json:"repo"`
	// Param is AppParamInfo
	Param map[string]interface{} `json:"param"`
}

type AppReadyInfo struct {
	Url     string   `json:"url"`
	Version string   `json:"version"`
	Exec    []string `json:"exec"`
}

type AppFactorInfo struct {
	NeedPublic bool `json:"needPublic"`
	// see pkg/cloud, only: none, aliyun
	CDNProvider string `json:"cdnProvider"`
}

type AppConfInfo struct {
	ConfInfo
	Url string `json:"url"`
}

type AppCodeInfo struct {
	CodeInfo
	Url string `json:"url"`
}

type AppRepoInfo struct {
	RepoInfo
	Url string `json:"url"`
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
