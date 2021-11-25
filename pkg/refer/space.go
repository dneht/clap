package refer

type SpaceInfo struct {
	Conf SpaceConfInfo `json:"conf,omitempty"`
	Code SpaceCodeInfo `json:"code,omitempty"`
	Repo SpaceRepoInfo `json:"repo,omitempty"`
	// Param is SpaceParamInfo
	Param map[string]interface{} `json:"param,omitempty"`
}

type ConfInfo struct {
	Type   string `json:"type,omitempty"`
	Space  string `json:"space,omitempty"`
	Group  string `json:"group,omitempty"`
	Secret string `json:"secret,omitempty"`
}

type SpaceConfInfo struct {
	ConfInfo
	Base string `json:"base,omitempty"`
}

type CodeInfo struct {
	Type   string `json:"type,omitempty"`
	Branch string `json:"branch,omitempty"`
}

type SpaceCodeInfo struct {
	CodeInfo
	Base string `json:"base,omitempty"`
}

type RepoInfo struct {
	Type string `json:"type,omitempty"`
}

type SpaceRepoInfo struct {
	RepoInfo
	Base string `json:"base,omitempty"`
}

type SpaceParamInfo struct {
	CommonParamInfo
	Domain string `json:"domain,omitempty"`
}
