package refer

type SpaceInfo struct {
	Conf SpaceConfInfo `json:"conf"`
	Code SpaceCodeInfo `json:"code"`
	Repo SpaceRepoInfo `json:"repo"`
	// Param is SpaceParamInfo
	Param map[string]interface{} `json:"param"`
}

type ConfInfo struct {
	Type   string `json:"type"`
	Space  string `json:"space"`
	Group  string `json:"group"`
	Secret string `json:"secret"`
}

type SpaceConfInfo struct {
	ConfInfo
	Base string `json:"base"`
}

type CodeInfo struct {
	Type   string `json:"type"`
	Branch string `json:"branch"`
}

type SpaceCodeInfo struct {
	CodeInfo
	Base string `json:"base"`
}

type RepoInfo struct {
	Type string `json:"type"`
}

type SpaceRepoInfo struct {
	RepoInfo
	Base string `json:"base"`
}

type SpaceParamInfo struct {
	CommonParamInfo
	Domain string `json:"domain,omitempty"`
}
