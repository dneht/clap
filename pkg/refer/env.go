package refer

type SyncInfo struct {
	Accept []SyncAcceptInfo `json:"accept,omitempty"`
	Issue  []SyncIssueInfo  `json:"issue,omitempty"`
}

type SyncAcceptInfo struct {
	Env    uint64 `json:"env,omitempty"`
	Secret uint64 `json:"secret,omitempty"`
}

type SyncIssueInfo struct {
	Env uint64 `json:"env,omitempty"`
	Url string `json:"url,omitempty"`
}

type DeployInfo struct {
	Host    string  `json:"host,omitempty"`
	K8SInfo K8SConf `json:"k8s,omitempty"`
	S3Info  S3Conf  `json:"s3,omitempty"`
	GitInfo GitConf `json:"git,omitempty"`
	RegInfo RegConf `json:"reg,omitempty"`
}

type FormatInfo struct {
	// From is FormatPlainInfo
	From map[string]string `json:"from,omitempty"`
	// Script is FormatPlainInfo
	Script map[string]string `json:"script,omitempty"`
	// Agent is FormatAgentInfo list
	Agent []map[string]string `json:"agent,omitempty"`
	// Spec is FormatSpecInfo map, this key is reference to AppInfo.Format
	Spec map[string]map[string]string `json:"spec,omitempty"`
}

type FormatPlainInfo struct {
	Static string `json:"static,omitempty"`
	Nginx  string `json:"nginx,omitempty"`
	Java   string `json:"java,omitempty"`
	Tomcat string `json:"tomcat,omitempty"`
	Go     string `json:"go,omitempty"`
	Python string `json:"python,omitempty"`
	Node   string `json:"node,omitempty"`
}

type FormatAgentInfo struct {
	Type  string `json:"type,omitempty"`
	Path  string `json:"path,omitempty"`
	Param string `json:"param,omitempty"`
}

type FormatSpecInfo struct {
	JvmFlags      string `json:"jvmFlags,omitempty"`
	RequestCpu    string `json:"requestCpu,omitempty"`
	RequestMemory string `json:"requestMemory,omitempty"`
	LimitCpu      string `json:"limitCpu,omitempty"`
	LimitMemory   string `json:"limitMemory,omitempty"`
}
