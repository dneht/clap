package refer

type K8SConf struct {
	Inner       bool   `json:"inner,omitempty"`
	Master      string `json:"master,omitempty"`
	Version     string `json:"version,omitempty"`
	ClusterName string `json:"clusterName,omitempty"`
	ClusterUser string `json:"clusterUser,omitempty"`
	ClientCert  string `json:"clientCert,omitempty"`
	ClientKey   string `json:"clientKey,omitempty"`
}

type S3Conf struct {
	Region     string `json:"region,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	AccessKey  string `json:"accessKey,omitempty"`
	SecretKey  string `json:"secretKey,omitempty"`
	ForcePath  bool   `json:"forcePath,omitempty"`
	DisableSSL bool   `json:"disableSSL,omitempty"`
}

type GitConf struct {
	Type       string `json:"type,omitempty"`
	Server     string `json:"server,omitempty"`
	AuthToken  string `json:"token,omitempty"`
	PublicKey  string `json:"public,omitempty"`
	PrivateKey string `json:"private,omitempty"`
}

type RegConf struct {
	Type    string            `json:"type,omitempty"`
	Server  string            `json:"server,omitempty"`
	AuthMap map[string]string `json:"auths,omitempty"`
}

type SrcConf map[string]interface{}
