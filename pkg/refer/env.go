package refer

type SyncInfo struct {
	Accept []SyncAcceptInfo `json:"accept"`
	Issue  []SyncIssueInfo  `json:"issue"`
}

type SyncAcceptInfo struct {
	Env    uint64 `json:"env"`
	Secret uint64 `json:"secret"`
}

type SyncIssueInfo struct {
	Env uint64 `json:"env"`
	Url string `json:"url"`
}

type DeployInfo struct {
	Host    string  `json:"host"`
	K8SInfo K8SConf `json:"k8s"`
	S3Info  S3Conf  `json:"s3"`
	GitInfo GitConf `json:"git"`
	RegInfo RegConf `json:"reg"`
}

type FormatInfo struct {
	// From is FormatPlainInfo
	From map[string]string `json:"from"`
	// Script is FormatPlainInfo
	Script map[string]string `json:"script"`
	// Agent is FormatAgentInfo list
	Agent []map[string]string `json:"agent"`
	// Spec is FormatSpecInfo map, this key is reference to AppInfo.Format
	Spec map[string]map[string]string `json:"spec"`
}

type FormatPlainInfo struct {
	Static string `json:"static"`
	Nginx  string `json:"nginx"`
	Java   string `json:"java"`
	Tomcat string `json:"tomcat"`
	Go     string `json:"go"`
	Python string `json:"python"`
	Node   string `json:"node"`
}

type FormatAgentInfo struct {
	Type  string `json:"type"`
	Path  string `json:"path"`
	Param string `json:"param"`
}

type FormatSpecInfo struct {
	JvmFlags      string `json:"jvmFlags"`
	RequestCpu    string `json:"requestCpu"`
	RequestMemory string `json:"requestMemory"`
	LimitCpu      string `json:"limitCpu"`
	LimitMemory   string `json:"limitMemory"`
}
