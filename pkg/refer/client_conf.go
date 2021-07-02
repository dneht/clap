package refer

type K8SConf struct {
	Inner       bool   `json:"inner"`
	Master      string `json:"master"`
	Version     string `json:"version"`
	ClusterName string `json:"clusterName"`
	ClusterUser string `json:"clusterUser"`
	ClientCert  string `json:"clientCert"`
	ClientKey   string `json:"clientKey"`
}

type S3Conf struct {
	Region     string `json:"region"`
	Endpoint   string `json:"endpoint"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey"`
	ForcePath  bool   `json:"forcePath"`
	DisableSSL bool   `json:"disableSSL"`
}

type GitConf struct {
	Type       string `json:"type"`
	Server     string `json:"server"`
	AuthToken  string `json:"token"`
	PublicKey  string `json:"public"`
	PrivateKey string `json:"private"`
}

type RegConf struct {
	Type    string            `json:"type"`
	Server  string            `json:"server"`
	AuthMap map[string]string `json:"auths"`
}

type SrcConf map[string]interface{}
